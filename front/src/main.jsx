import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import App from "./App";
import Navbar from "./components/NavBar";
import Login from "./pages/Login";
import Register from "./pages/Register";
import FileManager from "./pages/FileManager";
import PrivateRoute from "./components/PrivateRoute";
import NotFound from "./pages/NotFound";
import VantaBackground from "./components/VantaBackground";
import "./styles/global.scss";

ReactDOM.createRoot(document.getElementById("root")).render(
	<React.StrictMode>
		<Router>
			<VantaBackground />
			
			<Navbar />

			<main>
				<Routes>
					<Route path="/" element={<App />} />
					<Route path="/login" element={<Login />} />
					<Route path="/register" element={<Register />} />

					{/* Protected routes */}
					<Route element={<PrivateRoute />}>
						<Route path="/files" element={<FileManager />} />
					</Route>

					{/* 404 Not Found */}
					<Route path="*" element={<NotFound />} />
				</Routes>
			</main>
		</Router>
	</React.StrictMode>
);
