package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	ConfigLogFile   = "log_file"
	ConfigVerbosity = "verbosity"
)

var (
	logger *zerolog.Logger
)

func Init() error {
	levelName := viper.GetString(ConfigVerbosity)
	outputFilePath := viper.GetString(ConfigLogFile)

	level, err := zerolog.ParseLevel(levelName)
	if err != nil {
		return fmt.Errorf("failed to parse log level '%s': %v", levelName, err)
	}

	if level == zerolog.Disabled || outputFilePath == "" {
		l := zerolog.Nop()
		logger = &l
		return nil
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
	if err != nil {
		return fmt.Errorf("failed to open log file '%s': %v", outputFilePath, err)
	}

	zerolog.SetGlobalLevel(level)

	l := zerolog.New(outputFile).With().
		Timestamp().
		Logger()

	logger = &l

	return nil
}

func Get() *zerolog.Logger {
	if logger == nil {
		return &zerolog.Logger{}
	}

	return logger
}
