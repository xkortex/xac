/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xkortex/xac/kv/util"
)

// dirnameCmd represents the dirname command
var pathCmd = &cobra.Command{
	Use:   "path",
	Aliases: []string{"p"},

	Short: "Get the directory name of the kv store",
	Long: `Returns the dirname of the kv store or a namespace root dir.
If no key is given, the dirname of the kv store is returned. 
If a key is given, the path to the keyfile is returned. 
It is based on AppDirs UserDataDir. In Linux, this is $XDG_DATA_HOME`,
	Run: func(cmd *cobra.Command, args []string) {
		key := ""
		if len(args) == 1 {
			key = args[0]
		} else if len(args) > 1 {
			panic("/\\--/\\ Must have at most one argument (handling under construction)")
		}
		ns, _ := cmd.Flags().GetString("namespace")

		lookup_path := util.GetLookupPath(ns, key)
		fmt.Println(lookup_path)

	},
}

var appdirCmd = &cobra.Command{
	Use:   "appdir",
	Short: "Get the base application directory name of the kv store",
	Long: `Returns the dirname where the kv store is rooted. 
It is based on AppDirs UserDataDir. In Linux, this is $XDG_DATA_HOME`,
	Run: func(cmd *cobra.Command, args []string) {
		lookup_path := util.GetLookupPath("", "")
		fmt.Println(lookup_path)

	},
}

func init() {
	RootCmd.AddCommand(pathCmd)
	RootCmd.AddCommand(appdirCmd)
}
