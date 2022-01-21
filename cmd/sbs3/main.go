package main

import (
	"fmt"
	"github.com/natemarks/stayback/version"
	"github.com/rs/zerolog"
	"os"
)

func run() (err error) {
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Debug().Msgf("Starting")

	// validate target paths.  don't want to have to rerun a long job if we don't have to
	fmt.Println("Print job summary")
	fmt.Println("Press 'c' to continue")
	// backup to  tmp working directory
	// if all the targets are successful, copy the files from tmp to local backup. the names are base64 encodes of the absolute path of  the target directory, so it should only overwrite after a complete success (not a per job success)
	//and only keep the latest backup of each target.
	//  additionally, the whole working directory will be synced to the s3 backup path in a folder named for the job identifier

	fmt.Println("backup comp")
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
