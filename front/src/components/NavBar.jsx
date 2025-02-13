import { Link, useNavigate } from "react-router-dom";
import { isAuthenticated, logout } from "../utils/auth";
import "../styles/navbar.scss";

const Navbar = () => {
	const navigate = useNavigate();

	const handleLogout = () => {
		logout();
		navigate("/login");
	};

	return (
		<nav className="navbar">
			<h1 className="logo">GOLANG CDN</h1>
			<div className="nav-links">
				{isAuthenticated() ? (
					<>
						<Link to="/files">Fichiers</Link>
						<button onClick={handleLogout} className="logout-btn">Se déconnecter</button>
					</>
				) : (
					<>
						<Link to="/login">Se connecter</Link>
						<Link to="/register">S'inscrire</Link>
					</>
				)}
			</div>
		</nav>
	);
};

export default Navbar;
