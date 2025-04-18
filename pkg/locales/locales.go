package locales

import "time"

type Locale struct {
	Year, Day      int
	Month, Weekday string
}

func CreateEnglish(d *time.Time) Locale {
	return Locale{
		Year:    d.Year(),
		Day:     d.Day(),
		Month:   d.Format("January"),
		Weekday: d.Weekday().String(),
	}
}

func CreateGerman(d *time.Time) Locale {
	return Locale{
		Year:    d.Year(),
		Day:     d.Day(),
		Month:   data.German.Months[d.Month()-1],
		Weekday: data.German.Days[d.Weekday()],
	}
}

func Create(l string, d *time.Time) Locale {
	switch l {
	case "en":
		return CreateEnglish(d)
	case "de":
		return CreateGerman(d)
	default:
		return Locale{}
	}
}
