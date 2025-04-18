import { heroui } from "@heroui/react";

export default {
	content: [
		"./index.html",
		"./src/**/*.{tsx,jsx,html,js,ts}",
		"./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
	],
	theme: {
		extend: {},
	},
	darkMode: "selector",
	plugins: [heroui()],
};
