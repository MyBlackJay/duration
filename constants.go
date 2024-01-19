package isoduration

const (
	// PERIOD is constant that defines period designators
	PERIOD = 'P'
	// TIME is constant that defines time designators
	TIME = 'T'

	// YEAR is constant that defines year designators
	YEAR = 'Y'
	// MONTH is constant that defines month designators
	MONTH = 'M'
	// WEEK is constant that defines week designators
	WEEK = 'W'
	// DAY is constant that defines day designators
	DAY = 'D'
	// HOUR is constant that defines hour designators
	HOUR = 'H'
	// MINUTE is constant that defines minute designators
	MINUTE = 'M'
	// SECOND is constant that defines second designators
	SECOND = 'S'

	// DayHours is constant that defines number of hours in a day
	DayHours = 24
	// WeekDays is constant that defines number of days in a week
	WeekDays = 7
	// YearDays is constant that defines number of days in a year
	YearDays = 365
	// MonthDays is constant that defines number of days in a month
	MonthDays = YearDays / 12
)
