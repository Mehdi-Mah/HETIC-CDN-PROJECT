import { useEffect, useState } from "react";
import axios from "axios";
import { jwtDecode } from "jwt-decode";
import "../styles/files.scss";
import {handleFileDownload, createFolder, handleDelete} from "../utils/fileHandler"

const FileManager = () => {
	const [files, setFiles] = useState([]);
	const [username, setUsername] = useState("Utilisateur");
	const [expandedFolders, setExpandedFolders] = useState({});
	const [selectedFolderID, setSelectedFolderID] = useState(0);
	const [error, setError] = useState(null);
	const [uploadError, setUploadError] = useState(null);
	const [uploadFile, setUploadFile] = useState(null);
	const [uploadPath, setUploadPath] = useState("/");
	const [ROOT_FOLDER_NAME, setRootFolderName] = useState(`Fichiers de "${username}"`);
	const [newFolderName, setNewFolderName] = useState("");
	const [initPage, setInitPage] = useState(false)

	let idCounter = 0;

	// Get username from JWT token
	useEffect(() => {
		const token = localStorage.getItem("token");
		if (token) {
			try {
				const decoded = jwtDecode(token);
				if (decoded.username) {
					setUsername(decoded.username);

					const rootName = `Fichiers de "${decoded.username}"`;
					setRootFolderName(rootName);
				}
			} catch (err) {
				console.error("Erreur de décodage du token:", err);
			}
		}
		if(!initPage) {
			setInitPage(true);
			fetchFiles();
		}
	}, []);

	const fetchFiles = async () => {
		try {
			const response = await axios.get(`${import.meta.env.VITE_API_URL}/files`, {
				headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
			});

			// Assign a unique `id` to every file and folder
			if (Array.isArray(response.data)) {
				response.data = response.data.map((file) => ({id: file.id ?? idCounter++, ...file}));
				setFiles(response.data);
			} else {
				setFiles([]);
				setError("Format de réponse API non valide");
			}
		} catch (err) {
			setError("Échec de la récupération des données");
		}
	};

	// Handle File Selection
	const handleFileChange = (e) => {
		setUploadFile(e.target.files[0]);
	};

	// Handle File Upload
	const handleFileUpload = async () => {
		if (!uploadFile) {
			setUploadError("Aucun fichier sélectionné.");
			return;
		}
	
		const formData = new FormData();
		formData.append("file", uploadFile);
		formData.append("path", uploadPath); // Envoie le chemin du dossier sélectionné
	
		try {
			await axios.post(`${import.meta.env.VITE_API_URL}/upload-file`, formData, {
				headers: {
					Authorization: `Bearer ${localStorage.getItem("token")}`,
					"Content-Type": "multipart/form-data",
				},
			});
			setUploadFile(null);
			setUploadError(null);
			fetchFiles(); // Rafraîchir la liste après l'upload
		} catch (err) {
			setUploadError("Échec du téléchargement du fichier.");
		}
	};
	
	// Toggle folder expansion using `id`
	const toggleFolder = (folderID, folderPath) => {
		setExpandedFolders((prev) => ({
			...prev,
			[folderID]: !prev[folderID],
		}));

		// Update the path for file uploads
		if (!expandedFolders[folderID]) {
			setUploadPath(folderPath);
		}

		setSelectedFolderID(folderID);
		setUploadPath(folderPath);
	};

	// Convert files/folders into a structured tree format under the dynamic root folder
	const renderFileTree = () => {
		const tree = { [ROOT_FOLDER_NAME]: { _files: [], id: idCounter++, path: "/" } };
		const folderIds = {};

		files.forEach((file) => {
			if (file.path === "" || file.path === "/") {
				tree[ROOT_FOLDER_NAME]._files.push(file);
				return;
			}

			const parts = file.path.split("/").filter(Boolean);
			let current = tree[ROOT_FOLDER_NAME];

			for (let i = 0; i < parts.length; i++) {
				const part = parts[i];

				if (!current[part]) {
					const folderPath = parts.slice(0, i + 1).join("/");
					if (!folderIds[folderPath]) {
						folderIds[folderPath] = idCounter++;
					}

					current[part] = { _files: [], id: folderIds[folderPath], path: folderPath };
				}

				if (!current[part]._files) {
					current[part]._files = [];
				}

				if (i === parts.length - 1 && file.type === "file") {
					current[part]._files.push(file);
				}

				current = current[part];
			}
		});

		return tree;
	};

	// Recursive function to render folders and files
	const renderFoldersAndFiles = (structure, depth = 0) => {
		return Object.keys(structure)
			.filter((folderName) => folderName !== "_files" && folderName !== "id" && folderName !== "path")
			.map((folderName) => {
				const folder = structure[folderName];
	
				const isRoot = folderName === ROOT_FOLDER_NAME;
				const folderID = folder.id;
				const folderPath = folder.path;
				const isExpanded = isRoot ? true : (expandedFolders[folderID] ?? false);
				const folderIcon = isRoot ? "🏠" : isExpanded ? "📂" : "📁";
				const isSelected = folderID === selectedFolderID;
	
				const hasFiles = folder._files?.length > 0;
				const hasSubfolders = Object.keys(folder).some((key) => key !== "_files" && key !== "id" && key !== "path");
				const isEmpty = !hasFiles && !hasSubfolders;
				
				return (
					<div className="folder-content" key={folderID} depth={depth}>
							<div
							className="folder-item"
							{...(isRoot ? { root: "" } : {})}
							{...(isExpanded ? { expanded: "" } : {})}
							{...(isSelected ? { selected: "" } : {})}
							onClick={() => toggleFolder(folderID, folderPath)}
							>
								<span>
									{folderIcon} {folderName}
								</span>
								{!isRoot && 
									<button onClick={() => handleDelete({ name: folderName, path: folderPath, fetchFiles })} className="button-delete">
										❌
									</button>								
								}
							</div>
						{isExpanded && (
							<div className="nested-content">
								{renderFoldersAndFiles(folder, depth + 1)}
								{folder._files?.map((file) => (
									<div className="file-item" key={file.id}>
										<span className="file" onClick={() => handleFileDownload(file)}>
											📄 {file.name}
										</span>
									</div>
								))}
								
								{isEmpty && <div className="folder-item" empty="">📂 Dossier vide</div>}
							</div>
						)}
					</div>
				);
			});
	};


	return (
		<div className="file-manager">
			<div className="title">
				<h2>Gestionnaire de Fichiers</h2>
			</div>

			{/* Error Messages */}
			{error && <p className="error">{error}</p>}

			{/* Folder Tree with Expand/Collapse */}
			<div className="folder-tree">{renderFoldersAndFiles(renderFileTree())}</div>

			{/* File Upload */}
			<div className="file-upload">
				<label htmlFor="file-upload">
					<input id="file-upload" type="file" onChange={handleFileChange} />
					<button>Choisir un fichier</button>
				</label>
				<button onClick={handleFileUpload}>Téléverser</button>
				{uploadError && <p className="error">{uploadError}</p>}
			</div>

			<div className="folder-creation">
				<input
					type="text"
					placeholder="Nom du dossier"
					value={newFolderName}
					onChange={(e) => setNewFolderName(e.target.value)}
				/>
				<button onClick={() => createFolder(newFolderName, setNewFolderName, uploadPath, fetchFiles)}>Créer Dossier</button>
			</div>
		</div>
	);
};

export default FileManager;
