package isoduration

import (
	"strconv"
	"time"
)

// supported designators
var (
	periodDesignators = [4]rune{YEAR, MONTH, WEEK, DAY}
	timeDesignators   = [3]rune{HOUR, MINUTE, SECOND}
)

// designators methods
var (
	periodDesignatorsDef = map[rune]periodDesignatorFunc{
		YEAR: {
			get: func(d *PeriodDuration) time.Duration {
				return time.Duration(float64(time.Hour) * DayHours * YearDays * d.years)
			},
			set:    func(d *PeriodDuration, v float64) { d.years = v },
			string: func(d *PeriodDuration) string { return strconv.FormatFloat(d.years, 'f', -1, 64) + "Y" },
			checkSet: func(d *PeriodDuration) bool {
				if d.years == 0 {
					return false
				} else {
					return true
				}
			},
		},
		MONTH: {
			get: func(d *PeriodDuration) time.Duration {
				return time.Duration(float64(time.Hour) * DayHours * MonthDays * d.months)
			},
			set:    func(d *PeriodDuration, v float64) { d.months = v },
			string: func(d *PeriodDuration) string { return strconv.FormatFloat(d.months, 'f', -1, 64) + "M" },
			checkSet: func(d *PeriodDuration) bool {
				if d.months == 0 {
					return false
				} else {
					return true
				}
			},
		},
		DAY: {
			get:    func(d *PeriodDuration) time.Duration { return time.Duration(float64(time.Hour) * DayHours * d.days) },
			set:    func(d *PeriodDuration, v float64) { d.days = v },
			string: func(d *PeriodDuration) string { return strconv.FormatFloat(d.days, 'f', -1, 64) + "D" },
			checkSet: func(d *PeriodDuration) bool {
				if d.days == 0 {
					return false
				} else {
					return true
				}
			},
		},
		WEEK: {
			get: func(d *PeriodDuration) time.Duration {
				return time.Duration(float64(time.Hour*DayHours*WeekDays) * d.weeks)
			},
			set:    func(d *PeriodDuration, v float64) { d.weeks = v },
			string: func(d *PeriodDuration) string { return strconv.FormatFloat(d.weeks, 'f', -1, 64) + "W" },
			checkSet: func(d *PeriodDuration) bool {
				if d.weeks == 0 {
					return false
				} else {
					return true
				}
			},
		},
	}

	timeDesignatorsDef = map[rune]timeDesignatorFunc{
		HOUR: {
			get:    func(td *TimeDuration) time.Duration { return time.Duration(float64(time.Hour) * td.hours) },
			set:    func(td *TimeDuration, v float64) { td.hours = v },
			string: func(td *TimeDuration) string { return strconv.FormatFloat(td.hours, 'f', -1, 64) + "H" },
			checkSet: func(td *TimeDuration) bool {
				if td.hours == 0 {
					return false
				} else {
					return true
				}
			},
		},
		MINUTE: {
			get:    func(td *TimeDuration) time.Duration { return time.Duration(float64(time.Minute) * td.minutes) },
			set:    func(td *TimeDuration, v float64) { td.minutes = v },
			string: func(td *TimeDuration) string { return strconv.FormatFloat(td.minutes, 'f', -1, 64) + "M" },
			checkSet: func(td *TimeDuration) bool {
				if td.minutes == 0 {
					return false
				} else {
					return true
				}
			},
		},
		SECOND: {
			get:    func(td *TimeDuration) time.Duration { return time.Duration(float64(time.Second) * td.seconds) },
			set:    func(td *TimeDuration, v float64) { td.seconds = v },
			string: func(td *TimeDuration) string { return strconv.FormatFloat(td.seconds, 'f', -1, 64) + "S" },
			checkSet: func(td *TimeDuration) bool {
				if td.seconds == 0 {
					return false
				} else {
					return true
				}
			},
		},
	}
)

// designatorFunc interface for parser
type designatorFunc interface {
	periodDesignatorFunc | timeDesignatorFunc
}

// periodDesignatorFunc defines the available methods available for working with period designators
type periodDesignatorFunc struct {
	get      func(*PeriodDuration) time.Duration
	string   func(*PeriodDuration) string
	set      func(*PeriodDuration, float64)
	checkSet func(*PeriodDuration) bool
}

// timeDesignatorFunc defines the available methods available for working with time designators
type timeDesignatorFunc struct {
	get      func(*TimeDuration) time.Duration
	string   func(*TimeDuration) string
	set      func(*TimeDuration, float64)
	checkSet func(duration *TimeDuration) bool
}
