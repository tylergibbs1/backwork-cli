package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	apiKey    string
	baseURL   string
	output    string
	buildTime string
)

var rootCmd = &cobra.Command{
	Use:   "backwork",
	Short: "Backwork CLI - Medicare coverage policies and prior authorization",
	Long: `Backwork CLI provides access to Medicare coverage policies, prior authorization
requirements, and medical code lookups from the command line.

Get your API key from: https://backworkhealth.com/dashboard`,
}

func SetBuildInfo(version, builtAt string) {
	rootCmd.Version = version
	buildTime = builtAt
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetVersionTemplate("backwork {{.Version}}\n")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.backwork.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Backwork API key (or set BACKWORK_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "https://backworkhealth.com/api/v1", "API base URL")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format (table, json, yaml)")

	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".backwork")
	}

	viper.SetEnvPrefix("BACKWORK")
	viper.AutomaticEnv()

	// Anyone who configured the CLI before the Verity -> Backwork rename still
	// exports VERITY_*. viper.BindEnv (v1.21.0) tries the names in order and
	// takes the first one that is set, so the new name wins and the old one
	// keeps working instead of the CLI behaving as if it were unconfigured.
	// Remove these once the old variables are out of everyone's shell profiles.
	viper.BindEnv("api_key", "BACKWORK_API_KEY", "VERITY_API_KEY")
	viper.BindEnv("base_url", "BACKWORK_BASE_URL", "VERITY_BASE_URL")
	viper.BindEnv("output", "BACKWORK_OUTPUT", "VERITY_OUTPUT")

	if err := viper.ReadInConfig(); err != nil {
		if _, notFound := err.(viper.ConfigFileNotFoundError); notFound && cfgFile == "" {
			// Same reason as the env vars above: users have ~/.verity.yaml on
			// disk today. Silently ignoring it would leave the CLI looking
			// configured-but-empty. Drop this once those files are gone.
			viper.SetConfigName(".verity")
			viper.ReadInConfig()
		}
	}
}

func getAPIKey() string {
	key := viper.GetString("api_key")
	if key == "" {
		fmt.Fprintln(os.Stderr, "Error: API key is required. Set BACKWORK_API_KEY or use --api-key flag")
		os.Exit(1)
	}
	return key
}

func getBaseURL() string {
	return viper.GetString("base_url")
}

func getOutput() string {
	return viper.GetString("output")
}
