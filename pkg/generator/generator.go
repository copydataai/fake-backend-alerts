package generator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/copydataai/fake-backend-alerts/pkg/provider"
)

// Generator handles generating alerts for scenarios
type Generator struct {
	provider provider.Provider
}

// ScenarioOptions configures how a scenario is run
type ScenarioOptions struct {
	Count       int
	Interval    int
	Concurrency int
	Params      map[string]string
}

// ScenarioResult contains the results of a scenario run
type ScenarioResult struct {
	Sent     int
	Failed   int
	Duration time.Duration
	Rate     float64
}

// Scenario represents a predefined alert generation scenario
type Scenario struct {
	Name        string
	Description string
	Generator   func(context.Context, *Generator, ScenarioOptions) (ScenarioResult, error)
}

// NewGenerator creates a new alert generator
func NewGenerator(p provider.Provider) *Generator {
	return &Generator{
		provider: p,
	}
}

// RunScenario runs a predefined scenario
func (g *Generator) RunScenario(ctx context.Context, name string, opts ScenarioOptions) (ScenarioResult, error) {
	scenario, ok := scenarios[name]
	if !ok {
		return ScenarioResult{}, fmt.Errorf("unknown scenario: %s", name)
	}

	return scenario.Generator(ctx, g, opts)
}

// ListScenarios returns the list of available scenarios
func ListScenarios() []Scenario {
	var result []Scenario
	for _, s := range scenarios {
		result = append(result, s)
	}
	return result
}

// Available scenarios
var scenarios = map[string]Scenario{
	"escalating": {
		Name:        "escalating",
		Description: "Gradually increases alert severity over time",
		Generator:   generateEscalatingScenario,
	},
	"random": {
		Name:        "random",
		Description: "Generates alerts with random severities and priorities",
		Generator:   generateRandomScenario,
	},
	"burst": {
		Name:        "burst",
		Description: "Sends alerts in bursts with pauses in between",
		Generator:   generateBurstScenario,
	},
	"mixed": {
		Name:        "mixed",
		Description: "Mix of different alert types and severities",
		Generator:   generateMixedScenario,
	},
}

// generateEscalatingScenario generates alerts with escalating severity
func generateEscalatingScenario(ctx context.Context, g *Generator, opts ScenarioOptions) (ScenarioResult, error) {
	severities := []string{"info", "warning", "error", "critical"}
	priorities := []string{"low", "medium", "high", "critical"}
	
	total := opts.Count
	interval := time.Duration(opts.Interval) * time.Millisecond
	result := ScenarioResult{}
	
	start := time.Now()
	
	for i := 0; i < total; i++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			// Calculate severity index based on progress
			severityIndex := (i * len(severities)) / total
			if severityIndex >= len(severities) {
				severityIndex = len(severities) - 1
			}
			
			priorityIndex := (i * len(priorities)) / total
			if priorityIndex >= len(priorities) {
				priorityIndex = len(priorities) - 1
			}
			
			alert := provider.Alert{
				ID:        fmt.Sprintf("escalating-%d", i),
				Message:   fmt.Sprintf("Escalating alert #%d", i),
				Severity:  severities[severityIndex],
				Priority:  priorities[priorityIndex],
				Source:    "scenario-escalating",
				Timestamp: time.Now(),
				Details: map[string]interface{}{
					"scenario": "escalating",
					"progress": float64(i) / float64(total),
					"index":    i,
				},
			}
			
			err := g.provider.SendAlert(ctx, alert)
			if err != nil {
				result.Failed++
			} else {
				result.Sent++
			}
			
			// Sleep for the interval unless it's the last alert
			if i < total-1 {
				time.Sleep(interval)
			}
		}
	}
	
	result.Duration = time.Since(start)
	if result.Duration.Seconds() > 0 {
		result.Rate = float64(result.Sent) / result.Duration.Seconds()
	}
	
	return result, nil
}

// generateRandomScenario generates alerts with random properties
func generateRandomScenario(ctx context.Context, g *Generator, opts ScenarioOptions) (ScenarioResult, error) {
	severities := []string{"info", "warning", "error", "critical"}
	priorities := []string{"low", "medium", "high", "critical"}
	
	// Set up worker pool
	wg := &sync.WaitGroup{}
	jobs := make(chan int, opts.Count)
	results := make(chan bool, opts.Count)
	
	// Start the workers
	for w := 0; w < opts.Concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				sevIdx := rand.Intn(len(severities))
				priIdx := rand.Intn(len(priorities))
				
				alert := provider.Alert{
					ID:        fmt.Sprintf("random-%d", i),
					Message:   fmt.Sprintf("Random alert #%d", i),
					Severity:  severities[sevIdx],
					Priority:  priorities[priIdx],
					Source:    "scenario-random",
					Timestamp: time.Now(),
					Details: map[string]interface{}{
						"scenario": "random",
						"index":    i,
					},
				}
				
				err := g.provider.SendAlert(ctx, alert)
				results <- (err == nil)
			}
		}()
	}
	
	// Send jobs to the workers
	start := time.Now()
	go func() {
		for i := 0; i < opts.Count; i++ {
			jobs <- i
			// Add a small delay between job submissions to control the rate
			if opts.Interval > 0 {
				time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
			}
		}
		close(jobs)
	}()
	
	// Collect results
	result := ScenarioResult{}
	for i := 0; i < opts.Count; i++ {
		if <-results {
			result.Sent++
		} else {
			result.Failed++
		}
	}
	
	// Wait for all workers to finish
	wg.Wait()
	close(results)
	
	result.Duration = time.Since(start)
	if result.Duration.Seconds() > 0 {
		result.Rate = float64(result.Sent) / result.Duration.Seconds()
	}
	
	return result, nil
}

// generateBurstScenario generates alerts in bursts with pauses in between
func generateBurstScenario(ctx context.Context, g *Generator, opts ScenarioOptions) (ScenarioResult, error) {
	severities := []string{"info", "warning", "error", "critical"}
	burstSize := 10
	if v, ok := opts.Params["burst_size"]; ok {
		fmt.Sscanf(v, "%d", &burstSize)
	}
	
	pauseDuration := 2000 // milliseconds
	if v, ok := opts.Params["pause_duration"]; ok {
		fmt.Sscanf(v, "%d", &pauseDuration)
	}
	
	start := time.Now()
	result := ScenarioResult{}
	
	for i := 0; i < opts.Count; i++ {
		// Determine if this is a burst boundary
		if i > 0 && i%burstSize == 0 {
			fmt.Printf("Pausing for %d ms after burst\n", pauseDuration)
			time.Sleep(time.Duration(pauseDuration) * time.Millisecond)
		}
		
		sevIdx := rand.Intn(len(severities))
		alert := provider.Alert{
			ID:        fmt.Sprintf("burst-%d", i),
			Message:   fmt.Sprintf("Burst alert #%d", i),
			Severity:  severities[sevIdx],
			Priority:  "high",
			Source:    "scenario-burst",
			Timestamp: time.Now(),
			Details: map[string]interface{}{
				"scenario":   "burst",
				"index":      i,
				"burstIndex": i % burstSize,
			},
		}
		
		err := g.provider.SendAlert(ctx, alert)
		if err != nil {
			result.Failed++
		} else {
			result.Sent++
		}
		
		// Sleep within a burst
		if i%burstSize != burstSize-1 && i < opts.Count-1 {
			time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
		}
	}
	
	result.Duration = time.Since(start)
	if result.Duration.Seconds() > 0 {
		result.Rate = float64(result.Sent) / result.Duration.Seconds()
	}
	
	return result, nil
}

// generateMixedScenario generates a mix of different alert types
func generateMixedScenario(ctx context.Context, g *Generator, opts ScenarioOptions) (ScenarioResult, error) {
	templates := []struct {
		severity  string
		priority  string
		message   string
		category  string
	}{
		{"info", "low", "System startup complete", "system"},
		{"info", "low", "User logged in", "user"},
		{"warning", "medium", "High CPU usage detected", "performance"},
		{"warning", "medium", "Low disk space", "system"},
		{"error", "high", "Database connection failed", "database"},
		{"error", "high", "Authentication failure", "security"},
		{"critical", "critical", "Service unavailable", "service"},
		{"critical", "critical", "Security breach detected", "security"},
	}
	
	start := time.Now()
	result := ScenarioResult{}
	
	// Set up worker pool
	wg := &sync.WaitGroup{}
	jobs := make(chan int, opts.Count)
	results := make(chan bool, opts.Count)
	
	// Start the workers
	for w := 0; w < opts.Concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				// Select a random template
				tmpl := templates[rand.Intn(len(templates))]
				
				alert := provider.Alert{
					ID:        fmt.Sprintf("mixed-%d", i),
					Message:   fmt.Sprintf("%s (%d)", tmpl.message, i),
					Severity:  tmpl.severity,
					Priority:  tmpl.priority,
					Source:    "scenario-mixed",
					Timestamp: time.Now(),
					Details: map[string]interface{}{
						"scenario": "mixed",
						"category": tmpl.category,
						"index":    i,
					},
				}
				
				err := g.provider.SendAlert(ctx, alert)
				results <- (err == nil)
			}
		}()
	}
	
	// Send jobs to the workers
	go func() {
		for i := 0; i < opts.Count; i++ {
			jobs <- i
			// Add a small delay between job submissions
			if opts.Interval > 0 {
				time.Sleep(time.Duration(opts.Interval) * time.Millisecond)
			}
		}
		close(jobs)
	}()
	
	// Collect results
	for i := 0; i < opts.Count; i++ {
		if <-results {
			result.Sent++
		} else {
			result.Failed++
		}
	}
	
	// Wait for all workers to finish
	wg.Wait()
	close(results)
	
	result.Duration = time.Since(start)
	if result.Duration.Seconds() > 0 {
		result.Rate = float64(result.Sent) / result.Duration.Seconds()
	}
	
	return result, nil
}
