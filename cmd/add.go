package cmd

import (
	"github.com/fradev/gdaddydns/gdaddydns"
	"github.com/spf13/cobra"
)

var ttl int
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "The 'add' subcommand add a new entry to the domain dns.",
	Long: `The 'add' subcommand add a new entry to the domain dns.
	Adding the dns entry 'name'  to the domain passed as argument. 
	`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

		if err := gdaddydns.AddEntry(ValidateDomainConfFile(domain), data, name, flagTypeEnum, ttl, goDaddyUrl); err != nil {
			er(err)
		}

	},
}

func init() {
	addCmd.PersistentFlags().StringVar(&domain, "domain", "", "Domain (required)")
	addCmd.MarkPersistentFlagRequired("domain")
	addCmd.PersistentFlags().Var(&flagTypeEnum, "type", `Entry Type. Allowed  "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT" (required)`)
	addCmd.MarkPersistentFlagRequired("type")
	addCmd.PersistentFlags().StringVar(&data, "data", "", "DNS Data Ip/FDQN to point (required)")
	addCmd.MarkPersistentFlagRequired("data")
	addCmd.PersistentFlags().StringVar(&name, "name", "", "Hostname (required)")
	addCmd.MarkPersistentFlagRequired("name")
	addCmd.PersistentFlags().IntVar(&ttl, "ttl", 600, "TTL of the record")
	rootCmd.AddCommand(addCmd)
}
