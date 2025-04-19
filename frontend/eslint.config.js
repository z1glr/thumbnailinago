import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import reactPlugin from "eslint-plugin-react";
import { defineConfig } from "eslint/config";
import tailwind from "eslint-plugin-tailwindcss";

export default defineConfig([
	{
		settings: {
			react: {
				version: "detect",
			},
			formComponents: ["Form"],
			linkComponents: ["Link"],
		},
	},
	{
		files: ["src/**/*.{js,mjs,cjs,ts,jsx,tsx}"],
		plugins: { js },
		extends: ["js/recommended"],
	},
	{
		files: ["src/**/*.{js,mjs,cjs,ts,jsx,tsx}"],
		languageOptions: { globals: { ...globals.browser, ...globals.node } },
	},
	{
		ignores: ["wailsjs/**", "dist/**"],
	},
	tseslint.configs.strict,
	reactPlugin.configs.flat.recommended,
	reactPlugin.configs.flat["jsx-runtime"],
	...tailwind.configs["flat/recommended"],
]);
