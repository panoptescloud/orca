package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/adamkirk/orca/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const configFileName = "orca.yaml"
const configFileOverrideEnv = "ORCA_CONFIG_PATH"
const configToolsPathOverrideEnv = "ORCA_TOOLS_PATH"

type runEHandlerFunc func(cmd *cobra.Command, args []string) error
type runHandlerFunc func(cmd *cobra.Command, args []string)

var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

var svcContainer *services = &services{}

var rootCmd = &cobra.Command{
	Use:   "orca",
	Short: "Orca is used to orchestrate complex development environments with docker compose.",
	Long: `This can be used to manage multiple projects in different repositories
so that they can be started/stopped as a singular unit, and provides common
utilities to interact with services form anywhere on the host.`,
	RunE: handleGroup,
}

func errorHandlerWrapper(f runEHandlerFunc, errorExitCode int) runHandlerFunc {
	return func(cmd *cobra.Command, args []string) {
		err := f(cmd, args)

		if err != nil {
			// Set as a debug level here as it should already be logged earlier
			// in the stack
			slog.Debug("unhandled error", "err", err)
			os.Exit(errorExitCode)
		}
	}
}

func handleGroup(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func getToolsDir() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configFile := fmt.Sprintf("%s/.orca/bin", homeDir)

	if override, found := os.LookupEnv(configToolsPathOverrideEnv); found {
		configFile = override
	}

	return configFile
}

func getOverlayDir() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return fmt.Sprintf("%s/.orca/overlays", homeDir)
}

func getTLSDir() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	dir := fmt.Sprintf("%s/.orca/tls", homeDir)

	return dir
}

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configFile := fmt.Sprintf("%s/.orca/%s", homeDir, configFileName)

	if override, found := os.LookupEnv(configFileOverrideEnv); found {
		configFile = override
	}

	return configFile
}

func getWorkingDir() string {
	wd, err := os.Getwd()

	cobra.CheckErr(err)

	return wd
}

func getWorkingDirParent() string {
	wd := getWorkingDir()

	dir := path.Dir(fmt.Sprintf("%s/../", wd))

	_, err := os.Stat(dir)

	cobra.CheckErr(err)

	return dir
}

func init() {
	cfg := svcContainer.GetConfig()

	err := cfg.LoadOrCreate()
	cobra.CheckErr(err)

	cobra.OnInitialize(bootstrap)

	// Persistent flags
	rootCmd.PersistentFlags().String("log-level", cfg.GetLoggingLevel(), "Log level to use, one of: debug, info, warn, error, none. Defaults to none, as most errors are already surfaced anyway.")
	rootCmd.PersistentFlags().String("log-format", cfg.GetLoggingFormat(), "log format to use")

	cobra.CheckErr(viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("log-level")))
	cobra.CheckErr(viper.BindPFlag("logging.format", rootCmd.PersistentFlags().Lookup("log-format")))
}

func addServiceOption(cmd *cobra.Command, required bool) {
	cmd.Flags().StringP("service", "s", "", "The service within the project to exec into")

	if required {
		cmd.MarkFlagRequired("service")
	}
}

func addWorkspaceOption(cmd *cobra.Command, required bool) {
	cmd.Flags().StringP("workspace", "w", "", "The name of the workspace to run this command for.")

	if required {
		cmd.MarkFlagRequired("workspace")
	}
}

func addProjectOption(cmd *cobra.Command) {
	cmd.Flags().StringP("project", "p", "", "The name of the project within the workspace to run this command for.")
}

func bootstrap() {
	cfg := svcContainer.GetConfig()

	// Tell viper to replace . in nested path with underscores
	// e.g. logging.level becomes LOGGING_LEVEL
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetEnvPrefix("orca")
	viper.AutomaticEnv()

	err := viper.Unmarshal(cfg.GetRuntimeConfig())
	cobra.CheckErr(err)

	h, err := logging.NewSlogHandler(cfg)

	cobra.CheckErr(err)

	slog.SetDefault(slog.New(h))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
