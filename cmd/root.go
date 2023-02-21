package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fradev/gdaddydns/gdaddydns"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	nocolor    bool
	entries    gdaddydns.Domains
	msg        string
	exx        error
	version    = "0.0.1"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gdaddydns",
	Version: version,
	Short:   "Simple Utility to manipulate Go Daddy DNS via API",
	Long: `Simple Utility to manipulate Go Daddy DNS via API.
Using access key and secret pair, list, add, removes, edit the DNS
entries.
It reads the configuration from a json file (default  ~/.gdaddydns.json)
or passed via args. The format of the file must be the following:
'
{
  "domains": [
    {"name": "example.com", "api_key": "EXAMPLE_KEY", "api_secret": "EXAMPLE_SECRET"},
    {"name": "me.com", "api_key": "ME_KEY", "api_secret": "ME_SECRET"},
    {"name": "xxxx.net", "api_key": "XXXX_KEY", "api_secret": "XXXX_SECRET"}
   ]
}
'
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if nocolor {
			gdaddydns.SetNoColor()
		}
		if msg != "" {
			gdaddydns.PrintInfo(msg)
		}
		if exx != nil {
			er(exx)
		}
		if len(entries.Domains) == 0 {
			er(fmt.Errorf("configuration file is empty or bad formatted"))
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.gdaddydns.json)")
	rootCmd.PersistentFlags().BoolVar(&nocolor, "no-color", false, "No color output")
	rootCmd.PersistentFlags().StringVar(&goDaddyUrl, "godaddy-url", gdaddydns.GO_DADDY_API_SERVER, "GoDaddy API base URI")

}

func er(e error) {
	gdaddydns.PrintErrorMsg(e.Error())
	os.Exit(1)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(".gdaddydns.json")
		viper.SetConfigType("json")
		viper.AddConfigPath("$HOME")
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			exx = err
		} else {
			exx = fmt.Errorf("error while reading %s, check the file syntax", viper.ConfigFileUsed())
		}
	} else {
		msg = fmt.Sprintln("Using configuration file ", viper.ConfigFileUsed())
		if err = viper.Unmarshal(&entries); err != nil {
			exx = fmt.Errorf("error while unmarshal %s, check the syntax", viper.ConfigFileUsed())
		}

	}
}
func ValidateDomainConfFile(domain string) gdaddydns.Domain {
	var found bool = false
	var list []string
	var ret gdaddydns.Domain
	for _, s := range entries.Domains {
		if strings.EqualFold(s.Name, domain) {
			found = true
			ret = s
		}
		list = append(list, s.Name)
	}
	if !found {
		er(fmt.Errorf("%s domain not found the list  %s", domain, list))
	}
	return ret

}
