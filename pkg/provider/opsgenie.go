package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpsGenieProvider implements the Provider interface for OpsGenie
type OpsGenieProvider struct {
	apiKey   string
	endpoint string
}

// OpsGenieAlert represents the alert structure for OpsGenie
type OpsGenieAlert struct {
	Message     string                 `json:"message"`
	Description string                 `json:"description,omitempty"`
	Priority    string                 `json:"priority,omitempty"`
	Source      string                 `json:"source,omitempty"`
	Entity      string                 `json:"entity,omitempty"`
	Alias       string                 `json:"alias,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// NewOpsGenieProvider creates a new OpsGenie provider
func NewOpsGenieProvider(apiKey, endpoint string) *OpsGenieProvider {
	if endpoint == "" {
		endpoint = "https://api.opsgenie.com/v2/alerts"
	}
	
	return &OpsGenieProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

// Name returns the provider name
func (p *OpsGenieProvider) Name() string {
	return "opsgenie"
}

// SendAlert sends an alert to OpsGenie
func (p *OpsGenieProvider) SendAlert(ctx context.Context, alert Alert) error {
	opsAlert := OpsGenieAlert{
		Message:     alert.Message,
		Description: fmt.Sprintf("Alert generated via AlertCLI at %s", alert.Timestamp.Format(time.RFC3339)),
		Priority:    mapPriority(alert.Priority),
		Source:      alert.Source,
		Entity:      alert.Source,
		Alias:       alert.ID,
		Details:     alert.Details,
	}
	
	data, err := json.Marshal(opsAlert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %v", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+p.apiKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send alert: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}
	
	return nil
}

// mapPriority maps generic priority to OpsGenie priority
func mapPriority(priority string) string {
	switch priority {
	case "critical":
		return "P1"
	case "high":
		return "P2"
	case "medium":
		return "P3"
	case "low":
		return "P4"
	default:
		return "P3"
	}
}
