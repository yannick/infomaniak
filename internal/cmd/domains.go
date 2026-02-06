package cmd

import (
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Manage Infomaniak domains",
}

func init() {
	rootCmd.AddCommand(domainsCmd)
}
