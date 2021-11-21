package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// add deploy Cmd to rootCmd
	rootCmd.AddCommand(deployCmd)
	// add sub Cmds of deploy to deployCmd
	deployCmd.AddCommand(deployWebCmd)
	deployCmd.AddCommand(deployAPICmd)
	deployCmd.AddCommand(deployDatabaseCmd)

}

// Note Run argument is not provided since we don’t want deploy command to do anything own.
// ‘deploy’ should be executed together with it’s sub-commands: web, api or database.
var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"dep", "depl"},
	Short:   "Deploy artifacts (web, api or database)",
	Long:    `This command can be used together with web, api or database sub-commands to deploy respective artifacts`,
}

var deployWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Deploy web artifacts",
	Long:  `This command can be used to deploy web artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'deployer deploy web' placeholder command")
	},
}

var deployAPICmd = &cobra.Command{
	Use:   "api",
	Short: "Deploy API artifacts",
	Long:  `This command can be used to deploy API artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'deployer deploy api' placeholder command")
	},
}

var deployDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Deploy database artifacts",
	Long:  `This command can be used to deploy database artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'deployer deploy database' placeholder command")
	},
}
