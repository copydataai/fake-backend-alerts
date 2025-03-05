package provider

import (
	"context"
	"fmt"
	"time"
)

// Alert represents a generic alert structure that can be adapted for different providers
type Alert struct {
	ID        string
	Message   string
	Severity  string
	Source    string
	Priority  string
	Details   map[string]interface{}
	Timestamp time.Time
}

// Provider is the interface that all alert providers must implement
type Provider interface {
	// SendAlert sends a single alert to the provider
	SendAlert(ctx context.Context, alert Alert) error
	
	// Name returns the provider name
	Name() string
}

// GetProvider returns a provider implementation based on the name
func GetProvider(name, apiKey, endpoint string) (Provider, error) {
	switch name {
	case "opsgenie":
		return NewOpsGenieProvider(apiKey, endpoint), nil
	case "pagerduty":
		return NewPagerDutyProvider(apiKey, endpoint), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
