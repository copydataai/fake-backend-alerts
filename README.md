# Fake Alert Backend & CLI
> A backend service and CLI to test and stress test alert management providers like OpsGenie and PagerDuty

## CLI Usage

### Build the CLI

```bash
go build -o alertcli ./cmd/alertcli
```

### Send Individual Alerts

Send a single alert to OpsGenie:
```bash
./alertcli send --provider opsgenie --api-key YOUR_API_KEY --message "Test alert from CLI" --severity critical
```

Send a single alert to PagerDuty:
```bash
./alertcli send --provider pagerduty --api-key YOUR_ROUTING_KEY --message "Test alert from CLI" --severity critical
```

### Run Stress Test Scenarios

Run the escalating severity scenario with 100 alerts:
```bash
./alertcli scenario --provider opsgenie --api-key YOUR_API_KEY --name escalating --count 100 --interval 200 --concurrency 5
```

Run the random alert scenario:
```bash
./alertcli scenario --provider pagerduty --api-key YOUR_ROUTING_KEY --name random --count 500 --interval 50 --concurrency 10
```

List available scenarios:
```bash
./alertcli scenario list
```

### Available Scenarios

1. **escalating** - Gradually increases alert severity over time
2. **random** - Generates alerts with random severities and priorities
3. **burst** - Sends alerts in bursts with pauses in between
4. **mixed** - Mix of different alert types and severities

## REST API (Fake Backend Service)

### PagerDuty-compatible endpoints:

- POST /pagerduty/v1/incidents - Create new incident
- PUT /pagerduty/v1/incidents/{id} - Update incident
- GET /pagerduty/v1/incidents - List incidents

### OpsGenie-compatible endpoints:

- POST /opsgenie/v1/alerts - Create alert
- GET /opsgenie/v1/alerts - List alerts
- POST /opsgenie/v1/alerts/{id}/acknowledge - Acknowledge alert
- POST /opsgenie/v1/alerts/{id}/close - Close alert

### Alert Generator control endpoints:

- POST /generator/run/{scenario} - Run a predefined scenario
- POST /generator/alert - Generate a single alert
- GET /generator/scenarios - List available scenarios
