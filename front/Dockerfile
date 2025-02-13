# Étape 1 : Build de l'application Vite
FROM node:20-alpine AS builder

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers package.json et package-lock.json
COPY ./package.json ./package-lock.json ./

# Installer les dépendances
RUN npm install --frozen-lockfile

# Copier le reste du code
COPY . .

# Construire l'application
RUN npm run build

# Étape 2 : Image de production avec Nginx
FROM nginx:alpine

# Copier les fichiers de build dans le dossier Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copier le fichier de configuration personnalisé de Nginx
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Exposer le bon port (80)
EXPOSE 80

# Démarrer Nginx
CMD ["nginx", "-g", "daemon off;"]
