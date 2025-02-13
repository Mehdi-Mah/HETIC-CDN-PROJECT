import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { isAuthenticated } from "./utils/auth";

const App = () => {
	const navigate = useNavigate();

	useEffect(() => {
		if (isAuthenticated()) {
			navigate("/files"); // Redirect to file manager if logged in
		} else {
			navigate("/login"); // Otherwise, go to login
		}
	}, [navigate]);

	return null; // No need to render anything, just handling redirection
};

export default App;
