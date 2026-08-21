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
	viper.BindEnv("api_key", "BACKWORK_API_KEY")
	viper.BindEnv("base_url", "BACKWORK_BASE_URL")
	viper.BindEnv("output", "BACKWORK_OUTPUT")

	viper.ReadInConfig()
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
