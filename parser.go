package isoduration

import (
	"strconv"
	"strings"
	"unicode"
)

// parseLetter checks the found des, converts the types and fills the corresponding field of the pDuration input structure with a value.
// If an error occurs in data processing, returns an error
func parseLetter[DF designatorFunc, D *PeriodDuration | *TimeDuration](state, des rune, nums string, st D, desDef map[rune]DF) error {
	if d, ok := desDef[des]; ok {
		if nums == "" {
			return NewDesignatorValueNotFoundError(state, des)
		}

		if v, err := strconv.ParseFloat(nums, 64); err != nil {
			return NewIncorrectIsoFormatError(nums)
		} else if state == PERIOD {
			tDes := any(d).(periodDesignatorFunc)
			tD := any(st).(*PeriodDuration)

			if !tDes.checkSet(tD) {
				tDes.set(tD, v)

				return nil
			}
		} else if state == TIME {
			tDes := any(d).(timeDesignatorFunc)
			tD := any(st).(*TimeDuration)

			if !tDes.checkSet(tD) {
				tDes.set(tD, v)

				return nil
			}
		}
		return NewDesignatorMetError(des)
	}
	return NewIncorrectDesignatorError(state, des)
}

// parse parses an input string in ISO 8601 duration format without
// a sign and returns *Duration and an error if the string could not be parsed
func parse(duration string, multiplier float64) (*Duration, error) {
	dt := &PeriodDuration{}
	tm := &TimeDuration{}
	state := rune(0)
	fact := rune(0)
	buffer := strings.Builder{}

	for i, char := range duration {
		switch {
		case i == 0 && char == PERIOD:
			state = PERIOD
		case state == PERIOD && char == TIME:
			if buffer.String() != "" {
				return nil, NewDesignatorNotFoundError(state, buffer.String())
			}
			state = TIME
		case char == TIME || char == PERIOD:
			return nil, IsNotIsoFormatError
		case unicode.IsLetter(char):
			var err error

			if state == PERIOD {
				err = parseLetter(state, char, buffer.String(), dt, periodDesignatorsDef)
			} else {
				err = parseLetter(state, char, buffer.String(), tm, timeDesignatorsDef)
			}

			if err != nil {
				return nil, err
			}
			fact = state
			buffer.Reset()
		default:
			buffer.WriteRune(char)
		}
	}

	if state == 0 {
		return nil, IsNotIsoFormatError
	} else if buffer.String() != "" {
		return nil, NewDesignatorNotFoundError(state, buffer.String())
	} else if state == PERIOD && state != fact {
		return nil, PeriodIsEmptyError
	} else if state == TIME && state != fact {
		return nil, TimeIsEmptyError
	}

	return &Duration{
		period:     dt,
		time:       tm,
		multiplier: multiplier,
	}, nil
}

// ParseDuration is the main method for parsing a string in ISO format.
// Returns *Duration and an error if the string could not be parsed
// For example: P10Y5M2W1DT1H1.5M50S or -P10Y5M2W1DT1H1.5M50S
func ParseDuration(duration string) (*Duration, error) {
	if duration == "" {
		return nil, IsNotIsoFormatError
	}

	multiplier := float64(1)

	switch prefix := duration[0]; string(prefix) {
	case "-":
		multiplier = -1
		duration = duration[1:]
	case "+":
		duration = duration[1:]
	}

	return parse(duration, multiplier)
}

// MustParseDuration is the main method for parsing a string in ISO format. Returns *Duration.
// If the string cannot be parsed, it returns a panic
// For example: P10Y5M2W1DT1H1.5M50S or -P10Y5M2W1DT1H1.5M50S
func MustParseDuration(duration string) *Duration {
	d, err := ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return d
}
