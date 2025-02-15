# Étape 1: Build du binaire Go
FROM golang:1.22.3 AS builder

WORKDIR /app

# Copier les fichiers Go et les dépendances
COPY go.mod go.sum ./
RUN go mod tidy

# Copier tout le code source
COPY . .

# Vérifier que les fichiers Go sont bien copiés
RUN ls -la /app

# Compiler l'application Go avec la bonne architecture
RUN GOOS=linux GOARCH=amd64 go build -o /app/cdn-api ./cmd/cdn/main.go

# Étape 2: Création de l'image finale avec Alpine
FROM alpine:latest

WORKDIR /root/

# Installer libc6-compat pour éviter les erreurs de binaires
RUN apk add --no-cache libc6-compat

# Copier le binaire compilé depuis le builder
COPY --from=builder /app/cdn-api .

# Vérifier que le binaire est bien présent
RUN ls -la /root/

# Donner les permissions d'exécution au binaire
RUN chmod +x /root/cdn-api

# Exposer le port API
EXPOSE 8081

# Lancer l'application
CMD ["/root/cdn-api"]
