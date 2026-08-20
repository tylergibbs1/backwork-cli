package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylergibbs1/backwork-cli/pkg/client"
)

var evaluateCmd = &cobra.Command{
	Use:   "evaluate [policy-id]",
	Short: "Evaluate coverage for a policy",
	Long:  "Evaluate whether a procedure is covered under a specific policy given patient criteria",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		policyID := args[0]
		c := client.New(getAPIKey(), getBaseURL())

		reqBody := map[string]interface{}{
			"policy_id": policyID,
		}

		age, _ := cmd.Flags().GetInt("age")
		if age > 0 {
			reqBody["age"] = age
		}

		gender, _ := cmd.Flags().GetString("gender")
		if gender != "" {
			reqBody["gender"] = gender
		}

		diagnosis, _ := cmd.Flags().GetStringSlice("diagnosis")
		if len(diagnosis) > 0 {
			reqBody["diagnosis_codes"] = diagnosis
		}

		procedure, _ := cmd.Flags().GetString("procedure")
		if procedure != "" {
			reqBody["procedure_code"] = procedure
		}

		modifier, _ := cmd.Flags().GetString("modifier")
		if modifier != "" {
			reqBody["modifier"] = modifier
		}

		pos, _ := cmd.Flags().GetString("pos")
		if pos != "" {
			reqBody["place_of_service"] = pos
		}

		var result map[string]interface{}
		if err := c.Post("/coverage/evaluate", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printEvaluateResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(evaluateCmd)

	evaluateCmd.Flags().Int("age", 0, "Patient age")
	evaluateCmd.Flags().String("gender", "", "Patient gender (M, F)")
	evaluateCmd.Flags().StringSliceP("diagnosis", "d", []string{}, "Diagnosis codes (ICD-10)")
	evaluateCmd.Flags().StringP("procedure", "p", "", "Procedure code (CPT/HCPCS)")
	evaluateCmd.Flags().StringP("modifier", "m", "", "Procedure modifier")
	evaluateCmd.Flags().String("pos", "", "Place of service code")
}

func printEvaluateResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	covered := fmt.Sprintf("%v", data["covered"])
	if covered == "true" {
		fmt.Printf("Coverage: COVERED\n")
	} else {
		fmt.Printf("Coverage: NOT COVERED\n")
	}

	if confidence, ok := data["confidence"]; ok {
		fmt.Printf("Confidence: %v\n", confidence)
	}

	if reasons, ok := data["reasons"].([]interface{}); ok && len(reasons) > 0 {
		fmt.Println("\nReasons:")
		for _, reason := range reasons {
			fmt.Printf("  - %v\n", reason)
		}
	}

	if policyID, ok := data["policy_id"].(string); ok && policyID != "" {
		fmt.Printf("\nPolicy: %s\n", policyID)
	}

	if criteria, ok := data["matched_criteria"].([]interface{}); ok && len(criteria) > 0 {
		fmt.Println("\nMatched Criteria:")
		for _, c := range criteria {
			if cMap, ok := c.(map[string]interface{}); ok {
				if section, ok := cMap["section"].(string); ok {
					fmt.Printf("  [%s] ", section)
				}
				if text, ok := cMap["text"].(string); ok {
					if len(text) > 120 {
						text = text[:117] + "..."
					}
					fmt.Printf("%s\n", text)
				}
			}
		}
	}
}
