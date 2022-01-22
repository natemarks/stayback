package main

import "github.com/natemarks/stayback/cmd/sbcompress/cmd"

//	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
//	logger.Info().Msgf("Starting")
func main() {
	cmd.Execute()
}
