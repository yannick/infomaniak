package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

// SetVersion sets the application version from the build-time value.
func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:     "infomaniak",
	Short:   "Manage Infomaniak domains",
	Long:    "infomaniak is a CLI tool to manage domains and nameservers via the Infomaniak API.",
	Version: version,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "config file (default $HOME/.infomaniak.yaml)")
	rootCmd.PersistentFlags().String("token", "", "Infomaniak API token")
	rootCmd.PersistentFlags().Bool("json", false, "output as JSON")
	rootCmd.PersistentFlags().Bool("simple", false, "output simplified plain text")
	rootCmd.MarkFlagsMutuallyExclusive("json", "simple")

	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
}

func initConfig() {
	if cfgFile := rootCmd.PersistentFlags().Lookup("config").Value.String(); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("find home directory", "error", err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".infomaniak")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("INFOMANIAK")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			slog.Warn("reading config file", "error", err)
		}
	}
}

// Execute runs the root command.
func Execute() error {
	rootCmd.Version = version
	return rootCmd.Execute()
}
