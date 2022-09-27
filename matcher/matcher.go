package matcher

import (
	"fmt"
	"github.com/madelyne-io/madelyne/matcher/vm"
	"reflect"
	"regexp"
	"strconv"
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
	ErrNotSlice       = fmt.Errorf("Provided value is not slice type.")

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

	program := strings.Join(splitted[2:], "@")

	if splitted[0] != "" {
		return matchText(value, pattern)
	}
	switch splitted[1] {
	case "string":
		return matchString(value, program)
	case "number", "double", "integer":
		return matchNumber(value, program)
	case "boolean":
		return matchBool(value)
	case "uuid":
		return matchUuid(value)
	case "array":
		return matchArray(value, program)
	}
	return fmt.Errorf("%w Got: %s", ErrInvalidPattern, pattern)
}

func matchText(value interface{}, pattern string) error {
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("%w Got: %v", ErrNotString, value)
	}
	r := regexp.MustCompile("@([^@]+)@")
	matches := r.FindAllStringSubmatch(pattern, -1)
	p := regexp.QuoteMeta(pattern)
	for _, match := range matches {
		var replace string
		switch match[1] {
		case "string":
			replace = "(.+)"
		case "double":
			replace = "(\\-?[0-9]+[\\.|\\,][0-9]*)"
		case "number":
			replace = "(\\-?[0-9]+[\\.|\\,]?[0-9]*)"
		case "integer":
			replace = "(\\-?[0-9]+)"
		case "uuid":
			replace = "([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})"
		default:
			continue
		}
		p = strings.Replace(p, match[0], replace, -1)
	}
	v := regexp.MustCompile(p)
	if !v.MatchString(value.(string)) {
		return fmt.Errorf("%w Got: %v", ErrInvalidValue, value)
	}
	return nil
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
		"before":      fn_string_before,
		"after":       fn_string_after,
	})
	if err != nil {
		return err
	}
	return match(value)
}

func matchNumber(value interface{}, program string) error {
	_, ok := value.(float64)
	if !ok {
		valueAsString, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w Got: %v", ErrNotNumber, value)
		}
		s, err := strconv.ParseFloat(valueAsString, 64)
		if err != nil {
			return fmt.Errorf("%w Got: %v", ErrNotNumber, value)
		}
		value = s
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

func matchArray(value interface{}, program string) error {
	kind := reflect.ValueOf(value).Kind()
	if kind != reflect.Slice {
		return fmt.Errorf("%w Got: %v", ErrNotSlice, value)
	}
	if len(program) == 0 {
		return nil
	}
	match, err := vm.BuildProgramMatcher(program, map[string]func(value interface{}, args []interface{}) error{
		"repeat": fn_array_repeat,
	})
	if err != nil {
		return err
	}
	return match(value)
}
