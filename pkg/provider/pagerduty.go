package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PagerDutyProvider implements the Provider interface for PagerDuty
type PagerDutyProvider struct {
	apiKey   string
	endpoint string
}

// PagerDutyEvent represents the event structure for PagerDuty
type PagerDutyEvent struct {
	RoutingKey  string                `json:"routing_key"`
	EventAction string                `json:"event_action"`
	Payload     PagerDutyEventPayload `json:"payload"`
}

// PagerDutyEventPayload represents the event payload for PagerDuty
type PagerDutyEventPayload struct {
	Summary   string                 `json:"summary"`
	Source    string                 `json:"source"`
	Severity  string                 `json:"severity"`
	Timestamp string                 `json:"timestamp"`
	Component string                 `json:"component,omitempty"`
	Group     string                 `json:"group,omitempty"`
	Class     string                 `json:"class,omitempty"`
	Details   map[string]interface{} `json:"custom_details,omitempty"`
}

// NewPagerDutyProvider creates a new PagerDuty provider
func NewPagerDutyProvider(apiKey, endpoint string) *PagerDutyProvider {
	if endpoint == "" {
		endpoint = "https://events.pagerduty.com/v2/enqueue"
	}

	return &PagerDutyProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

// Name returns the provider name
func (p *PagerDutyProvider) Name() string {
	return "pagerduty"
}

// SendAlert sends an alert to PagerDuty
func (p *PagerDutyProvider) SendAlert(ctx context.Context, alert Alert) error {
	event := PagerDutyEvent{
		RoutingKey:  p.apiKey,
		EventAction: "trigger",
		Payload: PagerDutyEventPayload{
			Summary:   alert.Message,
			Source:    alert.Source,
			Severity:  mapSeverity(alert.Severity),
			Timestamp: alert.Timestamp.Format(time.RFC3339),
			Component: "AlertCLI",
			Group:     "Testing",
			Class:     "stress-test",
			Details:   alert.Details,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send alert: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return fmt.Errorf("request failed with status: %s, body: %s", resp.Status, body)
	}

	return nil
}

// mapSeverity maps generic severity to PagerDuty severity
func mapSeverity(severity string) string {
	switch severity {
	case "critical":
		return "critical"
	case "error":
		return "error"
	case "warning":
		return "warning"
	case "info":
		return "info"
	default:
		return "warning"
	}
}
