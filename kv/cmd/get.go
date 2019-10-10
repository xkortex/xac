/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xkortex/xac/kv/util"
	"log"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Aliases: []string{"g"},

	Short: "Get a value from the store",
	Long: `Attempts to get a value from the store given the provided key`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic("/\\--/\\ Must have exactly one arguments (handling under construction)")
		}
		ns, _ := cmd.Flags().GetString("namespace")
		silent, _ := cmd.Flags().GetBool("silent")

		key := args[0]
		lookup_path := util.GetLookupPath(ns, key)
		util.Vprint(lookup_path)
		val, err := util.Read_value(lookup_path, key)
		if err != nil && !silent {
			log.Fatal(err)
		}
		// So bash is smart enough to strip whitespace so even though this adds
		// a newline, it still seems to work for loading and
		fmt.Println(val)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	//RootCmd.Flags().StringP("namespace", "n", "", "namespace of kv store")

}
