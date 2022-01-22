package cmd

import (
	"errors"
	"fmt"
	"github.com/natemarks/stayback/backup"
	"github.com/natemarks/stayback/version"
	"github.com/rs/zerolog"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool = false

var rootCmd = &cobra.Command{
	Use:   "sbcompress",
	Short: "Back up a list of directories to tarballs in a destination path",
	Long: `
Each tarball is placed in the destination directory and named with the base64 encode
of its source absolute path. The tar is first created in a subdirectory of the destination
directory named for the job ID. sbcompress uses the timestamp for the job id.


this command:
sbcompress --destination /my/usb/drive /home/usr/pictures /opt/coolstuff
would backup:
/home/usr/pictures -> /my/usb/drive/L2hvbWUvdXNyL3BpY3R1cmVzCg==.tar.gz
/opt/coolstuff -> /my/usb/drive/L29wdC9jb29sc3R1ZmYK.tar.gz

The tar commands would first use the job ID temporary directory
tar -czvf /my/usb/drive/${JOBID}/L2hvbWUvdXNyL3BpY3R1cmVzCg==.tar.gz /home/usr/pictures
tar -czvf /my/usb/drive/${JOBID}/L29wdC9jb29sc3R1ZmYK.tar.gz /opt/coolstuff

If the tar succeeds, the files are moved to the destination directory, overwriting existing files`,

	RunE: func(cmd *cobra.Command, args []string) error {
		logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if Verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		job := backup.Job{
			Source:          "",
			Id:              "sf",
			HomeDirectory:   "/Users/nmarks",
			BackupDirectory: "/Users/nmarks/.stayback",
			S3Bucket:        "",
			EncryptedDirs:   nil,
			UnEncryptedDirs: []string{"/users/nmarks/.ssh", "/users/nmarks/.aws"},
		}
		dirsExist := job.TargetDirsExist()
		err := EvaluateDirCheck(dirsExist, &logger)
		return err
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// EvaluateDirCheck use the logger while evaluating each entry and return err  ir s
func EvaluateDirCheck(report map[string]bool, log *zerolog.Logger) (err error) {
	var result bool = true
	for k, v := range report {
		if v {
			log.Debug().Msgf("success: %s exists", k)
		} else {
			log.Info().Msgf("fail: %s does not exist", k)
			result = false
		}
	}
	if !result {
		return errors.New("one or more source directories does not exist")
	}

	return nil
}
