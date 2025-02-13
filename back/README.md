# HETIC-CDN-PROJECT

### PROJECT-ARCHITECTURE
```
go-cdn/
├── .github/
│   └── workflows/
│       └── deploy.yml       # Pipeline CI-CD
├── cmd/
│   └── cdn/
│       └── main.go          # Point d'entrée principal
├── pkg/
│   ├── proxy/
│   │   ├── proxy.go         # Code du proxy
│   │   └── proxy_test.go   # Tests unitaires
│   ├── lb/
│   │   ├── balancer.go      # Load Balancer
│   │   └── balancer_test.go
│   └── cache/
│       ├── cache.go         # Système de cache
│       └── cache_test.go
├── configs/
│   ├── prometheus.yml       # Configuration monitoring
│   └── production.env       # Variables d'environnement
├── deployments/
│   ├── docker-compose.yml   # Pour les tests locaux
│   └── kubernetes/          # Fichiers K8s (futur)
├── scripts/
│   ├── setup_server.sh      # Script d'installation
│   └── health_check.sh      # Vérification déploiement
├── test/
│   └── integration/         # Tests d'intégration
├── docs/
│   └── ARCHITECTURE.md      # Documentation technique
├── .gitignore               # Fichiers à ignorer
├── go.mod                   # Dépendances Go
├── Dockerfile               # Configuration Docker
├── Makefile                 # Commandes utiles
└── README.md                # Guide du projet
```