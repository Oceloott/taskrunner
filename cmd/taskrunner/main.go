package main

import (
	"flag"
	"fmt"
	"os"
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

	// TODO: charger les tâches, appeler Orchestrate, écrire le rapport.
	fmt.Fprintf(os.Stderr, "taskrunner: file=%s workers=%d verbose=%t\n", *file, *workers, *verbose)
}
