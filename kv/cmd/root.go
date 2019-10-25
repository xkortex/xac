/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"github.com/Wessie/appdirs"
	"github.com/spf13/viper"
	"github.com/xkortex/xac/kv/util"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	developer     string
	defaultCfgDir string
)
const defaultCfgName = "kv.yml"

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "kv",
	Short: "Utility for getting and setting key-value pairs",
	Long: `Does what it says on the tin. Bare-bone, no-nonsense kv store. 
Keys are stored as paths. 
Examples:
    $ kv foo=bar                  # Set foo to bar
    $ echo spam | kv foo          # set foo to spam
    $ kv foo                      # Get value of foo
    spam`,
	Run: func(cmd *cobra.Command, args []string) {
		util.Vprint("root called")
		util.Vprint(args)
		ns, _ := cmd.PersistentFlags().GetString("namespace")
		util.Vprint(ns)
		//if err := cmd.Usage(); err != nil {
		//	log.Fatalf("Error executing root command: %v", err)
		//}
		log.Fatal(cmd.SilenceErrors, cmd.SilenceUsage)

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
	cobra.OnInitialize(initConfig)
	defaultCfgDir = appdirs.UserConfigDir("kv", "", "", false)
	defaultCfgFile := filepath.Join(defaultCfgDir, "config.yml")
	//RootCmd.AddCommand(RootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RootCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		defaultCfgFile,
		"config file, based in UserConfigDir", )

	RootCmd.PersistentFlags().StringP("namespace", "n", "", "namespace of kv store")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RootCmd.PersistentFlags().BoolP("silent", "s", false, "Suppress errors")
	RootCmd.PersistentFlags().BoolP("stdin", "-", false, "Read from standard in")
	RootCmd.Flags().BoolP("verbose", "v", false, "Verbose tracing (in progress)")
	RootCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(defaultCfgDir)
		viper.SetConfigName(defaultCfgName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// If we just didn't find it, that's fine. Otherwise, we probably want
		// to know if, e.g., the file was corrupt.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// todo: make these debug logs
			log.Printf("No config read: %v", err)
		}
	} else {
		log.Printf("Using config file %q", viper.ConfigFileUsed())
	}
}
