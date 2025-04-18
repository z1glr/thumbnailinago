import {
	faCog,
	faFloppyDisk,
	faFolderOpen,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
	addToast,
	Button,
	ButtonGroup,
	Navbar,
	NavbarContent,
	NavbarItem,
	Tooltip,
} from "@heroui/react";
import Settings from "./Settings";
import { useState } from "react";
import { OpenTemplate, RefreshPreview } from "../wailsjs/go/main/App";
import { svgStore } from "./zustand";
import { main } from "../wailsjs/go/models";

export default function NavigationBar({
	exportDisabled,
	onGenerate,
}: {
	exportDisabled: boolean;
	onGenerate: () => void;
}) {
	const [showSettings, setShowSettings] = useState(false);

	const templateName = svgStore((state) => state.name);
	const setSVG = svgStore((state) => state.setSVG);

	async function openTemplate() {
		let result: main.FrontendTemplate;

		try {
			result = await OpenTemplate();
		} catch (err) {
			addToast({
				title: "Can't open template",
				description: `${err}`,
				color: "danger",
			});

			return;
		}

		if (result.SVG.length > 0) {
			setSVG(result);
		}
	}

	async function refreshPreview() {
		let result: main.FrontendTemplate;
		try {
			result = await RefreshPreview();
		} catch (err) {
			addToast({
				title: "Can't open template",
				description: `${err}`,
				color: "danger",
			});

			return;
		}

		if (result.SVG.length > 0) {
			setSVG(result);
		}
	}

	return (
		<>
			<Navbar
				classNames={{
					wrapper: "max-w-full",
				}}
			>
				<NavbarContent>
					<ButtonGroup>
						<Button isIconOnly onPress={openTemplate}>
							<Tooltip content="open template">
								<NavbarItem>
									<FontAwesomeIcon icon={faFolderOpen} />
								</NavbarItem>
							</Tooltip>
						</Button>
						<Button isIconOnly>
							<Tooltip content="open settings">
								<NavbarItem>
									<Button isIconOnly onPress={() => setShowSettings(true)}>
										<FontAwesomeIcon icon={faCog} />
									</Button>
								</NavbarItem>
							</Tooltip>
						</Button>
					</ButtonGroup>
				</NavbarContent>

				<NavbarContent justify="center">
					<NavbarItem className="font-mono italic">{templateName}</NavbarItem>
				</NavbarContent>

				<NavbarContent justify="end">
					<NavbarItem>
						<Tooltip content="generate thumbnails">
							<Button
								aria-label="generate thumbnails"
								isIconOnly
								onPress={onGenerate}
								disabled={exportDisabled}
							>
								<FontAwesomeIcon icon={faFloppyDisk} />
							</Button>
						</Tooltip>
					</NavbarItem>
				</NavbarContent>
			</Navbar>
			;
			<Settings
				isOpen={showSettings}
				onOpenChange={(state) => {
					setShowSettings(state);
					if (!state) refreshPreview();
				}}
			/>
		</>
	);
}
