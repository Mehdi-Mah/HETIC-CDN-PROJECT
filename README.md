# HETIC-CDN-PROJECT

## Groupe 9

- HOUENOU Johannes 
- Axel Lanyan 
- Mehdi Mahmoudi
- Nadhmi KHADHRI 
- Arouna Konake
- Adam Ramouki

## 🎯 Fonctionnalités

- Création compte
- Authentification
- Upload fichier
- Téléchargement fichier
- Création dossier
- Suppression dossier et fichier
- Limitation à 10 niveaux de profondeur pour la hiérarchie des dossiers
- Load Balancing et Scalabilité
- Proxy et Cache
- Conteneurisation du CDN avec Docker 

## 📥 Installation

### Cloner les repositories + 🐳 Docker

```bash
git clone https://github.com/Mehdi-Mah/cdn-go-project.git
cd cdn-go-project
docker compose up --build -d
```

### Configuration des variables d'environnement
```bash
#.env Back
MONGO_URI=mongodb://username:password@localhost:27017/namedb?authSource=admin
MONGO_USERNAME=username
MONGO_PASSWORD=password
MONGO_DATABASE=namedb

#.env Front
VITE_API_URL=http://localhost:8081
```

## 🔌 API Endpoints
```bash
POST /register        - Inscription utilisateur
POST /login           - Connexion utilisateur
GET  /files           - Récupére tous les fichiers et dossier du user connecté
POST /create-folder   - Création de dossier
POST /upload-file     - Upload un fichier
GET /download         - Récupére contenu fichier pour le télécharger
DELETE /delete        - Supression Fichier ou Dossier


```
