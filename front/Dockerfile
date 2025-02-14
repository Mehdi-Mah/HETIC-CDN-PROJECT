# Utiliser une image Node.js pour la construction
FROM node:18 AS builder

WORKDIR /app

# Copier les fichiers du projet
COPY package.json package-lock.json ./
RUN npm install

COPY . .

# Construire le projet
RUN npm run build

# Utiliser une image nginx pour servir l'application
FROM nginx:latest

COPY --from=builder /app/dist /usr/share/nginx/html

# Exposer le port 80
EXPOSE 80
