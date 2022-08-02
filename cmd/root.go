package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cadyrov/occam/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct{}

// nolint: gochecknoglobals
var (
	cfgFileName = ""
	cfgFilePath = ""
	envPrefix   = ""
	cnf         = config.Config{}

	rootCmd = &cobra.Command{
		Use:   "root",
		Short: "root",
		Long:  `root`,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

// nolint:gochecknoinits
func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFileName, "config-file-name",
		"config", "set alternative config file name without extensions")
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config-file-path", ".",
		"set alternative relative config file path")
	rootCmd.PersistentFlags().StringVar(&envPrefix, "env-prefix", "OCCAM",
		"set alternative env prefix")

	rootCmd.AddCommand(cliCmd)
}

func initLogger() zerolog.Logger {
	var level zerolog.Level

	switch cnf.Project.Log.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "warning":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}

	if cnf.Project.Log.Output != "" {
		fl, err := os.Create(cnf.Project.Log.Output)
		if err != nil {
			return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
		}

		return zerolog.New(fl).With().Timestamp().Logger().Level(level)
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
}

func initConfig() {
	v := viper.New()

	v.SetConfigName(cfgFileName)

	v.AddConfigPath(cfgFilePath)

	var errViperCNF *viper.ConfigFileNotFoundError

	if err := v.ReadInConfig(); err != nil {
		if errors.As(err, &errViperCNF) {
			panic(err)
		}
	}

	v.SetEnvPrefix(envPrefix)

	v.AutomaticEnv()

	r := strings.NewReplacer(".", "_")

	v.SetEnvKeyReplacer(r)

	if err := v.Unmarshal(&cnf); err != nil {
		panic(err)
	}
}
