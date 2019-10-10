/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"github.com/xkortex/xac/kv/util"
	"log"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Aliases: []string{"del", "dlt"},

	Short: "Remove a key from the store",
	Long: `Deletes the key from the kv store`,
	Run: func(cmd *cobra.Command, args []string) {
		ns, _ := cmd.Flags().GetString("namespace")
		silent, _ := cmd.Flags().GetBool("silent")
		key := args[0]
		lookup_path := util.GetLookupPath(ns, key)
		util.Vprint(lookup_path)
		_, err := util.Pop_value(lookup_path, key)
		if err != nil && !silent {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
