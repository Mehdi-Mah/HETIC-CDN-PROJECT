# HETIC-CDN-PROJECT

### PROJECT-ARCHITECTURE
```
go-cdn/
├── cmd/
│   └── cdn/          # Main package (entrypoint)
├── pkg/
│   ├── proxy/        # Reverse proxy intelligent
│   ├── cache/        # Stratégies de cache
│   └── lb/           # Algorithmes de load balancing
├── configs/          # Fichiers de configuration
├── deployments/      # Docker-Compose/K8s manifests
├── Dockerfile        # Build optimisé multi-stage
└── go.mod            # Gestion des dépendances
```