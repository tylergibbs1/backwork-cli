package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tylergibbs1/backwork-cli/pkg/client"
)

var batchCmd = &cobra.Command{
	Use:   "batch [codes...]",
	Short: "Batch lookup multiple medical codes",
	Long:  "Look up multiple medical codes (CPT, HCPCS, ICD-10, NDC) in a single request",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		reqBody := map[string]interface{}{
			"codes": args,
		}

		system, _ := cmd.Flags().GetString("system")
		if system != "" {
			reqBody["code_system"] = system
		}

		include, _ := cmd.Flags().GetStringSlice("include")
		if len(include) > 0 {
			reqBody["include"] = strings.Join(include, ",")
		}

		var result map[string]interface{}
		if err := c.Post("/codes/batch", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printBatchResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringP("system", "s", "", "Code system (CPT, HCPCS, ICD-10, NDC)")
	batchCmd.Flags().StringSliceP("include", "i", []string{}, "Include additional data (rvu, policies)")
}

func printBatchResult(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No results found")
		return
	}

	fmt.Printf("%-12s %-10s %-8s %s\n", "CODE", "SYSTEM", "FOUND", "DESCRIPTION")
	fmt.Printf("%-12s %-10s %-8s %s\n", "----", "------", "-----", "-----------")

	for _, item := range data {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		code := fmt.Sprintf("%v", entry["code"])
		system := fmt.Sprintf("%v", entry["code_system"])
		found := fmt.Sprintf("%v", entry["found"])

		description := ""
		if desc, ok := entry["description"].(string); ok {
			description = desc
			if len(description) > 60 {
				description = description[:57] + "..."
			}
		}

		fmt.Printf("%-12s %-10s %-8s %s\n", code, system, found, description)
	}
}
