import { jwtDecode } from "jwt-decode";

export const isAuthenticated = () => {
	const token = localStorage.getItem("token");
	if (!token) return false;

	const expiry = localStorage.getItem("token_expiry");
	if (expiry && Date.now() >= expiry * 1000) {
		logout();
		return false;
	}

	return true;
};

export const logout = () => {
	localStorage.removeItem("token");
	localStorage.removeItem("token_expiry");
	window.location.href = "/login"; // Redirect to login
};
