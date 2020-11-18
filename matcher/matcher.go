package matcher

import (
	"fmt"
	"github.com/madelyne-io/madelyne/matcher/vm"
	"regexp"
	"strings"
)

var (
	ErrUnhandledType  = fmt.Errorf("Unhandled type")
	ErrInvalidPattern = fmt.Errorf("Invalid pattern.")
	ErrInvalidValue   = fmt.Errorf("Provided value is not what was expected.")
	ErrNotBool        = fmt.Errorf("Provided value is not bool type.")
	ErrNotNumber      = fmt.Errorf("Provided value is not number type.")
	ErrNotString      = fmt.Errorf("Provided value is not string type.")
	ErrNotUuid        = fmt.Errorf("Provided value is not uuid type.")

	ErrUnhandledFunction = fmt.Errorf("Unhandled function.")
	ErrInvalidFunctions  = fmt.Errorf("Functions are not written correctly.")
)

func Match(value interface{}, expected interface{}) error {
	expectedAsString, ok := expected.(string)
	if !ok {
		return matchValue(value, expected)
	}
	return matchPatter(value, expectedAsString)
}

func matchValue(value interface{}, expected interface{}) error {
	if value == expected {
		return nil
	}
	return ErrInvalidValue
}

func matchPatter(value interface{}, pattern string) error {
	splitted := strings.Split(pattern, "@")
	if len(splitted) < 3 {
		return matchValue(value, pattern)
	}

	switch splitted[1] {
	case "string":
		return matchString(value, splitted[2])
	case "number", "double", "integer":
		return matchNumber(value, splitted[2])
	case "boolean":
		return matchBool(value)
	case "uuid":
		return matchUuid(value)
	}
	return fmt.Errorf("%w Got: %s", ErrInvalidPattern, pattern)
}

func matchString(value interface{}, program string) error {
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	if len(program) == 0 {
		return nil
	}
	match, err := vm.BuildProgramMatcher(program, map[string]func(value interface{}, args []interface{}) error{
		"startsWith":  fn_string_startsWith,
		"endsWith":    fn_string_endsWith,
		"contains":    fn_string_contains,
		"notContains": fn_string_notContains,
		"isUrl":       fn_string_isUrl,
		"isDateTime":  fn_string_isDateTime,
		"isEmail":     fn_string_isEmail,
		"isEmpty":     fn_string_isEmpty,
		"isNotEmpty":  fn_string_isNotEmpty,
		"matchRegex":  fn_string_matchRegex,
		"oneOf":       fn_oneOf,
	})
	if err != nil {
		return err
	}
	return match(value)
}

func matchNumber(value interface{}, program string) error {
	_, ok := value.(float64)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotNumber, value)
	}
	if len(program) == 0 {
		return nil
	}
	match, err := vm.BuildProgramMatcher(program, map[string]func(value interface{}, args []interface{}) error{
		"greaterThan": fn_number_greaterThan,
		"lowerThan":   fn_number_lowerThan,
		"oneOf":       fn_oneOf,
	})
	if err != nil {
		return err
	}
	return match(value)
}

func matchBool(value interface{}) error {
	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotBool, value)
	}
	return nil
}

func matchUuid(value interface{}) error {
	valueAsString, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotUuid, value)
	}
	uuidRegexp := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	if !uuidRegexp.MatchString(valueAsString) {
		return fmt.Errorf("%w Got: %v", ErrNotUuid, value)
	}
	return nil
}
