package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/tylergibbs1/backwork-cli/pkg/client"
)

var drugsCmd = &cobra.Command{
	Use:   "drugs",
	Short: "Drug formulary commands",
	Long:  "Search commercial pharmacy-benefit formulary evidence",
}

var drugsFormularyCmd = &cobra.Command{
	Use:   "formulary [query]",
	Short: "Search drug formulary evidence",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		payer, _ := cmd.Flags().GetString("payer")
		limit, _ := cmd.Flags().GetInt("limit")
		values := url.Values{}
		values.Set("q", args[0])
		values.Set("payer", payer)
		values.Set("limit", fmt.Sprintf("%d", limit))

		var result map[string]interface{}
		if err := c.Get("/drugs/formulary?"+values.Encode(), &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printDrugFormularyResults(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(drugsCmd)
	drugsCmd.AddCommand(drugsFormularyCmd)

	drugsFormularyCmd.Flags().StringP("payer", "p", "all", "Payer/PBM source (all, cvs_caremark, express_scripts, uhc)")
	drugsFormularyCmd.Flags().IntP("limit", "l", 25, "Maximum results per source")
}

func printDrugFormularyResults(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No formulary evidence found")
		return
	}

	for _, item := range data {
		row := item.(map[string]interface{})
		fmt.Printf("%v (%v)\n", row["drug_name"], row["source"])
		fmt.Printf("  Payer: %v\n", row["payer_name"])
		if status, ok := row["coverage_status"].(string); ok && status != "" {
			fmt.Printf("  Coverage: %s\n", status)
		}
		if tier, ok := row["tier"].(string); ok && tier != "" {
			fmt.Printf("  Tier: %s\n", tier)
		}
		if matched, ok := row["matched_text"].(string); ok && matched != "" {
			fmt.Printf("  Match: %s\n", matched)
		}
		fmt.Println("---")
	}
}
