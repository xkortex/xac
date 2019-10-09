/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"fmt"
	"github.com/xkortex/xac/kv/util"
	"log"

	"github.com/spf13/cobra"
)

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop a value from the kv",
	Long: `Return a value according to the key. Remove that value from the store`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic("/\\--/\\ Must have exactly one argument (handling under construction)")
		}
		ns, _ := cmd.Flags().GetString("namespace")
		silent, _ := cmd.Flags().GetBool("silent")

		key := args[0]
		lookup_path := util.GetLookupPath(ns, key)
		util.Vprint(lookup_path)
		val, err := util.Pop_value(lookup_path, key)
		if err != nil  && !silent {
			log.Fatal(err)
		}
		fmt.Println(val)
	},
}

func init() {
	RootCmd.AddCommand(popCmd)
	//RootCmd.Flags().StringP("namespace", "n", "", "namespace of kv store")

}
