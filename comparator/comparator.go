package comparator

import (
	"fmt"
	"github.com/madelyne-io/madelyne/matcher"
	"reflect"
	"regexp"
	"strings"
)

type ComparatorError struct {
	Path []string
	Err  error
}

func ErrorAt(path []string, err error) *ComparatorError {
	return &ComparatorError{
		Path: path,
		Err:  err,
	}
}

func (e *ComparatorError) Error() string {
	return "at '" + strings.Join(e.Path, ".") + "' : " + e.Err.Error()
}
func (e *ComparatorError) Unwrap() error { return e.Err }

var (
	ErrMissingKey        = fmt.Errorf("A key present in validation json is missing in content json")
	ErrExtraKey          = fmt.Errorf("A key should not be present in content json")
	ErrNotMatching       = fmt.Errorf("Content does not match with pattern.")
	ErrTypeNotmatching   = fmt.Errorf("Actual type is not the expected one")
	ErrRessourceNotFound = fmt.Errorf("Cannot load external ressources")
)

type LoadExternalResourceFunction func(path string) (map[string]interface{}, error)
type MatcherFunction func(actual interface{}, expected interface{}) error

type Comparator interface {
	Compare(actual interface{}, expected interface{}) error
	GetCaptured() map[string]interface{}
	Reset()
}

func New(fileLocation string) Comparator {
	return &comparator{
		loadExternalData: getfileLoaderFunc(fileLocation + "/responses/"),
		valueMatcher:     matcher.Match,
		captured:         map[string]interface{}{},
		path:             []string{},
	}
}

type comparator struct {
	loadExternalData LoadExternalResourceFunction
	valueMatcher     MatcherFunction
	captured         map[string]interface{}
	path             []string
}

func (c *comparator) GetCaptured() map[string]interface{} {
	return c.captured
}

func (c *comparator) Reset() {
	c.captured = map[string]interface{}{}
}

func (c *comparator) Compare(actual interface{}, expected interface{}) error {
	actualKind := reflect.ValueOf(actual).Kind()
	expectedKind := reflect.ValueOf(expected).Kind()

	if actualKind != expectedKind {
		if expectedKind != reflect.String {
			return ErrorAt(c.path, fmt.Errorf("%w : got %v want %v or %v", ErrTypeNotmatching, expectedKind, reflect.String, reflect.Slice))
		}
		if actualKind == reflect.Map {
			return ErrorAt(c.path, fmt.Errorf("%w : got %v want %v", ErrTypeNotmatching, actualKind, reflect.Slice))
		}
		if actualKind == reflect.Slice {
			return c.compareWithExternalRessource(actual.([]interface{}), expected.(string))
		}
	}
	switch actualKind {
	case reflect.Map:
		return c.compareMap(actual.(map[string]interface{}), expected.(map[string]interface{}))
	case reflect.Slice:
		return c.compareSlice(actual.([]interface{}), expected.([]interface{}))
	}
	return c.matchAndCapture(actual, expected)
}

func (c *comparator) matchAndCapture(actual interface{}, expected interface{}) error {
	capturedName, realExpected := splitCapturedNameAndExpectedValue(expected)
	err := c.valueMatcher(actual, realExpected)
	if err != nil {
		return ErrorAt(c.path, err)
	}
	if len(capturedName) > 0 {
		c.captured[capturedName] = actual
	}
	return nil
}

func splitCapturedNameAndExpectedValue(expected interface{}) (string, interface{}) {
	expectedAsString, ok := expected.(string)
	if !ok {
		return "", expected
	}
	expression := regexp.MustCompile(`\#(.*?)\=\{\{(.*?)\}\}`)
	result := expression.FindStringSubmatch(expectedAsString)
	if len(result) == 3 {
		return result[1], result[2]
	}
	return "", expected
}

func sliceToMap(in []interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for i, v := range in {
		out[fmt.Sprintf("%d", i)] = v
	}
	return out
}

func (c *comparator) compareSlice(actual, expected []interface{}) error {
	return c.compareMap(sliceToMap(actual), sliceToMap(expected))
}

func (c *comparator) compareMap(actual, expected map[string]interface{}) error {
	err := c.checkAllExpectedKeyArePresent(actual, expected)
	if err != nil {
		return err
	}
	return c.searchForExtraKey(actual, expected)
}

func (c *comparator) compareWithExternalRessource(actual []interface{}, path string) error {
	ext, err := c.loadExternalData(path)
	c.path = append(c.path, fmt.Sprintf("[%s]", path))
	if err != nil {
		return ErrorAt(c.path, fmt.Errorf("%w : %v", ErrRessourceNotFound, err))
	}
	for i, v := range actual {
		c.path = append(c.path, fmt.Sprintf("%d", i))
		err := c.Compare(v, ext)
		if err != nil {
			return err
		}
		c.path = c.path[:len(c.path)-1]
	}
	c.path = c.path[:len(c.path)-1]
	return nil
}

func isKeyOptional(key string) bool {
	if len(key) < 2 {
		return false
	}
	return key[0] == '?'
}
func getRealKeyNameOfOptional(key string) string {
	return key[1:]
}

func buildOptionalKeyNameFromReal(key string) string {
	return "?" + key
}

func (c *comparator) checkAllExpectedKeyArePresent(actual, expected map[string]interface{}) error {
	for k, ev := range expected {
		c.path = append(c.path, k)
		keyIsOptional := false
		if isKeyOptional(k) {
			k = getRealKeyNameOfOptional(k)
			keyIsOptional = true
		}
		av, ok := actual[k]
		if !ok {
			if keyIsOptional {
				c.path = c.path[:len(c.path)-1]
				return nil
			}
			return ErrorAt(c.path, fmt.Errorf("%w :%s", ErrMissingKey, k))
		}
		err := c.Compare(av, ev)
		if err != nil {
			return err
		}
		c.path = c.path[:len(c.path)-1]
	}
	return nil
}

func (c *comparator) searchForExtraKey(actual, expected map[string]interface{}) error {
	for k := range actual {
		c.path = append(c.path, k)
		_, ok := expected[k]
		_, optionalOk := expected[buildOptionalKeyNameFromReal(k)]
		if !ok && !optionalOk {
			return ErrorAt(c.path, fmt.Errorf("%w :%s", ErrExtraKey, k))
		}
		c.path = c.path[:len(c.path)-1]
	}
	return nil
}
