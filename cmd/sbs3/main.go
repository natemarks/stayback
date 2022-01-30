// backup local targets to s3 in parallel
// fail fast for missing local targets and for bad s3 access
package main

// TODO: have working directories be per job id
import (
	"github.com/natemarks/stayback/backup"
	"github.com/natemarks/stayback/version"
	"github.com/rs/zerolog"
	"os"
)

func run() (err error) {
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Debug().Msgf("Starting")
	job, err := backup.NewJobFromFile("/Users/nmarks/.stayback.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	// log error for all of the targets that don't exist
	// if any targets didn't exist, log fatal
	err = job.TargetDirsExist(&logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	// log fatal if we fail to create the S3 job path
	// this is a pretty good access check
	err = job.CreateS3JobPath()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	for _, v := range job.EncryptedDirs {
		t := backup.TargetHandlerInput{
			Target:    v,
			Encrypt:   true,
			Id:        job.Id,
			Local:     job.BackupDirectory,
			Recipient: job.Recipient,
			S3Bucket:  job.S3Bucket,
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
		}
		err = backup.TargetHandler(t)
		if err != nil {
			return err
		}
	}
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
