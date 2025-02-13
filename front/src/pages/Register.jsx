import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import "../styles/auth.scss";

const Register = () => {
	const [username, setUsername] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [confirmPassword, setConfirmPassword] = useState("");
	const [agreeToTerms, setAgreeToTerms] = useState(false);
	const [error, setError] = useState(null);
	const navigate = useNavigate();

	const handleRegister = async (e) => {
		e.preventDefault();
		setError(null);

		if (password !== confirmPassword) {
			setError("Les mots de passe ne correspondent pas!");
			return;
		}

		if (!agreeToTerms) {
			setError("Vous devez accepter les termes et conditions.");
			return;
		}

		try {
			await axios.post(`${import.meta.env.VITE_API_URL}/register`, {
				username, // Send username
				email,
				password,
			});

			navigate("/login"); // Redirect to login after successful registration
		} catch (err) {
			setError("Inscription échouée. Veuillez réessayer.");
		}
	};

	return (
		<div className="auth-container">
			<h2>Inscription</h2>
			<form onSubmit={handleRegister}>
				<input
					type="text"
					placeholder="Nom d'utilisateur"
					value={username}
					onChange={(e) => setUsername(e.target.value)}
					required
				/>
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
				<input
					type="password"
					placeholder="Confirmation du mot de passe"
					value={confirmPassword}
					onChange={(e) => setConfirmPassword(e.target.value)}
					required
				/>

				{/* Terms and Conditions Checkbox */}
				<label className="terms">
					<input
						type="checkbox"
						checked={agreeToTerms}
						onChange={(e) => setAgreeToTerms(e.target.checked)}
					/>
					Je suis d'accord avec les <a href="/terms">termes et conditions</a>.
				</label>

				<button type="submit" disabled={!agreeToTerms}>S'inscrire</button>

				{error && <p className="error">{error}</p>}
			</form>
		</div>
	);
};

export default Register;
