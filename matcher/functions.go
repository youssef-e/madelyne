package matcher

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidParameters = fmt.Errorf("Invalid parameters.")

	ErrGreaterThan = fmt.Errorf("Provided value is lower than what it should be.")
	ErrLowerThan   = fmt.Errorf("Provided value is greater than what is should be.")

	ErrNotContains   = fmt.Errorf("Provided string does not contain the expected substring.")
	ErrNotStartsWith = fmt.Errorf("Provided string doesn't start with the expected substring.")
	ErrNotEndsWith   = fmt.Errorf("Provided string doesn't end with the expected substring.")
	ErrNotDateTime   = fmt.Errorf("Provided string is not a datetime.")
	ErrNotEmail      = fmt.Errorf("Provided string is not an email.")
	ErrNotEmpty      = fmt.Errorf("The provided string is not empty.")
	ErrEmpty         = fmt.Errorf("The provided string is empty.")
	ErrNotUrl        = fmt.Errorf("The provided string is not an URL.")
	ErrNotMatchRegex = fmt.Errorf("Provided string doesn't match the regex.")
	ErrInvalidRegex  = fmt.Errorf("Provided regex is invalid.")
	ErrContains      = fmt.Errorf("Provided string contains a value it shouldn't.")

	ErrOneOf = fmt.Errorf("None of the functions provided in OneOf were validated.")
)

func fn_string_startsWith(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 startsWith : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	start, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 startsWith : %w param 0 must a string", ErrInvalidParameters)
	}
	if len(valueAsString) < len(start) {
		return fmt.Errorf("2 startsWith : %w got %s, must start with %s", ErrNotStartsWith, valueAsString, start)
	}
	if valueAsString[:len(start)] != start {
		return fmt.Errorf(
			"3 startsWith : %w got %s, must start with %s", ErrNotStartsWith, valueAsString, start)
	}
	return nil
}

func fn_string_endsWith(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 endsWith : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	end, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 endsWith : %w param 0 must a string", ErrInvalidParameters)
	}
	if len(valueAsString) < len(end) {
		return fmt.Errorf("2 endsWith : %w got %s, must end with %s", ErrNotEndsWith, valueAsString, end)
	}
	if valueAsString[len(valueAsString)-len(end):] != end {
		return fmt.Errorf(
			"3 endsWith : %w got %s, must end with %s", ErrNotEndsWith, valueAsString, end)
	}
	return nil
}

func fn_string_contains(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 contains : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	sub, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 contains : %w param 0 must a string", ErrInvalidParameters)
	}
	if len(valueAsString) < len(sub) {
		return fmt.Errorf("2 contains : %w got %s, must end with %s", ErrNotContains, valueAsString, sub)
	}
	if !strings.Contains(valueAsString, sub) {
		return fmt.Errorf(
			"3 contains : %w got %s, must contains %s", ErrNotContains, valueAsString, sub)
	}
	return nil
}

func fn_string_notContains(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 notContains : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	sub, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 notContains : %w param 0 must a string", ErrInvalidParameters)
	}
	if len(valueAsString) < len(sub) {
		return nil
	}
	if strings.Contains(valueAsString, sub) {
		return fmt.Errorf(
			"3 notContains : %w got %s, must not contains %s", ErrContains, valueAsString, sub)
	}
	return nil
}

func fn_string_isUrl(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 0 {
		return fmt.Errorf("0 isUrl : %w want 0 parameters got %d", ErrInvalidParameters, len(args))
	}
	_, err := url.ParseRequestURI(valueAsString)
	if err != nil {
		return fmt.Errorf("1 isUrl : %w Got: %s", ErrNotUrl, value)
	}
	return nil
}

func fn_string_isDateTime(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 0 {
		return fmt.Errorf("0 isDateTime : %w want 0 parameters got %d", ErrInvalidParameters, len(args))
	}
	_, err := time.Parse(time.RFC3339, valueAsString)
	if err != nil {
		return fmt.Errorf("1 isDateTime : %w Got: %s", ErrNotDateTime, value)
	}
	return nil
}

var emailRegexp = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])")

func fn_string_isEmail(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 0 {
		return fmt.Errorf("0 isEmail : %w want 0 parameters got %d", ErrInvalidParameters, len(args))
	}
	if !emailRegexp.Match([]byte(valueAsString)) {
		return fmt.Errorf("%w Got: %v", ErrNotEmail, value)
	}
	return nil
}

func fn_string_isEmpty(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 0 {
		return fmt.Errorf("0 isEmpty : %w want 0 parameters got %d", ErrInvalidParameters, len(args))
	}
	if len(valueAsString) != 0 {
		return fmt.Errorf("%w Got: %v", ErrNotEmpty, value)
	}
	return nil
}

func fn_string_isNotEmpty(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 0 {
		return fmt.Errorf("0 isNotEmpty : %w want 0 parameters got %d", ErrInvalidParameters, len(args))
	}
	if len(valueAsString) == 0 {
		return fmt.Errorf("%w Got: %v", ErrEmpty, value)
	}
	return nil
}

func fn_string_matchRegex(value interface{}, args []interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 matchRegex : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	re, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 matchRegex : %w param 0 must a string", ErrInvalidParameters)
	}
	compiledRe, err := regexp.Compile(re)
	if err != nil {
		return fmt.Errorf("2 %w : %s,", ErrInvalidRegex, re)
	}
	if !compiledRe.Match([]byte(valueAsString)) {
		return fmt.Errorf("3 %w Regex: %s, string: %s", ErrNotMatchRegex, re, valueAsString)
	}
	return nil
}

func fn_oneOf(value interface{}, args []interface{}) error {
	for i, a := range args {
		if a == nil {
			return nil
		}
		_, ok := a.(error)
		if !ok {
			return fmt.Errorf("oneOf : %d : %w expect a function got %#v", i, ErrInvalidParameters, a)
		}
	}
	return ErrOneOf
}

func fn_number_greaterThan(value interface{}, args []interface{}) error {
	valueAsFloat, ok := value.(float64)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotNumber, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 greaterThan : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	number, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 greaterThan : %w param 0 must a string", ErrInvalidParameters)
	}

	converted, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return fmt.Errorf("2 greaterThan %w cannot parse %v to float", ErrInvalidParameters, number)
	}

	if valueAsFloat < converted {
		return fmt.Errorf("3 greaterThan : %w  got %v < %v", ErrGreaterThan, valueAsFloat, converted)
	}
	return nil
}

func fn_number_lowerThan(value interface{}, args []interface{}) error {
	valueAsFloat, ok := value.(float64)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotNumber, value)
	}
	if len(args) != 1 {
		return fmt.Errorf("0 lowerThan : %w want 1 parameters got %d", ErrInvalidParameters, len(args))
	}
	number, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("1 lowerThan : %w param 0 must a string", ErrInvalidParameters)
	}

	converted, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return fmt.Errorf("2 lowerThan %w cannot parse %v to float", ErrInvalidParameters, number)
	}

	if valueAsFloat > converted {
		return fmt.Errorf("3 lowerThan : %w  got %v > %v", ErrLowerThan, valueAsFloat, converted)
	}
	return nil
}
