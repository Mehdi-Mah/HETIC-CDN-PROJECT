import axios from "axios";


export const handleFileDownload = async ({file, path}) => {
    if (!file || file.type !== "file") {
        console.error("Ce n'est pas un fichier téléchargeable !");
        return;
    }
    console.log("Téléchargement du fichier :", file.name);

    try {
        const response = await axios.get(`${import.meta.env.VITE_API_URL}/download`, {
            params: { path, name: file.name },
            headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
            responseType: "blob", // Pour récupérer un fichier en binaire
        });

        console.log("Fichier reçu :", response);

        // Créer une URL pour le blob
        const blob = new Blob([response.data], { type: response.headers["content-type"] });
        const link = document.createElement("a");
        link.href = URL.createObjectURL(blob);
        link.download = file.name;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    } catch (err) {
        console.error("Erreur lors du téléchargement du fichier :", err);
    }
};


export const createFolder = async (newFolderName, setNewFolderName, uploadPath, fetchFiles) => {
    if (!newFolderName.trim()) {
        alert("Veuillez entrer un nom de dossier.");
        return;
    }

    try {
        await axios.post(`${import.meta.env.VITE_API_URL}/create-folder`, {
            folderName: newFolderName,
            path: uploadPath, // Utilise le chemin sélectionné ou racine
        }, {
            headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });

        setNewFolderName(""); // Réinitialise l'input
        fetchFiles(); // Rafraîchir la liste des fichiers/dossiers
    } catch (error) {
        alert("Erreur lors de la création du dossier.");
    }
};

export const handleDelete = async (item) => {
    const {fetchFiles} = item;
    if (!item) {
        console.error("Aucun élément sélectionné !");
        return;
    }

    const confirmDelete = window.confirm(`Voulez-vous vraiment supprimer "${item.name}" ?`);
    if (!confirmDelete) return;

    console.log("Suppression de :", item.name);

    try {
        await axios.delete(`${import.meta.env.VITE_API_URL}/delete`, {
            data: { path: item.path, type: item.type, name: item.name }, // Spécifier `data` pour DELETE
            headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });

        console.log("Suppression réussie !");
        fetchFiles(); // Rafraîchir la liste après suppression
    } catch (err) {
        console.error("Erreur lors de la suppression :", err);
    }
};

