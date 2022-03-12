// restore job from S3 to stayback local working directory
// given a job id (ex. 20220306-070110 )
// see if that direcotry exists in the local stayback directory and fail
package main

import (
	"os"
	"path"
	"strings"

	"github.com/natemarks/stayback/shell"

	"github.com/natemarks/stayback/backup"
	"github.com/natemarks/stayback/version"
	"github.com/rs/zerolog"
)

func JobIdFromKey(key string) (jobId string, err error) {
	return "", err
}

func KeyToJobId(key string) string {
	parts := strings.Split(key, "/")
	return parts[1]
}

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
	latestKey, err := job.LatestKeyFromS3()
	// set the job id form the current time
	job.Id = KeyToJobId(latestKey)

	// Setup the restore directory. Fatal if it already exists
	err = job.MakeRestoreDir()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	logger.Debug().Msgf("downloading: %s -> %s", job.S3Uri(), path.Join(job.BackupDirectory, job.Id))
	job.Restore()
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
