package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/panoptescloud/orca/internal/hostsys"
	"github.com/panoptescloud/orca/internal/logging"
	"github.com/panoptescloud/orca/internal/workspaces"
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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Commands for managing configuration",
	RunE:  handleGroup,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the config",
	Run:   errorHandlerWrapper(handleConfigShow, 1),
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the path to config being used.",
	Run:   errorHandlerWrapper(handleConfigPath, 1),
}

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "Commands related to managing workspaces.",
	RunE:  handleGroup,
}

var wsSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to another workspace",
	Run:   errorHandlerWrapper(handleWsSwitch, 1),
}

var wsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises a workspace.",
	Long: `Initialises a new workspace with orca. Either a local directory or a git repository can be supplied.
When using a local directory, will locate an orca workspace config and clone the relevant projects.
When using a git url, will clone the given repository, and then clone any other required repositories based on an orca workspace file within the repo.`,
	Run: errorHandlerWrapper(handleWsInit, 1),
}

var wsLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all available workspaces.",
	Run:   errorHandlerWrapper(handleWsLs, 1),
}

var wsCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones all the projects required for this workspace.",
	Long: `Based on the workspace config will clone each repository required by the workspace.
This can be run at any time to clone any projects that have not already been cloned.
The project option allows you to clone only a specific project.`,
	Run: errorHandlerWrapper(handleWsClone, 1),
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

var sysInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a tool system tool that is needed.",
	Long: fmt.Sprintf(`Tools are installed 'locally' rather than globally, they will be stored within %s.

The first argument must be one of: %s`, getToolsDir(), hostsys.AllAvailableToolsCsv()),
	Run:          errorHandlerWrapper(handleSysInstall, 1),
	SilenceUsage: true,
}

var sysSelfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Updates this tool.",
	Long:  `By default will update to the latest available version. A specific version can be specified if a specific version is required.`,
	Run:   errorHandlerWrapper(handleSysSelfUpdate, 1),
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Starts the workspace or project.",
	Long:  `...TBD...`,
	Run:   errorHandlerWrapper(handleUp, 1),
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

	// Version
	versionCmd.Flags().Bool("short", false, "Show only the version, excluding commit and date information.")
	rootCmd.AddCommand(versionCmd)

	// Sys
	sysSelfUpdateCmd.Flags().String("to", "", "The version you wish to switch to. If left blank will download latest avaialable")
	sysCmd.AddCommand(sysSelfUpdateCmd)

	sysCmd.AddCommand(sysCheckCmd)
	sysCmd.AddCommand(sysInstallCmd)
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

	// config
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)

	// workspaces
	wsCmd.AddCommand(wsSwitchCmd)

	wsInitCmd.Flags().StringP("source", "s", getWorkingDir(), "The directory of the workspace configuration.")
	wsInitCmd.Flags().StringP("target", "t", getWorkingDirParent(), "The directory to store the workspace projects.")
	wsInitCmd.Flags().StringP("config", "c", workspaces.DefaultWorkspaceFileName, "The name of the workspace config file within the source.")
	wsCmd.AddCommand(wsInitCmd)

	wsCmd.AddCommand(wsLsCmd)

	wsCloneCmd.Flags().StringP("target", "t", "", `The directory in which to clone the project(s). 
If multiple projects are being cloned, then it will place them in {target}/{repo name}.
If a single project is being clone then it will be cloned into {target}.`)
	addWorkspaceOption(wsCloneCmd)
	addProjectOption(wsCloneCmd)
	wsCloneCmd.MarkFlagRequired("workspace")
	wsCmd.AddCommand(wsCloneCmd)

	rootCmd.AddCommand(wsCmd)

	// up
	addWorkspaceOption(upCmd)
	addProjectOption(upCmd)

	rootCmd.AddCommand(upCmd)
}

func addWorkspaceOption(cmd *cobra.Command) {
	cmd.Flags().StringP("workspace", "w", "", "The name of the workspace to run this command for.")
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
