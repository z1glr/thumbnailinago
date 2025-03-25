import { heroui } from "@heroui/react";

export default {
	content: [
		"./index.html",
		"./src/**/*.{tsx}",
		"./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
	],
	theme: {
		extend: {},
	},
	darkMode: "selector",
	plugins: [heroui()],
};
