# Fake backend 
> A fake backend service to test and try each alert manager in the market


### REST 
API Endpoints

PagerDuty-compatible endpoints:

POST /pagerduty/v1/incidents - Create new incident
PUT /pagerduty/v1/incidents/{id} - Update incident
GET /pagerduty/v1/incidents - List incidents


OpsGenie-compatible endpoints:

POST /opsgenie/v1/alerts - Create alert
GET /opsgenie/v1/alerts - List alerts
POST /opsgenie/v1/alerts/{id}/acknowledge - Acknowledge alert
POST /opsgenie/v1/alerts/{id}/close - Close alert


Alert Generator control endpoints:

POST /generator/run/{scenario} - Run a predefined scenario
POST /generator/alert - Generate a single alert
GET /generator/scenarios - List available scenarios
