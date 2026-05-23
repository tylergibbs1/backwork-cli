package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/backworkai/verity-cli/pkg/client"
	"github.com/spf13/cobra"
)

var claimsCmd = &cobra.Command{
	Use:   "claims",
	Short: "Claim validation commands",
	Long:  "Validate coverage and denial risk for procedure codes before submission",
}

var claimsValidateCmd = &cobra.Command{
	Use:   "validate [procedure-codes...]",
	Short: "Validate coverage and denial risk",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		reqBody := map[string]interface{}{
			"procedure_codes": args,
		}

		payer, _ := cmd.Flags().GetString("payer")
		planType, _ := cmd.Flags().GetString("plan-type")
		lineOfBusiness, _ := cmd.Flags().GetString("line-of-business")
		diagnosisCodes, _ := cmd.Flags().GetStringSlice("diagnosis")
		modifiers, _ := cmd.Flags().GetStringSlice("modifier")
		state, _ := cmd.Flags().GetString("state")
		dateOfService, _ := cmd.Flags().GetString("date-of-service")
		siteOfService, _ := cmd.Flags().GetString("site-of-service")
		providerSpecialty, _ := cmd.Flags().GetString("provider-specialty")
		ageCategory, _ := cmd.Flags().GetString("age-category")
		sex, _ := cmd.Flags().GetString("sex")
		idempotencyKey, _ := cmd.Flags().GetString("idempotency-key")
		legacy, _ := cmd.Flags().GetBool("legacy")

		if payer != "" {
			reqBody["payer"] = payer
		}
		if planType != "" {
			reqBody["plan_type"] = planType
		}
		if lineOfBusiness != "" {
			reqBody["line_of_business"] = lineOfBusiness
		}
		if len(diagnosisCodes) > 0 {
			reqBody["diagnosis_codes"] = diagnosisCodes
		}
		if len(modifiers) > 0 {
			reqBody["modifiers"] = modifiers
		}
		if state != "" {
			reqBody["state"] = state
		}
		if dateOfService != "" {
			reqBody["date_of_service"] = dateOfService
		}
		if siteOfService != "" {
			reqBody["site_of_service"] = siteOfService
		}
		if providerSpecialty != "" {
			reqBody["provider_specialty"] = providerSpecialty
		}
		if ageCategory != "" {
			reqBody["age_category"] = ageCategory
		}
		if sex != "" {
			reqBody["sex_when_policy_relevant"] = sex
		}

		path := "/claims/validate"
		if legacy {
			path = "/claim-validation"
		}

		headers := map[string]string{}
		if idempotencyKey != "" {
			headers["X-Idempotency-Key"] = idempotencyKey
		}

		var result map[string]interface{}
		if err := c.PostWithHeaders(path, reqBody, &result, headers); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printClaimValidationResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(claimsCmd)
	claimsCmd.AddCommand(claimsValidateCmd)

	claimsValidateCmd.Flags().String("payer", "", "Payer or policy source label")
	claimsValidateCmd.Flags().String("plan-type", "", "Plan type")
	claimsValidateCmd.Flags().String("line-of-business", "", "Line of business")
	claimsValidateCmd.Flags().StringSliceP("diagnosis", "d", []string{}, "Diagnosis codes (ICD-10)")
	claimsValidateCmd.Flags().StringSliceP("modifier", "m", []string{}, "Procedure modifiers")
	claimsValidateCmd.Flags().StringP("state", "s", "", "Two-letter state code")
	claimsValidateCmd.Flags().String("date-of-service", "", "Date of service (YYYY-MM-DD)")
	claimsValidateCmd.Flags().String("site-of-service", "", "Site of service")
	claimsValidateCmd.Flags().String("provider-specialty", "", "Provider specialty")
	claimsValidateCmd.Flags().String("age-category", "", "Age category")
	claimsValidateCmd.Flags().String("sex", "", "Sex when policy relevant")
	claimsValidateCmd.Flags().String("idempotency-key", "", "Unique request identifier for safe retries")
	claimsValidateCmd.Flags().Bool("legacy", false, "Use deprecated /claim-validation endpoint")
}

func printClaimValidationResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Coverage Status: %v\n", data["coverage_status"])
	fmt.Printf("Prior Auth Required: %s\n", formatNullableBool(data["prior_auth_required"]))
	fmt.Printf("Denial Risk: %v\n", data["denial_risk"])
	fmt.Printf("Overall Risk: %v\n", data["overall_risk"])
	fmt.Printf("Confidence: %v\n", data["confidence"])

	if requirements, ok := data["documentation_requirements"].([]interface{}); ok && len(requirements) > 0 {
		fmt.Println("\nDocumentation Requirements:")
		for _, item := range requirements {
			fmt.Printf("  - %v\n", item)
		}
	}

	if gaps, ok := data["known_gaps"].([]interface{}); ok && len(gaps) > 0 {
		fmt.Println("\nKnown Gaps:")
		for _, item := range gaps {
			fmt.Printf("  - %v\n", item)
		}
	}

	if issues, ok := data["issues"].([]interface{}); ok && len(issues) > 0 {
		fmt.Println("\nIssues:")
		for _, item := range issues {
			fmt.Printf("  - %v\n", item)
		}
	}

	if policies, ok := data["matched_policies"].([]interface{}); ok && len(policies) > 0 {
		fmt.Println("\nMatched Policies:")
		for _, item := range policies {
			policy, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			fmt.Printf("  - %v", policy["policy_id"])
			if title, ok := policy["title"].(string); ok && title != "" {
				fmt.Printf(": %s", title)
			}
			if jurisdiction, ok := policy["jurisdiction"].(string); ok && jurisdiction != "" {
				fmt.Printf(" (%s)", jurisdiction)
			}
			fmt.Println()
		}
	}
}

func formatNullableBool(value interface{}) string {
	switch v := value.(type) {
	case bool:
		if v {
			return "YES"
		}
		return "NO"
	default:
		return "UNKNOWN"
	}
}
