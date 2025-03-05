package cmd

import (
	"fmt"
	"time"

	"github.com/copydataai/fake-backend-alerts/pkg/provider"
	"github.com/spf13/cobra"
)

var (
	providerName string
	apiKey       string
	endpoint     string
	severity     string
	message      string
	source       string
	priority     string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an alert to a provider",
	Long:  `Send an individual alert to a specified provider like OpsGenie, PagerDuty, etc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := provider.GetProvider(providerName, apiKey, endpoint)
		if err != nil {
			return fmt.Errorf("failed to initialize provider: %v", err)
		}

		alert := provider.Alert{
			ID:        fmt.Sprintf("alert-%d", time.Now().Unix()),
			Message:   message,
			Severity:  severity,
			Source:    source,
			Priority:  priority,
			Timestamp: time.Now(),
		}

		if err := p.SendAlert(cmd.Context(), alert); err != nil {
			return fmt.Errorf("failed to send alert: %v", err)
		}

		fmt.Printf("Successfully sent alert to %s\n", providerName)
		return nil
	},
}

func init() {
	sendCmd.Flags().StringVar(&providerName, "provider", "", "Provider name (required): opsgenie, pagerduty")
	sendCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for the provider")
	sendCmd.Flags().StringVar(&endpoint, "endpoint", "", "Custom endpoint URL (optional)")
	sendCmd.Flags().StringVar(&severity, "severity", "warning", "Alert severity: critical, error, warning, info")
	sendCmd.Flags().StringVar(&message, "message", "Test alert", "Alert message")
	sendCmd.Flags().StringVar(&source, "source", "alertcli", "Alert source")
	sendCmd.Flags().StringVar(&priority, "priority", "medium", "Alert priority: low, medium, high, critical")

	sendCmd.MarkFlagRequired("provider")
}
