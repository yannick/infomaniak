package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yannick/informaniak/internal/api"
)

var domainsUpdateNSCmd = &cobra.Command{
	Use:   "update-ns <domain>",
	Short: "Update nameservers for a domain",
	Args:  cobra.ExactArgs(1),
	RunE:  runDomainsUpdateNS,
}

func init() {
	domainsUpdateNSCmd.Flags().StringSlice("nameservers", nil, "comma-separated list of nameservers")
	domainsUpdateNSCmd.Flags().Bool("verify", false, "verify nameserver availability before applying")
	_ = domainsUpdateNSCmd.MarkFlagRequired("nameservers")

	domainsCmd.AddCommand(domainsUpdateNSCmd)
}

func runDomainsUpdateNS(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("token is required: set via --token, config file, or $INFOMANIAK_TOKEN")
	}

	nameservers, err := cmd.Flags().GetStringSlice("nameservers")
	if err != nil {
		return fmt.Errorf("parse nameservers flag: %w", err)
	}

	verify, err := cmd.Flags().GetBool("verify")
	if err != nil {
		return fmt.Errorf("parse verify flag: %w", err)
	}

	client := api.NewClient(api.ClientConfig{Token: token})

	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	input := api.UpdateNameserversInput{
		Nameservers:          nameservers,
		VerifyNSAvailability: verify,
	}

	if err := client.UpdateNameservers(ctx, args[0], input); err != nil {
		return fmt.Errorf("update nameservers: %w", err)
	}

	jsonOut, _ := cmd.Flags().GetBool("json")

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]string{
			"domain": args[0],
			"status": "updated",
		})
	}

	fmt.Printf("Nameservers for %s updated successfully.\n", args[0])
	return nil
}
