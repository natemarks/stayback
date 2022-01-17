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

	fmt.Println("Print job summary")
	fmt.Println("Press 'c' to continue")
	fmt.Println("backup complete: bucket/stayback/20220101030405/")
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
