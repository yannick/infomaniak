package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yannick/informaniak/internal/api"
)

var domainsShowCmd = &cobra.Command{
	Use:   "show <domain>",
	Short: "Show details for a domain",
	Args:  cobra.ExactArgs(1),
	RunE:  runDomainsShow,
}

func init() {
	domainsCmd.AddCommand(domainsShowCmd)
}

func runDomainsShow(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("token is required: set via --token, config file, or $INFOMANIAK_TOKEN")
	}

	client := api.NewClient(api.ClientConfig{Token: token})

	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	domain, err := client.ShowDomain(ctx, args[0])
	if err != nil {
		return fmt.Errorf("show domain: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Name:\t%s\n", domain.Name)
	fmt.Fprintf(w, "TLD:\t%s\n", domain.TLD)
	fmt.Fprintf(w, "Premium:\t%v\n", domain.IsPremium)
	fmt.Fprintf(w, "Created:\t%s\n", time.Unix(domain.CreatedAt, 0).Format("2006-01-02"))
	fmt.Fprintf(w, "Expires:\t%s\n", time.Unix(domain.ExpiresAt, 0).Format("2006-01-02"))
	fmt.Fprintf(w, "DNS Anycast:\t%v\n", domain.Options.DNSAnycast)
	fmt.Fprintf(w, "DNSSEC:\t%v\n", domain.Options.DNSSEC)
	fmt.Fprintf(w, "Domain Privacy:\t%v\n", domain.Options.DomainPrivacy)
	return w.Flush()
}
