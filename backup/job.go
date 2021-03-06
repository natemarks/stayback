// backup provides the functions for backing up local paths
// The work is organized as a single job with tasks for each local target path
// The important job input is validated first so we can fail quickly
// than the  target handler is run for each target
package backup

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/natemarks/awsgrips/s3"
	"github.com/natemarks/stayback/shell"
	"github.com/rs/zerolog"
)

func CurrentTime() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

// Job is the backup job definition
type Job struct {
	Source          string   `json:"source"`          // ex. file, url, telepathy
	Id              string   `json:"id"`              // ex. 20220103-080910
	Recipient       string   `json:"recipient"`       // Recipient email identifies GPG public key for encryption
	HomeDirectory   string   `json:"homeDirectory"`   // assumed root for relative backup target paths
	BackupDirectory string   `json:"backupDirectory"` // local backup location with a .tmp working ssub-dir
	S3Bucket        string   `json:"s3Bucket"`        // backup s3 bucket
	EncryptedDirs   []string `json:"encryptedDirs"`   // list of abs or relative dirs to back up wit encryption
	UnEncryptedDirs []string `json:"unEncryptedDirs"` // list of absolute or relative dirs to back up without encryption
}

// TargetDirsExist Returns a merge map of the absolute directories and the boolean result of their existence
// We want to check every directory and report all problems, so they can all be solved at once
// log fatal if this fails
func (c *Job) TargetDirsExist(log *zerolog.Logger) (err error) {
	// the target lists might be a mix of absolute paths and paths relative to the home directory
	// clean all the path strings into strings that look like absolutes
	c.EncryptedDirs = cleanTargets(c.EncryptedDirs, c.HomeDirectory)
	c.UnEncryptedDirs = cleanTargets(c.UnEncryptedDirs, c.HomeDirectory)

	for _, v := range c.UnEncryptedDirs {
		_, fErr := os.Stat(v)
		if fErr != nil {
			err = errors.New("some target paths do not exist")
			log.Error().Msgf("target does not exist: %s", v)
			continue
		}
		log.Debug().Msgf("target exists: %s", v)
	}
	for _, v := range c.EncryptedDirs {
		_, fErr := os.Stat(v)
		if fErr != nil {
			err = errors.New("some target paths do not exist")
			log.Error().Msgf("target does not exist: %s", v)
			continue
		}
		log.Debug().Msgf("target exists: %s", v)
	}
	return err
}

// MakeRestoreDir Return an error if the directory already exists, otherwise create
func (c Job) MakeRestoreDir() (err error) {
	restoreDir := path.Join(c.BackupDirectory, c.Id)
	_, err = os.Stat(restoreDir)
	// Return error if restore directory exists
	if err == nil {
		msg := fmt.Sprintf("Restore directory already exists: %s", restoreDir)
		return errors.New(msg)
	}
	err = shell.MkdirP(restoreDir)
	return err
}

// Restore Download backup data to local
func (c Job) Restore() (err error) {
	restoreDir := path.Join(c.BackupDirectory, c.Id)
	_, err = shell.RunAndWait("aws", []string{"s3", "sync", c.S3Uri(), restoreDir})
	return err
}

// Return the S3 Uri for the job
func (c Job) S3Uri() (s3Uri string) {
	s3path := fmt.Sprintf("stayback/%s/", c.Id)
	return fmt.Sprintf("s3://%s/%s", c.S3Bucket, s3path)
}

// CreateS3JobPath creates the s3 destination path
// log fatal if this fails
func (c Job) CreateS3JobPath(log *zerolog.Logger) (err error) {
	s3path := fmt.Sprintf("stayback/%s/", c.Id)
	uri := fmt.Sprintf("s3://%s/%s", c.S3Bucket, s3path)
	err = s3.CreatePath(c.S3Bucket, s3path)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create s3 path: %s", uri)
	}
	log.Debug().Msgf("Created s3 path: %s", uri)
	return err
}

// Execute iterates the targets and runs the backups
func (c Job) Execute() (err error) {

	return err
}

// NewJobFromFile returns a JOb object from the json file path provided
func NewJobFromFile(fPath string) (result Job, err error) {
	if fPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return result, err
		}
		fPath = path.Join(homeDir, ".stayback")
	}

	file, err := ioutil.ReadFile(fPath)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(file, &result)
	return result, err
}

// makeAbsolute if the give dri doesn't begin with a path separator, join it onto the default root
func makeAbsolute(dir, defaultRoot string) string {
	var first rune
	for _, c := range dir {
		first = c
		break
	}
	if first == os.PathSeparator {
		return dir
	}
	return path.Join(defaultRoot, dir)
}

// cleanTargets converts each directory entry to an absolute path using a default root for relative directories
// then it sorts the list and removes the duplicates
func cleanTargets(tList []string, defaultRoot string) (oList []string) {
	var absList []string

	// go through the targets and ensure each is an absolute path
	for _, v := range tList {
		abs := makeAbsolute(v, defaultRoot)
		absList = append(absList, abs)
	}

	// sort the targets so duplicates will all be grouped together
	sort.Strings(absList)

	// loop through the absList and only copy each entry from absList to oList if it
	// is different from the previous entry
	for i, v := range absList {
		// skip the predecessor comparison on the first item
		if i == 0 {
			oList = append(oList, v)
		} else {
			if v != absList[i-1] {
				oList = append(oList, v)
			}
		}
	}

	return oList
}

// TargetHandlerInput is the input required to backup a single target
type TargetHandlerInput struct {
	Target    string          // Target path to be backed up
	Encrypt   bool            // Encrypt the target before uploading to s3
	Id        string          // Id job identifier
	Local     string          // Local job backup destination
	Recipient string          // Recipient email identifies GPG public key for encryption
	S3Bucket  string          // backup s3 bucket
	Logger    *zerolog.Logger // logger pointer
}

// TargetHandler a local target path
//- backup the target to a tarball in job directory.
//- encrypt the tarball in the job directory to *.tar.gz.asc (output should be ascii-armor and with --openpgp option)
//- overwrite-move the unencrypted tarball to the destination directory
//- upload the encrypted file to the 3 destination then delete the encrypted file
// each exceution should have a unique logging context so we known which backup props are generating  a given log
// message
func TargetHandler(input TargetHandlerInput) (err error) {
	// the absolute path of the target is converted ot base64 and that's used for the tarball base file name
	// /home/myhome/.ssh  -> L2hvbWUvbXlob21lLy5zc2gK
	jobDir := path.Join(input.Local, input.Id)
	baseFileName := base64.StdEncoding.EncodeToString([]byte(input.Target))
	tempTarball := path.Join(jobDir, baseFileName+".tar.gz")
	localTarball := path.Join(input.Local, baseFileName+".tar.gz")

	// compress the target to the a file in the job dir
	input.Logger.Debug().Msgf("compressing %s -> %s", input.Target, tempTarball)
	_, err = shell.RunAndWait("tar", []string{"-cpzvf", tempTarball, input.Target})
	if err != nil {
		input.Logger.Error().Err(err).Msgf("failed: compressing %s -> %s", input.Target, tempTarball)
		return err
	}
	input.Logger.Debug().Msgf("success: compressing %s -> %s", input.Target, tempTarball)

	// delete pre-existing tarball in local directory
	input.Logger.Debug().Msgf("deleting old local backup:  %s", localTarball)
	_, err = shell.RunAndWait("rm", []string{"-f", localTarball})
	if err != nil {
		input.Logger.Error().Err(err).Msgf("failed to delete old tarball: %s", localTarball)
		return err
	}
	input.Logger.Debug().Msgf("deleted old tarball: %s", localTarball)

	// copy new tarball from job to local
	input.Logger.Debug().Msgf("copying new tarball %s -> %s", tempTarball, localTarball)
	_, err = shell.RunAndWait("cp", []string{tempTarball, localTarball})
	if err != nil {
		input.Logger.Error().Err(err).Msgf("failed to copy new tarball %s -> %s", tempTarball, localTarball)
		return err
	}
	input.Logger.Debug().Msgf("copied new tarball %s -> %s", tempTarball, localTarball)

	// gpg --openpgp --batch --yes --output \
	//  "${1}.gpg" --encrypt --recipient "${recipient}" "${1}"
	if input.Encrypt {
		// encrypt job/tarball -> job/tarball.asc
		input.Logger.Debug().Msgf("encrypting %s -> %s", tempTarball, tempTarball+".asc")
		_, err = shell.RunAndWait("gpg", []string{
			"--openpgp",
			"--armor",
			"--batch",
			"--yes",
			"--encrypt",
			"--recipient",
			input.Recipient,
			tempTarball,
		})
		if err != nil {
			input.Logger.Error().Err(err).Msgf("failed: encrypting %s -> %s", tempTarball, tempTarball+".asc")
			return err
		}
		// delete the unencrypted tarball from the job directory
		input.Logger.Debug().Msgf("deleting unencrypted temp tarball:  %s", tempTarball)
		_, err = shell.RunAndWait("rm", []string{"-f", tempTarball})
		if err != nil {
			input.Logger.Error().Err(err).Msgf("failed to delete unencrypted temp tarball: %s", localTarball)
			return err
		}
		input.Logger.Debug().Msgf("deleted unencrypted temp tarball: %s", localTarball)

	}

	return err
}

// LatestKeyFromS3 find the most recent S3 backup job
func (c Job) LatestKeyFromS3() (key string, err error) {
	var latestOject types.Object
	results, err := s3.ListObjects(c.S3Bucket, "stayback/")
	if err != nil {
		return "", err
	}
	for _, cc := range results {
		if (latestOject.LastModified == nil) || (cc.LastModified.After(*latestOject.LastModified)) {
			latestOject = cc
		}
	}

	return *latestOject.Key, err
}
