// restore job from S3 to stayback local working directory
// given a job id (ex. 20220306-070110 )
// see if that direcotry exists in the local stayback directory and fail
package main

import (
	"fmt"
	"os"
	"path"

	"github.com/natemarks/stayback/backup"
	"github.com/natemarks/stayback/version"
	"github.com/rs/zerolog"
)

func JobIdFromKey(key string) (jobId string, err error) {
	return "", err
}

func run() (err error) {
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Debug().Msgf("Starting")
	job, err := backup.NewJobFromFile("/Users/nmarks/.stayback.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	latestKey, err := job.LatestKeyFromS3()
	fmt.Printf(latestKey)
	// set the job id form the current time
	job.Id = "20220306-070110"

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
