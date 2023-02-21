package cmd

import (
	"github.com/fradev/gdaddydns/gdaddydns"
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "The 'domains' subcommand list all the domains inside the config file passed.",
	Long:  `The 'domains' subcommand list all the domains inside the config file passed.`,
	Run: func(cmd *cobra.Command, args []string) {
		gdaddydns.PrintDomains(entries.Domains)

	},
}

func init() {
	rootCmd.AddCommand(domainsCmd)
}
