package isoduration

import (
	"errors"
	"fmt"
	"reflect"
)

// is checks type matching
var is = func(err, target error) bool {
	if err == nil {
		return false
	}

	if reflect.TypeOf(err) != reflect.TypeOf(target) {
		return false
	}

	return reflect.DeepEqual(err, target)
}

// IsNotIsoFormatError occurs when parsing a string in ISO 8601 duration format, the error cannot be clearly identified
var IsNotIsoFormatError = errors.New("incorrect ISO 8601 duration format")

// TimeIsEmptyError occurs when a time designator is found in a line, but its value is not found
// For example: P10YT
var TimeIsEmptyError = errors.New("incorrect ISO 8601 T duration format, designator T found, but value is empty")

// PeriodIsEmptyError occurs when a period designator is found in a line, but its value is not found
// For example: P
var PeriodIsEmptyError = errors.New("incorrect ISO 8601 P duration format, designator P found, but value is empty")

// IncorrectIsoFormatError occurs when a token is found in a string that cannot be converted to the float64 type
// For example: P10,5Y
type IncorrectIsoFormatError struct {
	text string
	in   string
}

// Error defines error output
func (i *IncorrectIsoFormatError) Error() string {
	return fmt.Sprintf(i.text, i.in)
}

// Is checks for object matching
func (i *IncorrectIsoFormatError) Is(err error) bool {
	return is(i, err)
}

// NewIncorrectIsoFormatError creates new IncorrectIsoFormatError
func NewIncorrectIsoFormatError(in string) *IncorrectIsoFormatError {
	return &IncorrectIsoFormatError{"incorrect ISO 8601 duration format, invalid tokens %s", in}
}

// IncorrectDesignatorError occurs when an unknown designator is encountered in the line.
// For example: PT10Y or P1V
type IncorrectDesignatorError struct {
	text       string
	designator rune
	state      rune
}

// Error defines error output
func (i *IncorrectDesignatorError) Error() string {
	return fmt.Sprintf(i.text, i.state, i.designator)
}

// Is checks for object matching
func (i *IncorrectDesignatorError) Is(err error) bool {
	return is(i, err)
}

// NewIncorrectDesignatorError creates new IncorrectDesignatorError
func NewIncorrectDesignatorError(state, designator rune) *IncorrectDesignatorError {
	return &IncorrectDesignatorError{"incorrect ISO 8601 duration %c format, invalid designator %c", designator, state}
}

// DesignatorNotFoundError occurs when there is no designator in the line after the token
// For example P10T10H.
type DesignatorNotFoundError struct {
	text  string
	state rune
	after string
}

// Error defines error output
func (i *DesignatorNotFoundError) Error() string {
	return fmt.Sprintf(i.text, i.state, i.after)
}

// Is checks for object matching
func (i *DesignatorNotFoundError) Is(err error) bool {
	return is(i, err)
}

// NewDesignatorNotFoundError creates new DesignatorNotFoundError
func NewDesignatorNotFoundError(state rune, after string) *DesignatorNotFoundError {
	return &DesignatorNotFoundError{"incorrect ISO 8601 duration %c format, designator not found after token %s", state, after}
}

// DesignatorValueNotFoundError occurs when no value is found for the designator
// For example: PT1HM
type DesignatorValueNotFoundError struct {
	text       string
	state      rune
	designator rune
}

// Error defines error output
func (i *DesignatorValueNotFoundError) Error() string {
	return fmt.Sprintf(i.text, i.state, i.designator)
}

// Is checks for object matching
func (i *DesignatorValueNotFoundError) Is(err error) bool {
	return is(i, err)
}

// NewDesignatorValueNotFoundError creates new DesignatorValueNotFoundError
func NewDesignatorValueNotFoundError(state, designator rune) *DesignatorValueNotFoundError {
	return &DesignatorValueNotFoundError{"incorrect ISO 8601 duration %c format, %c designator's value not found", state, designator}
}

// DesignatorMetError occurs when this designator has already been processed previously.
// Affect: if you write string like PT10H0M10M it defines value as 10
// For example: P1Y5Y
type DesignatorMetError struct {
	text       string
	designator rune
}

// Error defines error output
func (i *DesignatorMetError) Error() string {
	return fmt.Sprintf(i.text, i.designator)
}

// Is checks for object matching
func (i *DesignatorMetError) Is(err error) bool {
	return is(i, err)
}

// NewDesignatorMetError creates new DesignatorMetError
func NewDesignatorMetError(designator rune) *DesignatorMetError {
	return &DesignatorMetError{"incorrect ISO 8601 duration format, the designator %c has already been processed", designator}
}
