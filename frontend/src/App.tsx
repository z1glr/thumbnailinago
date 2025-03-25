import { useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { Greet } from "../wailsjs/go/main/App";
import {
	Button,
	ButtonGroup,
	DateRangePicker,
	Navbar,
	NavbarContent,
	NavbarItem,
	TimeInput,
} from "@heroui/react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
	faCog,
	faFile,
	faFileUpload,
	faFloppyDisk,
	faFolderOpen,
} from "@fortawesome/free-solid-svg-icons";

function App() {
	const [resultText, setResultText] = useState(
		"Please enter your name below 👇"
	);
	const [name, setName] = useState("");
	const updateName = (e: any) => setName(e.target.value);
	const updateResultText = (result: string) => setResultText(result);

	function greet() {
		Greet(name).then(updateResultText);
	}

	return (
		<>
			<Navbar>
				<NavbarContent>
					<ButtonGroup>
						<Button isIconOnly>
							<NavbarItem>
								<FontAwesomeIcon icon={faFolderOpen} />
							</NavbarItem>
						</Button>
						<Button isIconOnly>
							<NavbarItem>
								<FontAwesomeIcon icon={faCog} />
							</NavbarItem>
						</Button>
					</ButtonGroup>
				</NavbarContent>

				<NavbarContent justify="center">
					<NavbarItem className="italic font-mono">
						name of the template
					</NavbarItem>
				</NavbarContent>

				<NavbarContent justify="end">
					<NavbarItem>
						<Button isIconOnly>
							<FontAwesomeIcon icon={faFloppyDisk} />
						</Button>
					</NavbarItem>
				</NavbarContent>
			</Navbar>
			<img src={logo} id="logo" alt="logo" />
			<div className="flex flex-col gap-2 w-16">
				<DateRangePicker />
				<TimeInput />
			</div>
		</>
	);
}

export default App;
