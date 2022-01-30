// backup provides the functions for backing up local paths
// The work is organized as a single job with tasks for each local target path
// The important job input is validated first so we can fail quickly
// than the  target handler is run for each target
package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/natemarks/awsgrips/s3"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"
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
func (c Job) TargetDirsExist(log *zerolog.Logger) (err error) {
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

// CreateS3JobPath creates the s3 destination path
// log fatal if this fails
func (c Job) CreateS3JobPath(log *zerolog.Logger) (err error) {
	path := fmt.Sprintf("stayback/%s/", c.Id)
	uri := fmt.Sprintf("s3://%s/%s", c.S3Bucket, path)
	err = s3.CreatePath(c.S3Bucket, path)
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
	Target    string `json:"target"`    // Target path to be backed up
	Encrypt   bool   `json:"encrypt"`   // Encrypt the target before uploading to s3
	Id        string `json:"id"`        // Id job identifier
	Local     string `json:"Local"`     // Local job backup destination
	Recipient string `json:"recipient"` // Recipient email identifies GPG public key for encryption
	S3Bucket  string `json:"s3Bucket"`  // backup s3 bucket
}

// TargetHandler a local target path
//- backup the target to a tarball in job directory.
//- encrypt the tarball in the job directory to *.tar.gz.asc (output should be ascii-armor and with --openpgp option)
//- overwrite-move the unencrypted tarball to the destination directory
//- upload the encrypted file to the 3 destination then delete the encrypted file
// each exceution should have a unique logging context so we known which backup props are generating  a given log
// message
func TargetHandler(input TargetHandlerInput) (err error) {
	return err
}
