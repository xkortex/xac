/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xkortex/xac/dug/dug"
	"github.com/xkortex/xac/kv/util"
	"os"
)

// RootCmd represents the root command
var arpScanCmd = &cobra.Command{
	Use:     "arpscan",
	Short:   "Dug: a better dig",
	Aliases: []string{"as"},

	Long: `arp scan on an interface 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.JSONFormatter{DisableTimestamp:true})
		log.SetOutput(os.Stdout)

		util.Vprint("root called")
		util.Vprint(args)
		all, _ := cmd.PersistentFlags().GetBool("all")

		timeout, _ := cmd.PersistentFlags().GetFloat64("timeout")
		delay, _ := cmd.PersistentFlags().GetFloat64("delay")
		util.Vprint("Timeout: ", timeout)
		util.Vprint("Delay: ", delay)
		if all {
			arpResults, err := dug.ScanAllInterfaces(timeout, delay)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(arpResults)
			return
		}

		interfaceName := args[0]

		arpResults, err := dug.ScanInterface(interfaceName, timeout, delay)
		if err != nil {
			log.Fatal(err)
		}
		util.Vprint(arpResults) // in prog
		return

	},
}

func init() {
	RootCmd.AddCommand(arpScanCmd)
	arpScanCmd.PersistentFlags().BoolP("all", "a", false, "Scan all interfaces")
	arpScanCmd.PersistentFlags().Float64P("timeout", "t", 1.0, "Timeout in seconds")
	arpScanCmd.PersistentFlags().Float64P("delay", "d", 0.0, "Delay between sends in seconds")

}
