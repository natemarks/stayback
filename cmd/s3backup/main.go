// backup local targets to s3 in parallel
// fail fast for missing local targets and for bad s3 access
package main

// TODO: have working directories be per job id
import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/natemarks/stayback/backup"
	"github.com/natemarks/stayback/shell"
	"github.com/natemarks/stayback/version"
	"github.com/rs/zerolog"
)

func run() (err error) {
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Debug().Msgf("Starting")
	configFile, err := shell.DefaultConfigFilePath()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	job, err := backup.NewJobFromFile(configFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	// set the job id form the current time
	t := time.Now()
	job.Id = fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	// log error for all of the targets that don't exist
	// if any targets didn't exist, log fatal
	err = job.TargetDirsExist(&logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	// log fatal if we fail to create the S3 job path
	// this is a pretty good access check
	err = job.CreateS3JobPath(&logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	// create the temporary job directory
	err = shell.MkdirP(path.Join(job.BackupDirectory, job.Id))
	for _, v := range job.EncryptedDirs {
		t := backup.TargetHandlerInput{
			Target:    v,
			Encrypt:   true,
			Id:        job.Id,
			Local:     job.BackupDirectory,
			Recipient: job.Recipient,
			S3Bucket:  job.S3Bucket,
			Logger:    &logger,
		}
		err = backup.TargetHandler(t)
		if err != nil {
			return err
		}
	}
	for _, v := range job.UnEncryptedDirs {
		t := backup.TargetHandlerInput{
			Target:    v,
			Encrypt:   false,
			Id:        job.Id,
			Local:     job.BackupDirectory,
			Recipient: job.Recipient,
			S3Bucket:  job.S3Bucket,
			Logger:    &logger,
		}
		err = backup.TargetHandler(t)
		if err != nil {
			return err
		}
	}
	// finally runupload like:
	///aws s3 cp  ~/.stayback/20220130-143927/ \
	// s3://com.imprivata.371143864265.us-east-1.personal/stayback/20220130-143927 --recursive
	jobDir := path.Join(job.BackupDirectory, job.Id)
	s3Uri := fmt.Sprintf("s3://%s/stayback/%s/", job.S3Bucket, job.Id)
	logger.Debug().Msgf("uploading files: %s -> %s", jobDir, s3Uri)
	_, err = shell.RunAndWait("aws", []string{"s3", "cp", jobDir, s3Uri, "--recursive"})
	if err != nil {
		logger.Error().Err(err).Msgf("failed to upload files: %s -> %s", jobDir, s3Uri)
		return err
	}
	logger.Debug().Msgf("uploadied files: %s -> %s", jobDir, s3Uri)

	return err
}

// main just wraps run and sets the exit code
func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
