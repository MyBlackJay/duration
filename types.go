package isoduration

import (
	"math"
	"time"
)

// PeriodDuration is period duration marks
type PeriodDuration struct {
	years  float64
	months float64
	days   float64
	weeks  float64
}

// TimeDuration is time duration marks
type TimeDuration struct {
	hours   float64
	minutes float64
	seconds float64
}

// Duration is basic duration structure
type Duration struct {
	period     *PeriodDuration
	time       *TimeDuration
	multiplier float64
}

// NewDuration creates new *Duration based on time and period marks
func NewDuration(years, months, days, weeks, hours, minutes, seconds float64, isNegative bool) *Duration {
	multiplier := float64(1)

	if isNegative {
		multiplier = -1
	}

	return &Duration{
		period:     &PeriodDuration{years, months, days, weeks},
		time:       &TimeDuration{hours, minutes, seconds},
		multiplier: multiplier,
	}
}

// NewFromTimeDuration creates new *Duration based on time.Duration
// Affect: This may have some rounding inaccuracies
func NewFromTimeDuration(t time.Duration) *Duration {
	pd := new(PeriodDuration)
	td := new(TimeDuration)
	multiplier := float64(1)

	if t < 0 {
		multiplier = -1
		t *= time.Duration(multiplier)
	}

	if t.Hours() >= DayHours*YearDays {
		pd.years = math.Floor(t.Hours() / (DayHours * YearDays))
		t -= time.Duration(float64(time.Hour) * pd.years * DayHours * YearDays)
	}

	if t.Hours() >= DayHours*MonthDays {
		pd.months = math.Floor(t.Hours() / (DayHours * MonthDays))
		t -= time.Duration(float64(time.Hour) * pd.months * DayHours * MonthDays)
	}

	if t.Hours() >= DayHours*WeekDays {
		pd.weeks = math.Floor(t.Hours() / (DayHours * WeekDays))
		t -= time.Duration(float64(time.Hour) * pd.weeks * DayHours * WeekDays)
	}

	if t.Hours() >= DayHours {
		pd.days = math.Floor(t.Hours() / DayHours)

		t -= time.Duration(float64(time.Hour) * pd.days * DayHours)
	}

	if t.Hours() >= 1 {
		td.hours = math.Floor(t.Hours())
		t -= time.Hour * time.Duration(td.hours)
	}

	if t.Minutes() >= 1 {
		td.minutes = math.Floor(t.Minutes())
		t -= time.Minute * time.Duration(td.minutes)
	}

	td.seconds = t.Seconds()

	return &Duration{
		period:     pd,
		time:       td,
		multiplier: multiplier,
	}
}

// Years returns years from *PeriodDuration
func (d *Duration) Years() float64 {
	return d.period.years * d.multiplier
}

// Months returns months from *PeriodDuration
func (d *Duration) Months() float64 {
	return d.period.months * d.multiplier
}

// Weeks returns weeks from *PeriodDuration
func (d *Duration) Weeks() float64 {
	return d.period.weeks * d.multiplier
}

// Days returns days from *PeriodDuration
func (d *Duration) Days() float64 {
	return d.period.days * d.multiplier
}

// Hours returns hours from *TimeDuration
func (d *Duration) Hours() float64 {
	return d.time.hours * d.multiplier
}

// Minutes returns minutes from *TimeDuration
func (d *Duration) Minutes() float64 {
	return d.time.minutes * d.multiplier
}

// Seconds returns seconds from *TimeDuration
func (d *Duration) Seconds() float64 {
	return d.time.seconds * d.multiplier
}

// ToTimeDuration turns *Duration into time.Duration
func (d *Duration) ToTimeDuration() time.Duration {
	var timeDuration time.Duration

	for _, v := range periodDesignators {
		timeDuration += periodDesignatorsDef[v].get(d.period)
	}
	for _, v := range timeDesignators {
		timeDuration += timeDesignatorsDef[v].get(d.time)
	}

	return timeDuration * time.Duration(d.multiplier)
}

// ToTimeDuration turns *Duration into a string in ISO 8601 duration format
func (d *Duration) String() string {
	prefix := "P"
	period := ""
	tm := ""

	for _, v := range periodDesignators {
		if tmp := periodDesignatorsDef[v].string(d.period); tmp != "0"+string(v) {
			period += tmp
		}
	}

	for _, v := range timeDesignators {
		if tmp := timeDesignatorsDef[v].string(d.time); tmp != "0"+string(v) {
			tm += tmp
		}
	}

	if d.multiplier == -1 && tm != "" && period != "" {
		prefix = "-" + prefix
	}

	switch {
	case tm == "" && period == "":
		tm = "T0S"
	case tm != "":
		tm = "T" + tm
	}

	return prefix + period + tm
}

// FormatTimeDuration represents time.Duration as a string in ISO 8601 duration format
// Affect: This may have some rounding inaccuracies
func FormatTimeDuration(d time.Duration) string {
	return NewFromTimeDuration(d).String()
}

// UnmarshalJSON designed to serialize a string in ISO 8601 duration format to *Duration, defined in user code via the json library
func (d *Duration) UnmarshalJSON(source []byte) error {
	value := string(source)

	if len(value) < 2 {
		return IsNotIsoFormatError
	}

	if value == "null" {
		return nil
	}

	if parsed, err := ParseDuration(value[1 : len(value)-1]); err == nil {
		*d = *parsed
		return nil
	} else {
		return err
	}
}

// MarshalJSON designed to deserialize *Duration to a string in ISO 8601 duration format, defined in user code via the json library
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

// UnmarshalYAML designed to serialize a string in ISO 8601 duration format to *Duration, defined in user code via the gopkg.in/yaml.v3 library
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return IsNotIsoFormatError
	}

	if str == "null" {
		return nil
	}

	if parsed, err := ParseDuration(str); err == nil {
		*d = *parsed
		return nil
	} else {
		return err
	}
}

// MarshalYAML designed to deserialize *Duration to a string in ISO 8601 duration format, defined in user code via the gopkg.in/yaml.v3 library
func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}
