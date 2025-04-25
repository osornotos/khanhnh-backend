package cmd

import (
	"io"
	"khanhnh-backend/internal"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "khanhnh-backend",
	Long:  ``,
}

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")

}

var appConfig internal.Config
var cfgFile string

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Stack().Err(err).Msgf("fail to read config file")
	}
	log.Info().Msgf("Using config file: %v", viper.ConfigFileUsed())
	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatal().Stack().Err(err).Msgf("fail to parse config file")
	}
	err := appConfig.Validate()
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("invalid config file")
	}

	log.Info().Msgf("Running app version: %+v", appConfig.Version)
	log.Info().Msgf("Running app with config: %+v", appConfig.Safe())

	//sentryWriter, err := sentryWriter()
	setupLogger(appConfig.Debug)
}

func sentryWriter() (io.Writer, error) {
	//sentryWriter, err := zlogsentry.New(appConfig.Sentry.Dsn,
	//	zlogsentry.WithEnvironment(appConfig.Env),
	//	zlogsentry.WithRelease(appConfig.Version),
	//	zlogsentry.WithLevels([]zerolog.Level{
	//		zerolog.WarnLevel,
	//		zerolog.ErrorLevel,
	//		zerolog.FatalLevel,
	//		zerolog.PanicLevel,
	//	}...),
	//)
	//
	//if err != nil {
	//	return err
	//}
	return nil, nil
}

func setupLogger(debugLevel bool, writers ...io.Writer) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	level := zerolog.InfoLevel
	if debugLevel {
		level = zerolog.DebugLevel
	}

	var logger zerolog.Logger

	if len(writers) > 0 {
		writers = append(writers, os.Stdout)
		multiWriter := zerolog.MultiLevelWriter(writers...)
		logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger().Level(level)
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Level(level)
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		level := 0
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				level += 1
				if level >= 3 {
					break
				}
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = logger
}

func SetCfgFile(fileDir string) {
	cfgFile = fileDir
}

func GetConfig() internal.Config {
	return appConfig
}
