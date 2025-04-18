import { create } from "zustand";
import { main } from "../wailsjs/go/models";

interface SVGStore {
	name: string;
	SVG: string;
	setSVG: (svg: main.FrontendTemplate) => void;
}

export const svgStore = create<SVGStore>()((set) => ({
	name: "",
	SVG: "",
	setSVG(svg) {
		set({
			name: svg.name,
			SVG: svg.SVG,
		});
	},
}));
