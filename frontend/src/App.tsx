import { KeyboardEvent, useState } from "react";
import { getLocalTimeZone, parseTime, today } from "@internationalized/date";
import {
	addToast,
	Button,
	DatePicker,
	DateRangePicker,
	DateValue,
	Listbox,
	ListboxItem,
	Modal,
	ModalBody,
	ModalContent,
	ModalFooter,
	ModalHeader,
	RangeValue,
	Spinner,
	TimeInput,
	TimeInputValue,
	Tooltip,
} from "@heroui/react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faAdd, faX } from "@fortawesome/free-solid-svg-icons";

import { GenerateThumbnails } from "../wailsjs/go/main/App";
import NavigationBar from "./NavigationBar";
import { svgStore } from "./zustand";

function App() {
	const [showSpinner, setShowSpinner] = useState(false);

	const svg = svgStore((state) => state.SVG);

	const [dateRange, setDateRange] = useState<RangeValue<DateValue> | null>({
		start: today(getLocalTimeZone()),
		end: today(getLocalTimeZone()).add({ weeks: 1 }),
	});
	const [time, setTime] = useState<TimeInputValue | null>(
		parseTime("10:00:00"),
	);

	const [customDates, setCustomDates] = useState<DateValue[]>([]);
	const [dateInput, setDateInput] = useState<DateValue | null>(null);

	function addCustomDate(e?: KeyboardEvent) {
		if (e === undefined || e.key === "Enter") {
			if (dateInput) {
				setCustomDates([...customDates, dateInput]);
			}
		}
	}

	async function generate() {
		if (dateRange && time) {
			setShowSpinner(true);

			let result: number = 0;

			try {
				result = await GenerateThumbnails({
					from: dateRange.start.toString(),
					to: dateRange.end.toString(),
					time: time.toString(),
					customDates: customDates.map((dt) => dt.toString()),
				});
			} catch (err) {
				addToast({
					title: "Can't generate thumbnails",
					description: `${err}`,
					color: "danger",
				});
			} finally {
				if (result > 0) {
					addToast({
						title: "Export successful",
						description: `Exported ${result} thumbnails`,
						color: "success",
					});
				}

				setShowSpinner(false);
			}
		}
	}

	return (
		<div className="flex h-screen w-screen flex-col">
			<NavigationBar exportDisabled={svg.length === 0} onGenerate={generate} />

			<div className="m-4 flex flex-1 justify-between gap-2">
				<div className="flex flex-col">
					<h3 className="text-lg font-semibold">Preview</h3>
					<div className="my-auto flex aspect-video h-[80vh] items-center justify-center overflow-hidden rounded-small border-small border-default-200">
						{svg.length > 0 ? (
							<img
								src={`data:image/svg+xml;utf-8,${encodeURIComponent(svg)}`}
							/>
						) : (
							<div>Open a template</div>
						)}
					</div>
				</div>
				<div className="flex flex-1 flex-col gap-2">
					<div className="flex gap-2">
						<DateRangePicker
							aria-label="creation range"
							label="Creation range"
							labelPlacement="inside"
							showMonthAndYearPickers
							value={dateRange}
							onChange={setDateRange}
						/>
						<TimeInput
							className="w-max"
							label="Time"
							value={time}
							onChange={setTime}
						/>
					</div>

					<h3>Custom Dates</h3>
					<div className="flex-1 rounded-small border-small border-default-200 px-1 py-2">
						<Listbox
							aria-label="custom dates"
							variant="light"
							emptyContent="No custom dates"
						>
							{customDates.map((item, index) => (
								<ListboxItem
									key={item.toString() + index.toString()}
									endContent={
										<Button
											aria-label={`delete custom date ${item.toString()}`}
											isIconOnly
											variant="light"
											size="sm"
											onPress={() =>
												setCustomDates(customDates.toSpliced(index, 1))
											}
										>
											<FontAwesomeIcon size="xs" icon={faX} />
										</Button>
									}
								>
									{item.toString()}
								</ListboxItem>
							))}
						</Listbox>
					</div>
					<div className="flex items-center gap-2">
						<Tooltip content="add date">
							<Button isIconOnly onPress={() => addCustomDate()}>
								<FontAwesomeIcon icon={faAdd} />
							</Button>
						</Tooltip>
						<DatePicker
							label="custom date"
							value={dateInput}
							onChange={setDateInput}
							onKeyUp={addCustomDate}
						/>
					</div>
				</div>
			</div>
			<Modal isOpen={showSpinner} backdrop="blur" hideCloseButton>
				<ModalContent>
					<ModalHeader>
						<h1>Exporting Thumbnails</h1>
					</ModalHeader>
					<ModalBody>
						<Spinner />
					</ModalBody>
					<ModalFooter />
				</ModalContent>
			</Modal>
		</div>
	);
}

export default App;
