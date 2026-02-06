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

var domainsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains for an account",
	RunE:  runDomainsList,
}

func init() {
	domainsCmd.AddCommand(domainsListCmd)
}

func runDomainsList(cmd *cobra.Command, _ []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("token is required: set via --token, config file, or $INFOMANIAK_TOKEN")
	}

	accountID := viper.GetString("account_id")
	if accountID == "" {
		return fmt.Errorf("account-id is required: set via --account-id, config file, or $INFOMANIAK_ACCOUNT_ID")
	}

	client := api.NewClient(api.ClientConfig{Token: token})

	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	domains, err := client.ListDomains(ctx, accountID)
	if err != nil {
		return fmt.Errorf("list domains: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTLD\tEXPIRES")
	for _, d := range domains {
		expires := time.Unix(d.ExpiresAt, 0).Format("2006-01-02")
		fmt.Fprintf(w, "%s\t%s\t%s\n", d.Name, d.TLD, expires)
	}
	return w.Flush()
}
