package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tylergibbs1/backwork-cli/pkg/client"
)

var complianceCmd = &cobra.Command{
	Use:   "compliance",
	Short: "Compliance acknowledgment commands",
	Long:  "List, acknowledge, and summarize policy changes for compliance workflows",
}

var complianceUnreviewedCmd = &cobra.Command{
	Use:   "unreviewed",
	Short: "List unreviewed policy changes",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		values := url.Values{}
		limit, _ := cmd.Flags().GetInt("limit")
		values.Set("limit", fmt.Sprintf("%d", limit))

		changeType, _ := cmd.Flags().GetString("change-type")
		if changeType != "" {
			values.Set("change_type", changeType)
		}

		cursor, _ := cmd.Flags().GetString("cursor")
		if cursor != "" {
			values.Set("cursor", cursor)
		}

		var result map[string]interface{}
		if err := c.Get("/compliance/unreviewed?"+values.Encode(), &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		printResult(result, printUnreviewedChanges)
	},
}

var complianceAckCmd = &cobra.Command{
	Use:   "ack [diff-id]",
	Short: "Acknowledge a policy change",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		diffID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: diff-id must be an integer\n")
			return
		}

		reqBody := map[string]interface{}{
			"diff_id": diffID,
		}

		notes, _ := cmd.Flags().GetString("notes")
		if notes != "" {
			reqBody["notes"] = notes
		}

		var result map[string]interface{}
		if err := c.Post("/compliance/ack", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		printResult(result, func(map[string]interface{}) {
			fmt.Printf("Acknowledged change %s\n", args[0])
		})
	},
}

var complianceBulkAckCmd = &cobra.Command{
	Use:   "bulk-ack [diff-ids...]",
	Short: "Acknowledge multiple policy changes",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		diffIDs := make([]int, 0, len(args))
		for _, arg := range args {
			diffID, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error: diff-id must be an integer: %s\n", arg)
				return
			}
			diffIDs = append(diffIDs, diffID)
		}

		reqBody := map[string]interface{}{
			"diff_ids": diffIDs,
		}

		notes, _ := cmd.Flags().GetString("notes")
		if notes != "" {
			reqBody["notes"] = notes
		}

		var result map[string]interface{}
		if err := c.Post("/compliance/ack/bulk", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		printResult(result, printBulkAckResult)
	},
}

var complianceStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get compliance statistics",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		var result map[string]interface{}
		if err := c.Get("/compliance/stats", &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		printResult(result, printComplianceStats)
	},
}

func init() {
	rootCmd.AddCommand(complianceCmd)
	complianceCmd.AddCommand(complianceUnreviewedCmd)
	complianceCmd.AddCommand(complianceAckCmd)
	complianceCmd.AddCommand(complianceBulkAckCmd)
	complianceCmd.AddCommand(complianceStatsCmd)

	complianceUnreviewedCmd.Flags().String("change-type", "", "Filter by change type")
	complianceUnreviewedCmd.Flags().String("cursor", "", "Pagination cursor")
	complianceUnreviewedCmd.Flags().IntP("limit", "l", 50, "Results per page (1-100)")
	complianceAckCmd.Flags().String("notes", "", "Optional acknowledgment notes")
	complianceBulkAckCmd.Flags().String("notes", "", "Optional acknowledgment notes")
}

func printResult(result map[string]interface{}, tablePrinter func(map[string]interface{})) {
	if getOutput() == "json" {
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
		return
	}

	tablePrinter(result)
}

func printUnreviewedChanges(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No unreviewed changes")
		return
	}

	for _, item := range data {
		change := item.(map[string]interface{})
		fmt.Printf("%v %v %v\n", change["diff_id"], change["policy_id"], change["change_type"])
		if summary, ok := change["change_summary"].(string); ok && summary != "" {
			fmt.Printf("  %s\n", summary)
		}
	}
}

func printBulkAckResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Bulk acknowledgment completed")
		return
	}

	fmt.Printf("Acknowledged: %v\n", data["acknowledged"])
	fmt.Printf("Already Acknowledged: %v\n", data["already_acked"])
	fmt.Printf("Total: %v\n", data["total"])
}

func printComplianceStats(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("No compliance stats found")
		return
	}

	fmt.Printf("Total Changes (30d): %v\n", data["total_changes_30d"])
	fmt.Printf("Acknowledged: %v\n", data["acknowledged_count"])
	fmt.Printf("Unreviewed: %v\n", data["unreviewed_count"])
	fmt.Printf("Acknowledgment Rate: %v%%\n", data["acknowledgment_rate"])
	fmt.Printf("Critical Unreviewed: %v\n", data["critical_unreviewed"])
}
