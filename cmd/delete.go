package cmd

import (
	"github.com/fradev/gdaddydns/gdaddydns"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "The 'del' subcommand delete a specific entry in the DNS.",
	Long:  `The 'del' subcommand delete a specific entry in the DNS.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := gdaddydns.DelEntry(ValidateDomainConfFile(domain), name, flagTypeEnum, goDaddyUrl); err != nil {
			er(err)
		}

	},
}

func init() {
	delCmd.PersistentFlags().StringVar(&domain, "domain", "", "Domain (required)")
	delCmd.MarkPersistentFlagRequired("domain")
	delCmd.PersistentFlags().StringVar(&name, "name", "", "Hostname (required)")
	delCmd.MarkPersistentFlagRequired("name")
	delCmd.PersistentFlags().Var(&flagTypeEnum, "type", `Entry Type. Allowed  "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT" (required)`)
	delCmd.MarkPersistentFlagRequired("type")
	rootCmd.AddCommand(delCmd)
}
