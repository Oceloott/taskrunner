package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"taskrunner/internal/metrics"
	"taskrunner/internal/orchestrator"
	"taskrunner/internal/task"
)

func main() {
	file := flag.String("file", "", "chemin vers le fichier JSON de tâches (obligatoire)")
	workers := flag.Int("workers", 3, "nombre de workers simultanés")
	verbose := flag.Bool("verbose", false, "affiche le statut des tâches en temps réel sur stderr")
	flag.Parse()

	if *file == "" {
		fmt.Fprintln(os.Stderr, "erreur: le flag -file est obligatoire")
		flag.Usage()
		os.Exit(2)
	}

	w, err := orchestrator.ValidateWorkers(*workers)
	if err != nil {
		fmt.Fprintln(os.Stderr, "avertissement:", err)
	}

	tasks, err := task.Load(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur de chargement:", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rep, err := orchestrator.Orchestrate(ctx, tasks, w, orchestrator.WithVerbose(*verbose))
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur d'exécution:", err)
		os.Exit(1)
	}

	if _, err := rep.WriteTo(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "erreur d'écriture du rapport:", err)
		os.Exit(1)
	}

	content := metrics.WriteMetrics(rep.Results)
	if err := os.WriteFile("METRICS.md", []byte(content), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "erreur d'écriture de METRICS.md:", err)
	}
}
