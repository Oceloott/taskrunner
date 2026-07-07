package orchestrator

import "fmt"

type OrchestratorConfig struct {
	Workers int
	Verbose bool
}

type Option func(*OrchestratorConfig)

func WithWorkers(n int) Option {
	return func(c *OrchestratorConfig) { c.Workers = n }
}

func WithVerbose(v bool) Option {
	return func(c *OrchestratorConfig) { c.Verbose = v }
}

func ValidateWorkers(n int) (int, error) {
	if n < 1 || n > 100 {
		return 3, fmt.Errorf("le nombre de workers doit être entre 1 et 100 (reçu %d)", n)
	}
	return n, nil
}
