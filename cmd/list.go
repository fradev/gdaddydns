package cmd

import (
	"github.com/fradev/gdaddydns/gdaddydns"
	"github.com/spf13/cobra"
)

var (
	fileDump string
	notable  bool
)
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "The 'list' subcommand list all the entries in the provided DNS .",
	Long: `The 'list' subcommand list all the entries in the provided DNS .
	It allows to filter the entries via 'type' arg and store the raw response (json) inside a file `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gdaddydns.ListEntries(ValidateDomainConfFile(domain), flagTypeEnum, fileDump, notable, goDaddyUrl); err != nil {
			er(err)
		}
	},
}

func init() {
	listCmd.PersistentFlags().StringVar(&domain, "domain", "", "Domain")
	listCmd.MarkPersistentFlagRequired("domain")
	listCmd.Flags().Var(&flagTypeEnum, "type", `Dns Type. Allowed  "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT"`)
	listCmd.Flags().StringVar(&fileDump, "file", "", "File to store the raw json (backup)")
	listCmd.Flags().BoolVar(&notable, "no-table", false, "No table output")

	rootCmd.AddCommand(listCmd)
}
