package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/panoptescloud/orca/internal/config"
	"github.com/panoptescloud/orca/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type runEHandlerFunc func(cmd *cobra.Command, args []string, config *config.Config) error
type runHandlerFunc func(cmd *cobra.Command, args []string)

var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

var appCfg *config.Config
var svcContainer *services = &services{}

var rootCmd = &cobra.Command{
	Use:   "orca",
	Short: "Orca is used to orchestrate complex development environments with docker compose.",
	Long: `This can be used to manage multiple projects in different repositories
so that they can be started/stopped as a singular unit, and provides common
utilities to interact with services form anywhere on the host.`,
	RunE: handleGroup,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version you're using.",
	Run:   errorHandlerWrapper(handleVersion, 1),
}

var gCmd = &cobra.Command{
	Use:   "g",
	Short: "Commands related to git.",
	RunE:  handleGroup,
}

var gCoCmd = &cobra.Command{
	Use:   "co",
	Short: "Checkout a branch for a git repository.",
	Long: `Searches for a branch with the name provided as an argument. If a single branch is
found, it will be checked out. If multiple are found will provide a list to select from.`,
	Run: errorHandlerWrapper(handleGCo, 1),
}

var gBranchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "Searches branches for the git repository.",
	Long: `Will search for any branches containing the given search term (case-insensitive).
If no search term is given, will list all branches.`,
	Run: errorHandlerWrapper(handleGBranches, 1),
}

var gPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a branch from origin.",
	Long:  `Will pull the currently checked out branch`,
	Run:   errorHandlerWrapper(handleGPull, 1),
}

var gRbiCmd = &cobra.Command{
	Use:   "rbi",
	Short: "Run an interactive rebase.",
	Long:  `Starts an interactive rebase, for the number of commits required.`,
	Run:   errorHandlerWrapper(handleGRebaseInteractively, 1),
}

var gPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes the branch to origin.",
	Long:  `Pushes the current branch to origin, using the current branches name as the target on the origin.`,
	Run:   errorHandlerWrapper(handleGPush, 1),
}

var gUndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Removes commits from the branch.",
	Long:  `This will (destructively) remove commits from the current branch. The number of commits to remove is defined by the 'number' option.`,
	Run:   errorHandlerWrapper(handleGUndo, 1),
}

var gLoglCmd = &cobra.Command{
	Use:   "logl",
	Short: "Shows the last X commits.",
	Long:  `Only shows the short commit sha and subject of each commit for the last X commits on the current branch.`,
	Run:   errorHandlerWrapper(handleGLogl, 1),
}

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Utility commands, mainly helper type commands for internal use.",
	RunE:  handleGroup,
}

var utilGenDocsCmd = &cobra.Command{
	Use:   "gen-docs",
	Short: "Generates markdown documentation for the CLI.",
	Run:   errorHandlerWrapper(handleUtilGenDocs, 1),
}

var sysCmd = &cobra.Command{
	Use:   "sys",
	Short: "Commands for handling the installation of this tool.",
	RunE:  handleGroup,
}

var sysCheckCmd = &cobra.Command{
	Use:          "check",
	Short:        "Checks for dependencies.",
	Long:         `Checks that the required system dependencies are installed and usable.`,
	Run:          errorHandlerWrapper(handleCheck, 1),
	SilenceUsage: true,
}

func errorHandlerWrapper(f runEHandlerFunc, errorExitCode int) runHandlerFunc {
	return func(cmd *cobra.Command, args []string) {
		err := f(cmd, args, appCfg)

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

func init() {
	cobra.OnInitialize(bootstrap)

	// Persistent flags
	rootCmd.PersistentFlags().String("log-level", logging.LogLevelNoneName, "Log level to use, one of: debug, info, warn, error, none. Defaults to none, as most errors are already surfaced anyway.")
	rootCmd.PersistentFlags().String("log-format", "text", "log format to use")

	cobra.CheckErr(viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("log-level")))
	cobra.CheckErr(viper.BindPFlag("logging.format", rootCmd.PersistentFlags().Lookup("log-format")))

	// Version
	versionCmd.Flags().Bool("short", false, "Show only the version, excluding commit and date information.")
	rootCmd.AddCommand(versionCmd)

	// Sys
	sysCmd.AddCommand(sysCheckCmd)
	rootCmd.AddCommand(sysCmd)

	// Git
	gCoCmd.Flags().BoolP("pull", "p", false, "Pulls the branch from origin after checking it out.")
	gCmd.AddCommand(gCoCmd)

	gCmd.AddCommand(gBranchesCmd)

	gRbiCmd.Flags().IntP("number", "n", 2, "The number of commits to include in the interactive rebase.")
	gCmd.AddCommand(gRbiCmd)

	gPushCmd.Flags().BoolP("force", "f", false, "Whether force push the branch.")
	gCmd.AddCommand(gPushCmd)

	gUndoCmd.Flags().BoolP("yes", "y", false, "Skip the confirmation prompt, and just delete them. #YOLO")
	gUndoCmd.Flags().IntP("number", "n", 1, "The number of commits to remove from the branch.")
	gCmd.AddCommand(gUndoCmd)

	gLoglCmd.Flags().IntP("number", "n", 10, "The number of commits to remove from the branch.")
	gCmd.AddCommand(gLoglCmd)

	gCmd.AddCommand(gPullCmd)

	rootCmd.AddCommand(gCmd)

	// Utils
	utilCmd.AddCommand(utilGenDocsCmd)

	rootCmd.AddCommand(utilCmd)
}

func buildConfig() *config.Config {
	if appCfg != nil {
		return appCfg
	}

	c := config.NewDefault()
	err := viper.Unmarshal(c)
	cobra.CheckErr(err)

	appCfg = c
	return appCfg
}

func bootstrap() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	currentDir, err := os.Getwd()
	// 	cobra.CheckErr(err)

	// 	viper.AddConfigPath(currentDir)
	// 	viper.SetConfigType("yaml")
	// 	viper.SetConfigName("config")
	// }

	// Tell viper to replace . in nested path with underscores
	// e.g. logging.level becomes LOGGING_LEVEL
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetEnvPrefix("orca")
	viper.AutomaticEnv()

	// err := viper.ReadInConfig()

	// cobra.CheckErr(err)

	cfg := buildConfig()

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
