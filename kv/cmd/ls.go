/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xkortex/xac/kv/util"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Aliases: []string{"l", "list"},

	Short: "List values",
	Long: `List values present in the store`,
	Run: func(cmd *cobra.Command, args []string) {
		ns, _ := cmd.Flags().GetString("namespace")

		lookup_path := util.GetLookupPath(ns, "")
		util.Vprint(lookup_path)
		util.List_all(lookup_path)
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
