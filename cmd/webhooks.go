package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tylergibbs1/backwork-cli/pkg/client"
)

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage webhooks",
	Long:  "Create, list, update, delete, and test webhook subscriptions",
}

var webhooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		var result map[string]interface{}
		if err := c.Get("/webhooks", &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printWebhooksList(result)
		}
	},
}

var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new webhook",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		url, _ := cmd.Flags().GetString("url")
		events, _ := cmd.Flags().GetString("events")

		reqBody := map[string]interface{}{
			"url":    url,
			"events": strings.Split(events, ","),
		}

		var result map[string]interface{}
		if err := c.Post("/webhooks", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printWebhookDetail(result)
		}
	},
}

var webhooksUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		webhookID := args[0]
		c := client.New(getAPIKey(), getBaseURL())

		reqBody := map[string]interface{}{}

		url, _ := cmd.Flags().GetString("url")
		if url != "" {
			reqBody["url"] = url
		}

		events, _ := cmd.Flags().GetString("events")
		if events != "" {
			reqBody["events"] = strings.Split(events, ",")
		}

		path := fmt.Sprintf("/webhooks/%s", webhookID)

		var result map[string]interface{}
		if err := c.Request("PATCH", path, reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printWebhookDetail(result)
		}
	},
}

var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		webhookID := args[0]
		c := client.New(getAPIKey(), getBaseURL())

		path := fmt.Sprintf("/webhooks/%s", webhookID)

		var result map[string]interface{}
		if err := c.Request("DELETE", path, nil, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			fmt.Printf("Webhook %s deleted successfully\n", webhookID)
		}
	},
}

var webhooksTestCmd = &cobra.Command{
	Use:   "test [id]",
	Short: "Send a test event to a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		webhookID := args[0]
		c := client.New(getAPIKey(), getBaseURL())

		path := fmt.Sprintf("/webhooks/%s/test", webhookID)

		var result map[string]interface{}
		if err := c.Post(path, nil, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printWebhookTestResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(webhooksCmd)
	webhooksCmd.AddCommand(webhooksListCmd)
	webhooksCmd.AddCommand(webhooksCreateCmd)
	webhooksCmd.AddCommand(webhooksUpdateCmd)
	webhooksCmd.AddCommand(webhooksDeleteCmd)
	webhooksCmd.AddCommand(webhooksTestCmd)

	webhooksCreateCmd.Flags().String("url", "", "Webhook endpoint URL")
	webhooksCreateCmd.Flags().String("events", "", "Comma-separated event types (e.g., policy.created,policy.updated)")
	webhooksCreateCmd.MarkFlagRequired("url")
	webhooksCreateCmd.MarkFlagRequired("events")

	webhooksUpdateCmd.Flags().String("url", "", "New webhook endpoint URL")
	webhooksUpdateCmd.Flags().String("events", "", "New comma-separated event types")
}

func printWebhooksList(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No webhooks found")
		return
	}

	fmt.Printf("Found %d webhooks:\n\n", len(data))
	for _, w := range data {
		webhook := w.(map[string]interface{})
		fmt.Printf("ID: %v\n", webhook["id"])
		fmt.Printf("URL: %v\n", webhook["url"])
		if events, ok := webhook["events"].([]interface{}); ok {
			eventStrs := make([]string, len(events))
			for i, e := range events {
				eventStrs[i] = fmt.Sprintf("%v", e)
			}
			fmt.Printf("Events: %s\n", strings.Join(eventStrs, ", "))
		}
		if status, ok := webhook["status"].(string); ok {
			fmt.Printf("Status: %s\n", status)
		}
		if createdAt, ok := webhook["created_at"].(string); ok {
			fmt.Printf("Created: %s\n", createdAt)
		}
		fmt.Println("---")
	}
}

func printWebhookDetail(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("ID: %v\n", data["id"])
	fmt.Printf("URL: %v\n", data["url"])
	if events, ok := data["events"].([]interface{}); ok {
		eventStrs := make([]string, len(events))
		for i, e := range events {
			eventStrs[i] = fmt.Sprintf("%v", e)
		}
		fmt.Printf("Events: %s\n", strings.Join(eventStrs, ", "))
	}
	if status, ok := data["status"].(string); ok {
		fmt.Printf("Status: %s\n", status)
	}
	if secret, ok := data["secret"].(string); ok && secret != "" {
		fmt.Printf("Secret: %s\n", secret)
	}
	if createdAt, ok := data["created_at"].(string); ok {
		fmt.Printf("Created: %s\n", createdAt)
	}
}

func printWebhookTestResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Test Result: %v\n", data["status"])
	if statusCode, ok := data["status_code"]; ok {
		fmt.Printf("Response Code: %v\n", statusCode)
	}
	if duration, ok := data["duration_ms"]; ok {
		fmt.Printf("Duration: %vms\n", duration)
	}
	if errMsg, ok := data["error"].(string); ok && errMsg != "" {
		fmt.Printf("Error: %s\n", errMsg)
	}
}
