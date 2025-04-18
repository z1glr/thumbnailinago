package locales

type Months []string
type Days []string

type LanguageLocale struct {
	Months
	Days
}

type dataStruct struct {
	German LanguageLocale
}

var data = dataStruct{
	German: LanguageLocale{
		Months{
			"Januar",
			"Februar",
			"März",
			"April",
			"Mai",
			"Juni",
			"Juli",
			"August",
			"September",
			"Oktober",
			"November",
			"Dezember",
		},
		Days{
			"Sonntag",
			"Montag",
			"Dienstag",
			"Mittwoch",
			"Donnerstag",
			"Freitag",
			"Samstag",
		},
	},
}
