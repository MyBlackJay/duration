package isoduration

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

type T struct {
	D *Duration `json:"duration"`
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		input  T
		result []byte
	}{
		{
			input:  T{D: NewDuration(1, 1.5, 7, 1, 1, 60, 30, false)},
			result: []byte(`{"duration":"P1Y1.5M1W7DT1H60M30S"}`),
		},
		{
			input:  T{D: NewDuration(1, 1.5, 7, 1, 1, 60, 30, true)},
			result: []byte(`{"duration":"-P1Y1.5M1W7DT1H60M30S"}`),
		},
	}

	for i, v := range tests {
		body, err := json.Marshal(v.input)

		switch {
		case err == nil && string(body) == string(v.result):
			t.Logf("Test %d (input: %s) completed successfully", i, v.input)
		default:
			t.Errorf("Test %d (input: %s) failed. Expected: %s. Result: %s", i, v.input, v.result, body)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input   []byte
		result  T
		isError bool
	}{
		{
			input:   []byte(`{"duration": "P1Y1.5M1W7DT1H60M30S"}`),
			result:  T{D: NewDuration(1, 1.5, 7, 1, 1, 60, 30, false)},
			isError: false,
		},
		{
			input:   []byte(`{"duration": "-P1Y1.5M1W7DT1H60M30S"}`),
			result:  T{D: NewDuration(1, 1.5, 7, 1, 1, 60, 30, true)},
			isError: false,
		},
		{
			input:   []byte(`{"duration": "PP1Y1.5M1W7DT1H60M30S"}`),
			result:  T{D: NewDuration(1, 1.5, 7, 1, 1, 60, 30, false)},
			isError: true,
		},
		{
			input:   []byte(`{"duration": "P1Y1.5M1W7DT"}`),
			result:  T{D: NewDuration(1, 1.5, 7, 1, 1, 60, 30, false)},
			isError: true,
		},
	}
	for i, v := range tests {
		tr := &T{}
		err := json.Unmarshal(v.input, tr)
		switch {
		case (err != nil && v.isError) || (err == nil && reflect.DeepEqual(tr, &v.result)):
			t.Logf("Test %d (input: %s) completed successfully", i, v.input)
		default:
			t.Errorf("Test %d (input: %s) failed. Expected: %s. Result: %s", i, v.input, v.result.D, v.result.D)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input  *Duration
		result string
	}{
		{
			input:  NewDuration(2, 6, 32, 2, 1, 65, 30, false),
			result: "P2Y6M2W32DT1H65M30S",
		},
		{
			input:  NewDuration(2, 6, 32.6, 2, 1, 65.6, 5, true),
			result: "-P2Y6M2W32.6DT1H65.6M5S",
		},
		{
			input:  NewDuration(0, 0, 0, 0, 0, 0, 0, false),
			result: "PT0S",
		},
		{
			input:  NewDuration(0, 0, 0, 0, 0, 0, 0, true),
			result: "PT0S",
		},
	}
	for i, v := range tests {
		switch {
		case v.input.String() == v.result:
			t.Logf("Test %d (iso duration: %s) completed successfully", i, v.input)
		default:
			t.Errorf("Test %d (input: %s) failed. Expected: %s. Result: %s", i, v.input, v.input, v.result)
		}
	}

}

func TestNewFromTimeDuration(t *testing.T) {
	tests := []struct {
		input          time.Duration
		stringResult   string
		durationResult time.Duration
	}{
		{
			time.Duration(
				(float64(time.Hour) * DayHours * YearDays * 2) +
					(float64(time.Hour) * DayHours * MonthDays * 6) +
					(float64(time.Hour) * DayHours * WeekDays * 2) +
					(float64(time.Hour) * DayHours * 32) +
					(float64(time.Hour) * 1) +
					(float64(time.Minute) * 65) +
					(float64(time.Second) * 30),
			),
			"P2Y7M2W2DT2H5M30S",
			NewDuration(2, 6, 32, 2, 1, 65, 30, false).ToTimeDuration(),
		},
		{
			-1 * time.Duration(
				(float64(time.Hour)*DayHours*YearDays*1)+
					(float64(time.Hour)*DayHours*MonthDays*3)+
					(float64(time.Hour)*DayHours*WeekDays*1)+
					(float64(time.Hour)*DayHours*1)+
					(float64(time.Hour)*1)+
					(float64(time.Minute)*65)+
					(float64(time.Second)*30),
			),
			"-P1Y3M1W1DT2H5M30S",
			NewDuration(1, 3, 1, 1, 1, 65, 30, true).ToTimeDuration(),
		},
		{
			time.Duration(0),
			"PT0S",
			NewDuration(0, 0, 0, 0, 0, 0, 0, false).ToTimeDuration(),
		},
	}

	for i, v := range tests {
		dur := NewFromTimeDuration(v.input)

		switch {
		case dur.String() == v.stringResult && dur.ToTimeDuration() == v.durationResult:
			t.Logf("Test %d (iso duration: %s) completed successfully", i, dur)
		default:
			t.Errorf("Test %d (input: %s) failed. Expected: %s. Result: %s", i, v.stringResult, dur.ToTimeDuration(), v.durationResult)
		}
	}
}

func TestToTimeDuration(t *testing.T) {
	tests := []struct {
		object *Duration
		result time.Duration
	}{
		{
			NewDuration(1, 1, 1, 1, 1, 65, 30, false),
			time.Duration(
				(float64(time.Hour) * DayHours * YearDays * 1) +
					(float64(time.Hour) * DayHours * MonthDays * 1) +
					(float64(time.Hour) * DayHours * WeekDays * 1) +
					(float64(time.Hour) * DayHours * 1) +
					(float64(time.Hour) * 1) +
					(float64(time.Minute) * 65) +
					(float64(time.Second) * 30),
			),
		},
		{
			NewDuration(1, 1, 1, 1, 1, 65, 30, true),
			-1 * time.Duration(
				(float64(time.Hour)*DayHours*YearDays*1)+
					(float64(time.Hour)*DayHours*MonthDays*1)+
					(float64(time.Hour)*DayHours*WeekDays*1)+
					(float64(time.Hour)*DayHours*1)+
					(float64(time.Hour)*1)+
					(float64(time.Minute)*65)+
					(float64(time.Second)*30),
			),
		},
		{
			NewDuration(0, 0, 0, 0, 0, 0, 0, false),
			time.Duration(0),
		},
	}

	for i, v := range tests {
		r := v.object.ToTimeDuration()
		switch r {
		case v.result:
			t.Logf("Test %d (iso duration: %s) completed successfully", i, v.object)
		default:
			t.Errorf("Test %d (input: %s) failed. Expected: %s. Result: %s", i, v.object, r, v.result)
		}
	}
}

func TestParseLetter(t *testing.T) {
	testsTimeDuration := []struct {
		state, designator rune
		nums              string
		tm                *TimeDuration
		desDef            map[rune]timeDesignatorFunc
		result            *TimeDuration
		isError           bool
		err               error
	}{
		{
			state:      TIME,
			designator: SECOND,
			nums:       "0",
			tm:         &TimeDuration{},
			desDef:     timeDesignatorsDef,
			result:     &TimeDuration{seconds: 0},
			isError:    false,
			err:        nil,
		},
		{
			state:      TIME,
			designator: HOUR,
			nums:       "10",
			tm:         &TimeDuration{},
			desDef:     timeDesignatorsDef,
			result:     &TimeDuration{hours: float64(10)},
			isError:    false,
			err:        nil,
		},
		{
			state:      TIME,
			designator: DAY,
			nums:       "10",
			tm:         &TimeDuration{},
			desDef:     timeDesignatorsDef,
			result:     &TimeDuration{hours: float64(10)},
			isError:    true,
			err:        NewIncorrectDesignatorError(TIME, DAY),
		},
		{
			state:      TIME,
			designator: HOUR,
			nums:       "10, 5",
			tm:         &TimeDuration{},
			desDef:     timeDesignatorsDef,
			result:     &TimeDuration{hours: float64(10)},
			isError:    true,
			err:        NewIncorrectIsoFormatError("10, 5"),
		},
		{
			state:      TIME,
			designator: HOUR,
			nums:       "10",
			tm:         &TimeDuration{hours: 1},
			desDef:     timeDesignatorsDef,
			result:     &TimeDuration{hours: float64(1)},
			isError:    true,
			err:        NewDesignatorMetError(HOUR),
		},
	}

	testsPeriodDuration := []struct {
		state, designator rune
		nums              string
		tm                *PeriodDuration
		desDef            map[rune]periodDesignatorFunc
		result            *PeriodDuration
		isError           bool
		err               error
	}{
		{
			state:      PERIOD,
			designator: YEAR,
			nums:       "10.5",
			tm:         &PeriodDuration{},
			desDef:     periodDesignatorsDef,
			result:     &PeriodDuration{years: 10.5},
			isError:    false,
			err:        nil,
		},
		{
			state:      PERIOD,
			designator: HOUR,
			nums:       "10",
			tm:         &PeriodDuration{},
			desDef:     periodDesignatorsDef,
			result:     &PeriodDuration{years: float64(10)},
			isError:    true,
			err:        NewIncorrectDesignatorError(PERIOD, HOUR),
		},
		{
			state:      PERIOD,
			designator: YEAR,
			nums:       "10, 5",
			tm:         &PeriodDuration{},
			desDef:     periodDesignatorsDef,
			result:     nil,
			isError:    true,
			err:        NewIncorrectIsoFormatError("10, 5"),
		},
		{
			state:      PERIOD,
			designator: YEAR,
			nums:       "10",
			tm:         &PeriodDuration{years: 10},
			desDef:     periodDesignatorsDef,
			result:     nil,
			isError:    true,
			err:        NewDesignatorMetError(YEAR),
		},
	}
	t.Run("testsParseLetterForTimeDuration", func(t *testing.T) {
		for i, v := range testsTimeDuration {
			err := parseLetter(v.state, v.designator, v.nums, v.tm, v.desDef)

			switch {
			case err != nil && v.isError && errors.Is(err, v.err):
				t.Logf("Test %d (state: %c, designator: %c, nums: %s) completed successfully", i, v.state, v.designator, v.nums)
			case err == nil && v.result.hours == v.tm.hours:
				t.Logf("Test %d (state: %c, designator: %c, nums: %s) completed successfully", i, v.state, v.designator, v.nums)
			default:
				t.Errorf("Test %d (state: %c, designator: %c, nums: %s) failed.", i, v.state, v.designator, v.nums)
			}
		}
	})

	t.Run("testsParseLetterForPeriodDuration", func(t *testing.T) {
		for i, v := range testsPeriodDuration {
			err := parseLetter(v.state, v.designator, v.nums, v.tm, v.desDef)

			switch {
			case err != nil && v.isError && errors.Is(err, v.err):
				t.Logf("Test %d (state: %c, designator: %c, nums: %s) completed successfully", i, v.state, v.designator, v.nums)
			case err == nil && v.result.years == v.tm.years:
				t.Logf("Test %d (state: %c, designator: %c, nums: %s) completed successfully", i, v.state, v.designator, v.nums)
			default:
				t.Errorf("Test %d (state: %c, designator: %c, nums: %s) failed.", i, v.state, v.designator, v.nums)
			}
		}
	})

}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input   string
		result  *Duration
		isError bool
		err     error
	}{
		{
			input:   "PT0S",
			result:  NewDuration(0, 0, 0, 0, 0, 0, 0, false),
			isError: false,
			err:     nil,
		},
		{
			input:   "P1Y1.5M1W7DT1H60M30S",
			result:  NewDuration(1, 1.5, 7, 1, 1, 60, 30, false),
			isError: false,
			err:     nil,
		},
		{
			input:   "PT1H60M30S",
			result:  NewDuration(0, 0, 0, 0, 1, 60, 30, false),
			isError: false,
			err:     nil,
		},
		{
			input:   "P1Y1.5M1W7D",
			result:  NewDuration(1, 1.5, 7, 1, 0, 0, 0, false),
			isError: false,
			err:     nil,
		},
		{
			input:   "-P1Y1.5M1W7DT1H60M30S",
			result:  NewDuration(1, 1.5, 7, 1, 1, 60, 30, true),
			isError: false,
			err:     nil,
		},
		{
			input:   "PPP",
			result:  nil,
			isError: true,
			err:     IsNotIsoFormatError,
		},
		{
			input:   "P",
			result:  nil,
			isError: true,
			err:     PeriodIsEmptyError,
		},
		{
			input:   "PT",
			result:  nil,
			isError: true,
			err:     TimeIsEmptyError,
		},
		{
			input:   "PT10",
			result:  nil,
			isError: true,
			err:     NewDesignatorNotFoundError(TIME, "10"),
		},
		{
			input:   "P20T10H",
			result:  nil,
			isError: true,
			err:     NewDesignatorNotFoundError(PERIOD, "20"),
		},
		{
			input:   "PT10H9H",
			result:  nil,
			isError: true,
			err:     NewDesignatorMetError(HOUR),
		},
		{
			input:   "PT1HM",
			result:  nil,
			isError: true,
			err:     NewDesignatorValueNotFoundError(TIME, MINUTE),
		},
		{
			input:   "P1H",
			result:  nil,
			isError: true,
			err:     NewIncorrectDesignatorError(PERIOD, HOUR),
		},
	}

	for i, v := range tests {
		result, err := ParseDuration(v.input)

		switch {
		case err != nil && v.isError && errors.Is(err, v.err):
			t.Logf("Test %d (input: %s) completed successfully", i, v.input)
		case err == nil && reflect.DeepEqual(result, v.result):
			t.Logf("Test %d (input: %s) completed successfully", i, v.input)
		default:
			t.Errorf("Test %d (input: %s) failed. Expected: %s. Result: %s", i, v.input, v.result, result)
		}
	}
}
