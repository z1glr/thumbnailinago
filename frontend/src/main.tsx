import React from "react";
import { createRoot } from "react-dom/client";
import "./style.css";
import App from "./App";
import { HeroUIProvider, ToastProvider } from "@heroui/react";

const container = document.getElementById("root");

const root = createRoot(container as HTMLElement);

root.render(
	<React.StrictMode>
		<HeroUIProvider>
			<ToastProvider />
			<App />
		</HeroUIProvider>
	</React.StrictMode>,
);
