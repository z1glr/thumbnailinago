import {
	addToast,
	Checkbox,
	CheckboxGroup,
	Input,
	Modal,
	ModalBody,
	ModalContent,
	ModalFooter,
	ModalHeader,
	Select,
	SelectItem,
	SharedSelection,
} from "@heroui/react";

import { useEffect, useState } from "react";

import { GetSettings, SetSettings } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

const locales: { key: string; label: string }[] = [
	{
		key: "en",
		label: "English",
	},
	{
		key: "de",
		label: "German",
	},
];

const weekdays = [
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
];

export default function Settings({
	isOpen,
	onOpenChange,
}: {
	isOpen: boolean;
	onOpenChange: (isOpen: boolean) => void;
}) {
	const [locale, setLocale] = useState<SharedSelection>(new Set(["en"]));
	const [days, setDays] = useState<string[]>(["Sunday"]);
	const [replacementKey, setReplacementKey] = useState("SUNDAY_DATE");
	const [dateFormat, setDateFormat] = useState(
		"{{.Day}}. {{.Month}} {{.Year}}",
	);

	useEffect(() => {
		const loadSettings = async () => {
			let settings: main.FrontendSettings;
			try {
				settings = await GetSettings();
			} catch (err) {
				addToast({
					title: "Can't load settings from backend",
					description: `${err}`,
					color: "danger",
				});

				return;
			}

			setLocale(new Set([settings.locale]));
			setDays(settings.days);
			setReplacementKey(settings.replacementKey);
			setDateFormat(settings.dateFormat);
		};

		void loadSettings();
	}, []);

	async function closeSettings() {
		try {
			await SetSettings({
				locale: locale.currentKey ?? locales[0].key,
				days: days,
				replacementKey,
				dateFormat,
			});
		} catch (err) {
			addToast({
				title: "Can't save settings",
				description: `${err}`,
				color: "danger",
			});
		}

		onOpenChange(false);
	}

	return (
		<Modal
			isOpen={isOpen}
			onOpenChange={(isOpen) =>
				isOpen ? onOpenChange(isOpen) : closeSettings()
			}
		>
			<ModalContent>
				<ModalHeader>Settings</ModalHeader>

				<ModalBody>
					<Select
						aria-label="select locale"
						items={locales}
						label="Locale"
						selectedKeys={locale}
						onSelectionChange={setLocale}
					>
						{(locale) => (
							<SelectItem key={locale.key}>{locale.label}</SelectItem>
						)}
					</Select>

					<CheckboxGroup
						aria-label="select days"
						value={days}
						onValueChange={setDays}
					>
						{weekdays.map((day) => (
							<Checkbox aria-label={day} key={day} value={day}>
								{day}
							</Checkbox>
						))}
					</CheckboxGroup>

					<Input
						aria-label="replacement-key"
						value={replacementKey}
						onValueChange={setReplacementKey}
						label="Replacement key"
					/>
					<Input
						label="Date format"
						value={dateFormat}
						onValueChange={setDateFormat}
					/>
				</ModalBody>

				<ModalFooter />
			</ModalContent>
		</Modal>
	);
}
