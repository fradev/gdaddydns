package gdaddydns

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/table"
)

const GO_DADDY_API_SERVER = "https://api.godaddy.com"
const GO_DADDY_DOMAIN = "/v1/domains/{domain}/records/"
const GO_DADDY_DOMAIN_TYPE = GO_DADDY_DOMAIN + "{type}/"
const GO_DADDY_DOMAIN_TYPE_NAME = GO_DADDY_DOMAIN_TYPE + "{name}/"

// Domain type compose of a Name (example.com) and api key and secret
type Domain struct {
	Name       string
	Api_Key    string
	Api_Secret string
}

// Array of Domains
type Domains struct {
	Domains []Domain `json:"array"`
}

// Enum for the type of the DNS entry
type dnsType string

const (
	A     dnsType = "A"
	AAAA  dnsType = "AAAA"
	CNAME dnsType = "CNAME"
	MX    dnsType = "MX"
	NS    dnsType = "NS"
	SOA   dnsType = "SOA"
	SRV   dnsType = "SRV"
	TXT   dnsType = "TXT"
	NIL   dnsType = "NIL"
)

// to string method for the enum. Returns empty on NIL for
// the default value in the command line
func (e *dnsType) String() string {
	if *e == NIL {
		return ""
	}
	return string(*e)
}

// Setter from string for the enum type
func (e *dnsType) Set(v string) error {
	switch v {
	case "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT":
		*e = dnsType(v)
		return nil
	default:
		return errors.New(`must be one of  "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT"`)
	}
}

// Get type
func (e *dnsType) Type() string {
	return "dnsType"
}

// print the error message using color.Red
func PrintErrorMsg(msg string) {
	color.Red(msg)
}

// print an info message using color.Blue
func PrintInfo(msg string) {
	color.Blue(msg)
}

// disabe colors
func SetNoColor() {
	color.NoColor = true
}

// DNS entry struct returned by the godaddy api
// Json is the following
// [
//
//	{
//	 "data": "string",
//	 "name": "string",
//	 "port": 65535,
//	 "priority": 0,
//	 "protocol": "string",
//	 "service": "string",
//	 "ttl": 0,
//	 "type": "A",
//	 "weight": 0
//	}
//
// ]
type DNSEntry struct {
	Data     string  `json:"data"`
	Name     string  `json:"name"`
	TTL      int     `json:"ttl"`
	DNSType  dnsType `json:"type"`
	Weight   int     `json:"weight"`
	Service  string  `json:"service"`
	Priority int     `json:"priority"`
	Protocol string  `json:"protocol"`
	Port     int     `json:"port"`
}

// Get call for listing entries WITH DNS type  filtering
// https://developer.godaddy.com/doc/endpoint/domains#/v1/recordGet
func listEntriesWithType(req *resty.Request, domain string, dnsType string) (*resty.Response, error) {
	return req.SetPathParams(map[string]string{
		"domain": strings.ToLower(domain),
		"type":   dnsType,
	}).Get(GO_DADDY_DOMAIN_TYPE)
}

// Get call for listing entries WITHOUT DNS type filtering
// https://developer.godaddy.com/doc/endpoint/domains#/v1/recordGet
func listEntriesWithoutType(req *resty.Request, domain string) (*resty.Response, error) {
	return req.
		SetPathParam("domain", strings.ToLower(domain)).
		Get(GO_DADDY_DOMAIN)
}

// Get call for listing entries in the domain
// https://developer.godaddy.com/doc/endpoint/domains#/v1/recordGet
func ListEntries(entry Domain, dnsType dnsType, filedump string, notable bool, goDaddyURL string) error {
	var resp *resty.Response
	var restErr error
	client := resty.New()
	req := client.SetBaseURL(goDaddyURL).R().EnableTrace().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("sso-key %s:%s", entry.Api_Key, entry.Api_Secret))
	PrintInfo(fmt.Sprintf("Using Access Key %s and Secret %s for %s", entry.Api_Key, entry.Api_Secret, strings.ToLower(entry.Name)))
	if dnsType != NIL {
		PrintInfo(fmt.Sprintf("Base URL  is %s and path is %s ", goDaddyURL, GO_DADDY_DOMAIN))
		resp, restErr = listEntriesWithType(req, strings.ToLower(entry.Name), dnsType.String())
	} else {
		PrintInfo(fmt.Sprintf("Base URL  is %s and path is %s and type %s ", goDaddyURL, GO_DADDY_DOMAIN_TYPE, dnsType))
		resp, restErr = listEntriesWithoutType(req, strings.ToLower(entry.Name))

	}
	if restErr != nil {
		printResponseError(resp, restErr)
		return restErr
	}
	if len(filedump) > 0 {
		if err := os.WriteFile(filedump, resp.Body(), 0644); err != nil {
			return err
		} else {
			PrintInfo(fmt.Sprintf("Dumped response in  %s ", filedump))
		}
	} else {
		var dnsEntries []DNSEntry
		if err := json.Unmarshal(resp.Body(), &dnsEntries); err != nil {
			printResponseError(resp, restErr)
			return err
		}
		color.Set(color.FgGreen, color.Bold)
		defer color.Unset()
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		if notable {
			t.Style().Options.DrawBorder = false
			t.Style().Options.SeparateColumns = false
			t.Style().Options.SeparateFooter = false
			t.Style().Options.SeparateHeader = false
			t.Style().Options.SeparateRows = false
		} else {
			t.AppendHeader(table.Row{"#", "Data", "Name", "TTL", "Type"})
		}
		for i, e := range dnsEntries {
			if notable {
				t.AppendRow(table.Row{e.Name, e.Data, e.DNSType.String(), e.TTL})
			} else {
				t.AppendRow(table.Row{(i + 1), e.Name, e.Data, e.DNSType.String(), e.TTL})
				t.AppendSeparator()
			}
		}
		t.Render()
	}
	return nil
}

// Add an entry to the DNS
func AddEntry(entry Domain, data string, name string, dnsType dnsType, ttl int, goDaddyUrl string) error {
	PrintInfo(fmt.Sprintf("Pointing %s to %s (%s) in %s", name, data, dnsType.String(), strings.ToLower(entry.Name)))
	client := resty.New()
	dnsEntry := DNSEntry{data, name, ttl, dnsType, 0, "string", 0, "string", 65535}
	dnsEntries := []DNSEntry{dnsEntry}
	resp, err := client.SetBaseURL(goDaddyUrl).R().EnableTrace().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("sso-key %s:%s", entry.Api_Key, entry.Api_Secret)).
		SetPathParams(map[string]string{
			"domain": strings.ToLower(entry.Name),
			"type":   dnsType.String(),
			"name":   name,
		}).SetBody(dnsEntries).Put(GO_DADDY_DOMAIN_TYPE_NAME)
	if err != nil {
		printResponseError(resp, err)
		return err
	}
	if resp.StatusCode() == 200 {
		PrintInfo(fmt.Sprintf("Added %s OK", name))
		return nil

	} else {
		printResponseError(resp, nil)
		return fmt.Errorf("error adding %s", name)
	}
}

// Delete an Entry to the DNS
func DelEntry(entry Domain, name string, dnsType dnsType, goDaddyURL string) error {
	PrintInfo(fmt.Sprintf("Deleting %s (%s) in %s", name, dnsType.String(), strings.ToLower(entry.Name)))
	client := resty.New()
	resp, err := client.SetBaseURL(goDaddyURL).R().EnableTrace().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("sso-key %s:%s", entry.Api_Key, entry.Api_Secret)).
		SetPathParams(map[string]string{
			"domain": strings.ToLower(entry.Name),
			"type":   dnsType.String(),
			"name":   name,
		}).Delete(GO_DADDY_DOMAIN_TYPE_NAME)
	if err != nil {
		printResponseError(resp, err)
		return err
	}
	if resp.StatusCode() == 200 || resp.StatusCode() == 204 {
		PrintInfo(fmt.Sprintf("Deleted %s OK", name))
		return nil

	} else {
		printResponseError(resp, nil)
		return fmt.Errorf("error deleting %s", name)
	}
}

// Print domains configured in the config file
func PrintDomains(domains []Domain) {
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Domain Name"})
	for i, s := range domains {
		t.AppendRow(table.Row{(i + 1), s.Name})
		t.AppendSeparator()
	}
	t.Render()
}

// print the error response completely
func printResponseError(resp *resty.Response, err error) {
	PrintErrorMsg("Response Info:")
	PrintErrorMsg(fmt.Sprintf("  Error      :%v	", err))
	PrintErrorMsg(fmt.Sprintf("  Status Code:%d", resp.StatusCode()))
	PrintErrorMsg(fmt.Sprintf("  Status     :%s", resp.Status()))
	PrintErrorMsg(fmt.Sprintf("  Proto      :%s", resp.Proto()))
	PrintErrorMsg(fmt.Sprintf("  Time       :%q", resp.Time()))
	PrintErrorMsg(fmt.Sprintf("  Received At:%q", resp.ReceivedAt()))
	PrintErrorMsg(fmt.Sprintf("  Body       :%s\n", resp.Body()))
	PrintErrorMsg("")
}
