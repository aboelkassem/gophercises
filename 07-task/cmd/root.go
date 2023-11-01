package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here

	// 	fmt.Println(cmd.Short)

	// 	fmt.Println()

	// 	// print usage
	// 	if err := cmd.Usage(); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(1) // error happened
	// 	}
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitf("%v", err)
	}
}

// just to override and hide later
func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generate the autocompletion script for the specified shell",
	}
}

func init() {
	// hide help command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	completion := completionCommand()

	// mark completion hidden
	completion.Hidden = true
	rootCmd.AddCommand(completion)
}

func exitf(format string, a ...any) {
	fmt.Printf(format, a)
	os.Exit(1) // error happened
}
