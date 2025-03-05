package cmd

import (
	"fmt"

	"github.com/copydataai/fake-backend-alerts/pkg/generator"
	"github.com/copydataai/fake-backend-alerts/pkg/provider"
	"github.com/spf13/cobra"
)

var (
	scenarioName   string
	count          int
	interval       int
	concurrency    int
	scenarioParams map[string]string
)

var scenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Run an alert scenario",
	Long:  `Run a predefined alert scenario to generate multiple alerts for stress testing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := provider.GetProvider(providerName, apiKey, endpoint)
		if err != nil {
			return fmt.Errorf("failed to initialize provider: %v", err)
		}

		gen := generator.NewGenerator(p)
		
		fmt.Printf("Running scenario '%s' with %d alerts at %d ms intervals using %d concurrent workers\n", 
			scenarioName, count, interval, concurrency)
		
		result, err := gen.RunScenario(cmd.Context(), scenarioName, generator.ScenarioOptions{
			Count:       count,
			Interval:    interval,
			Concurrency: concurrency,
			Params:      scenarioParams,
		})
		
		if err != nil {
			return fmt.Errorf("scenario failed: %v", err)
		}

		fmt.Printf("Scenario complete: %d alerts sent, %d failed\n", result.Sent, result.Failed)
		fmt.Printf("Duration: %v\n", result.Duration)
		fmt.Printf("Rate: %.2f alerts/sec\n", result.Rate)
		
		return nil
	},
}

var listScenariosCmd = &cobra.Command{
	Use:   "list",
	Short: "List available scenarios",
	Run: func(cmd *cobra.Command, args []string) {
		scenarios := generator.ListScenarios()
		fmt.Println("Available scenarios:")
		for _, s := range scenarios {
			fmt.Printf("- %s: %s\n", s.Name, s.Description)
		}
	},
}

func init() {
	scenarioCmd.Flags().StringVar(&providerName, "provider", "", "Provider name (required): opsgenie, pagerduty")
	scenarioCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for the provider")
	scenarioCmd.Flags().StringVar(&endpoint, "endpoint", "", "Custom endpoint URL (optional)")
	scenarioCmd.Flags().StringVar(&scenarioName, "name", "escalating", "Scenario name: escalating, random, mixed")
	scenarioCmd.Flags().IntVar(&count, "count", 100, "Number of alerts to generate")
	scenarioCmd.Flags().IntVar(&interval, "interval", 100, "Interval between alerts in milliseconds")
	scenarioCmd.Flags().IntVar(&concurrency, "concurrency", 10, "Number of concurrent alert generators")
	
	scenarioCmd.MarkFlagRequired("provider")
	
	scenarioCmd.AddCommand(listScenariosCmd)
}
