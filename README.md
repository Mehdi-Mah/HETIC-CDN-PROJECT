# HETIC-CDN-PROJECT

# PROJECT-ARCHITECTURE
go-cdn/
├── cmd/
│   └── cdn/          # Main package (entrypoint)
├── pkg/
│   ├── proxy/        # Reverse proxy intelligent
│   ├── cache/        # Stratégies de cache (LRU, TTL)
│   └── lb/           # Algorithmes de load balancing
├── configs/          # Fichiers de configuration (YAML/ENV)
├── deployments/      # Docker-Compose/K8s manifests
├── Dockerfile        # Build optimisé multi-stage
└── go.mod            # Gestion des dépendances