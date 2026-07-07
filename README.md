# taskrunner

Orchestrateur de tâches concurrentes en ligne de commande. Lit un fichier de
tâches JSON, les exécute en parallèle via un pool de workers (avec timeouts et
retries), puis produit un rapport JSON.

Écrit en Go, **bibliothèque standard uniquement** (aucune dépendance externe).

## Utilisation

```bash
make build
./taskrunner -file tasks.json -workers 3
```

Flags :

| Flag | Défaut | Rôle |
|------|--------|------|
| `-file` | *(obligatoire)* | chemin vers le fichier JSON de tâches |
| `-workers` | `3` | nombre de workers simultanés (borné à [1, 100]) |
| `-verbose` | `false` | affiche le statut de chaque tâche en temps réel sur stderr |

Le rapport JSON est écrit sur **stdout**, les messages de progression sur
**stderr**. Un fichier `METRICS.md` est généré à la fin de l'exécution.

## Format du fichier de tâches

```json
{
  "tasks": [
    { "id": "t1", "type": "print", "params": { "message": "hello" }, "timeout": "2s", "retries": 0 },
    { "id": "t2", "type": "calc",  "params": { "value": 42 },        "timeout": "1s", "retries": 1 }
  ]
}
```

Types disponibles : `print`, `calc`, `download`, `fake`.

## Architecture

```
cmd/taskrunner/      point d'entrée (flags, signal SIGINT, câblage)
internal/
├── task/            interface Task, implémentations, TaskError, loader JSON
├── report/          Report (io.WriterTo), TaskResult
├── orchestrator/    Orchestrate (pool de workers), options, ValidateWorkers
└── metrics/         WriteMetrics -> METRICS.md
```

## Développement

```bash
make test     # lance les tests
make lint     # go vet + vérification gofmt
make run      # compile et exécute avec tasks.json
go test -race ./...   # tests avec détecteur de data races
```
