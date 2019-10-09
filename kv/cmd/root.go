/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"github.com/xkortex/xac/kv/util"
	"log"

	"github.com/spf13/cobra"
)

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "kv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.Vprint("root called")
		util.Vprint(args)
		ns, _ := cmd.PersistentFlags().GetString("namespace")
		util.Vprint(ns)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {
	//RootCmd.AddCommand(RootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RootCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().StringP("namespace", "n", "", "namespace of kv store")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RootCmd.PersistentFlags().BoolP("silent", "s", false, "Suppress errors")
	RootCmd.PersistentFlags().BoolP("stdin", "-", false, "Read from standard in")
	RootCmd.Flags().BoolP("verbose", "v", false, "Verbose tracing (in progress)")
}
