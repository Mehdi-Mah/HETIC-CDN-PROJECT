FROM golang:1.22.3-alpine

WORKDIR /app

# Copier les fichiers de module et télécharger les dépendances
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copier le reste du code source
COPY . .

# Compiler l'application
RUN go build -o cdn-app ./cmd/cdn/main.go

# Exposer le port utilisé par l'application
EXPOSE 8081

# Démarrer l'application
CMD ["./cdn-app"]
