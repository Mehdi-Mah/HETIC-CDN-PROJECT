import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { jwtDecode } from "jwt-decode";
import "../styles/auth.scss";

const Login = () => {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [error, setError] = useState(null);
	const navigate = useNavigate();

	const handleLogin = async (e) => {
		e.preventDefault();
		setError(null);

		try {
			const response = await axios.post(
				`${import.meta.env.VITE_API_URL}/login`,
				{ email, password }
			);

			const { token } = response.data;
			localStorage.setItem("token", token);
			const decoded = jwtDecode(token);
			localStorage.setItem("token_expiry", decoded.exp);

			navigate("/files"); // Redirect to the file manager
		} catch (err) {
			setError("Identifiants non valides. Veuillez réessayer.");
		}
	};

	return (
		<div className="auth-container">
			<h2>Connexion</h2>
			<form onSubmit={handleLogin}>
				<input
					type="email"
					placeholder="Adresse mail"
					value={email}
					onChange={(e) => setEmail(e.target.value)}
					required
				/>
				<input
					type="password"
					placeholder="Mot de passe"
					value={password}
					onChange={(e) => setPassword(e.target.value)}
					required
				/>
				<button type="submit">Se connecter</button>
				{error && <p className="error">{error}</p>}
			</form>
		</div>
	);
};

export default Login;
