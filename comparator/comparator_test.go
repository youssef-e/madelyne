package comparator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func getTestLoaderFunc(name string, data string, t *testing.T) LoadExternalResourceFunction {
	return func(path string) (map[string]interface{}, error) {
		if path != name {
			return nil, fmt.Errorf("failed")
		}
		result, ok := testReadJson(data, t).(map[string]interface{})
		if !ok {
			t.Fatalf("Data '%s' provide should be a map[string]interface{}", data)
		}
		return result, nil
	}
}

func testReadJson(data string, t *testing.T) interface{} {
	var result interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		t.Fatalf("invalid input json =>%s", err.Error())
	}
	return result
}

func TestNewComparator(t *testing.T) {
	c := New("test")
	cAsComparator, ok := c.(*comparator)
	if !ok {
		t.Fatalf("not expected type")
	}
	if cAsComparator.loadExternalData == nil {
		t.Fatalf("loadExternalJson func is nil")
	}
}

func TestComparator(t *testing.T) {
	tests := []struct {
		title           string
		left            string
		right           string
		externalLoarder LoadExternalResourceFunction
		capturedVars    map[string]interface{}
		expected        error
		expectedPath    []string
		pattern         string
		expectedPattern error
	}{
		{
			title:           "TestExtraKeyOnLeft",
			left:            `{"key1": "value1", "key2": "value2"}`,
			right:           `{"key1": "value1"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrExtraKey,
			expectedPath:    []string{"key2"},
		},
		{
			title:           "TestMissinfKeyOnLeft",
			left:            `{"key1": "value1"}`,
			right:           `{"key1": "value1", "key2": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrMissingKey,
			expectedPath:    []string{"key2"},
		},
		{
			title:           "TestComparingArrayWithMap",
			left:            `{"key1": "value1"}`,
			right:           `["key1"]`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrTypeNotmatching,
			expectedPath:    []string{},
		},
		{
			title:           "TestComparingMapWithArray",
			left:            `["key1"]`,
			right:           `{"key1": "value1"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrTypeNotmatching,
			expectedPath:    []string{},
		},
		{
			title:           "TestComparingArrayWithMapInSubKey",
			left:            `{"subkey":{"key1": "value1"}}`,
			right:           `{"subkey":["key1"]}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrTypeNotmatching,
			expectedPath:    []string{"subkey"},
		},
		{
			title:           "TestComparingMapWithArrayInSubKey",
			left:            `{"subkey":["key1"]}`,
			right:           `{"subkey":{"key1": "value1"}}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrTypeNotmatching,
			expectedPath:    []string{"subkey"},
		},
		{
			title:           "TestNotComparingAMapNorAnArray",
			left:            `1`,
			right:           `1`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestMissinfKeyOnLeft",
			left:            `{"key1": "value1", "key3": "value2"}`,
			right:           `{"key1": "value1", "key2": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrMissingKey,
			expectedPath:    []string{"key2"},
		},
		{
			title:           "TestSimpleValues",
			left:            `{"key": "value"}`,
			right:           `{"key": "value"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestSimpleValuesNotMatching",
			left:            `{"key": "value1"}`,
			right:           `{"key": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrNotMatching,
			expectedPath:    []string{"key"},
		},
		{
			title:           "TestSimpleValuesAllowPatternMatching",
			left:            `{"key": 1}`,
			right:           `{"key": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrNotMatching,
			expectedPath:    []string{"key"},
		},
		{
			title:           "TestSimpleArray",
			left:            `["key"]`,
			right:           `["key"]`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestLessSimpleArray",
			left:            `["key1","key2","key3","key4"]`,
			right:           `["key1","key2","key3","key4"]`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestShorterArray",
			left:            `["key1","key2","key3"]`,
			right:           `["key1","key2","key3","key4"]`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrMissingKey,
			expectedPath:    []string{"3"},
		},
		{
			title:           "TestLongerArray",
			left:            `["key1","key2","key3","key4"]`,
			right:           `["key1","key2","key3"]`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrExtraKey,
			expectedPath:    []string{"3"},
		},
		{
			title:           "TestTwoSimpleValues",
			left:            `{"key": "value", "key2": "value2"}`,
			right:           `{"key": "value", "key2": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestTwoNestedSimpleValues",
			left:            `{"key": "value", "key2": {"subkey":"subvalue"}}`,
			right:           `{"key": "value", "key2": {"subkey":"subvalue"}}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestTwoNestedArrayValues",
			left:            `{"key": "value", "key2": [{"subkey":"subvalue"},{"subkey":"subvalue"}]}`,
			right:           `{"key": "value", "key2": [{"subkey":"subvalue"},{"subkey":"subvalue"}]}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestNotAnArray",
			left:            `{"key1": [{"subKey": "not an array"}]}`,
			right:           `{"key1": "value1"}`,
			externalLoarder: getTestLoaderFunc("none", "{}", t),
			capturedVars:    map[string]interface{}{},
			expected:        ErrRessourceNotFound,
			expectedPath:    []string{"key1", "[value1]"},
		},
		{
			title:           "TestPartialFiles",
			left:            `{"key": [{"subkey": "subvalue"}]}`,
			right:           `{"key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestPartialFilesWithMultipleArrayElem",
			left:            `{"key": [{"subkey": "subvalue"},{"subkey": "subvalue"}]}`,
			right:           `{"key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestPartialFilesWithOneMissingKey",
			left:            `{"key": [{"subkey": "subvalue"},{"subother": "subvalue"}]}`,
			right:           `{"key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        ErrMissingKey,
			expectedPath:    []string{"key", "[array3]", "1", "subkey"},
		},
		{
			title:           "TestPartialFilesWithOneExtraKey",
			left:            `{"key": [{"subkey": "subvalue"},{"subkey": "subvalue","subother": "subvalue"}]}`,
			right:           `{"key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        ErrExtraKey,
			expectedPath:    []string{"key", "[array3]", "1", "subother"},
		},
		{
			title:           "TestPartialFilesNotMatching",
			left:            `{"key": [{"subkey": "subvalue"},{"subkey": "subvalue2"}]}`,
			right:           `{"key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        ErrNotMatching,
			expectedPath:    []string{"key", "[array3]", "1", "subkey"},
		},
		{
			title:           "TestCapturingVars",
			left:            `{"key1": "value", "key2": "value2"}`,
			right:           `{"key1": "#var1={{value}}", "key2": "#var2={{value2}}"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{"var1": "value", "var2": "value2"},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestTwoSimpleValuesWithOneOptional_allgood",
			left:            `{"key": "value", "key2": "value2"}`,
			right:           `{"key": "value", "?key2": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestTwoSimpleValuesWithOneOptional_optiondontraisean_error",
			left:            `{"key": "value"}`,
			right:           `{"key": "value", "?key2": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestTwoSimpleValuesWithOneOptional_notmatching",
			left:            `{"key": "value", "key2": "value1"}`,
			right:           `{"key": "value", "?key2": "value2"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        ErrNotMatching,
			expectedPath:    []string{"?key2"},
		},
		{
			title:           "TestOptionalPartialFiles_allgood",
			left:            `{"key": [{"subkey": "subvalue"}]}`,
			right:           `{"?key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestOptionalPartialFiles_optiondontraisean_error",
			left:            `{"other":"value"}`,
			right:           `{"other":"value","?key": "array3"}`,
			externalLoarder: getTestLoaderFunc("array3", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			expectedPath:    []string{},
		},
		{
			title:           "TestOptionalPartialFiles_notmatching",
			left:            `{"key": [{"subkey": "subothervalue"}]}`,
			right:           `{"?key": "array1"}`,
			externalLoarder: getTestLoaderFunc("array1", `{"subkey": "subvalue"}`, t),
			capturedVars:    map[string]interface{}{},
			expected:        ErrNotMatching,
			expectedPath:    []string{"?key", "[array1]", "0", "subkey"},
		},
		{
			title:           "TestCaptureValue1",
			left:            `{"key1": "value1"}`,
			right:           `{"key1": "value1"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{"pcre0": "value1", "pcre1": "1"},
			expected:        nil,
			pattern:         `value(\d)`,
			expectedPattern: nil,
		},
		{
			title:           "TestCaptureNotMatching",
			left:            `{"key1": "value"}`,
			right:           `{"key1": "value"}`,
			externalLoarder: nil,
			capturedVars:    map[string]interface{}{},
			expected:        nil,
			pattern:         `value([0-9])`,
			expectedPattern: ErrNotMatching,
		},
	}

	for i, tt := range tests {
		c := New("").(*comparator)
		c.loadExternalData = tt.externalLoarder
		c.valueMatcher = func(actual interface{}, expected interface{}) error {
			if actual == expected {
				return nil
			}
			return ErrNotMatching
		}

		c.Reset()

		err := c.Compare(testReadJson(tt.left, t), testReadJson(tt.right, t))
		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d:%s failed: got %v want %v", i, tt.title, err, tt.expected)
		}
		if err != nil {
			comparatorErr, ok := err.(*ComparatorError)
			if !ok {
				t.Fatalf("%d:%serror should be a *ComparatorError but got %T", i, tt.title, err)
			}
			if !reflect.DeepEqual(comparatorErr.Path, tt.expectedPath) {
				t.Fatalf("%d:%s path failed: got %#v want %#v", i, tt.title, comparatorErr.Path, tt.expectedPath)
			}
		}

		if tt.pattern != "" {
			err = c.Capture([]byte(tt.left), tt.pattern)
			if !errors.Is(err, tt.expectedPattern) {
				t.Fatalf("%d:%s failed: got %v want %v", i, tt.title, err, tt.expected)
			}
		}

		if tt.capturedVars == nil {
			continue
		}

		captured := c.GetCaptured()
		if len(captured) != len(tt.capturedVars) {
			t.Errorf("%#v", captured)
			t.Fatalf("%d:%s failed: have captured %v want %v", i, tt.title, len(captured), len(tt.capturedVars))
		}
		for key, val := range tt.capturedVars {
			capturedVal, ok := captured[key]
			if !ok {
				t.Fatalf("%d:%s A variable is missing from captured vars: %s", i, tt.title, key)
			} else if capturedVal != val {
				t.Fatalf("%d:%s Captured var is not what was expected. Expected: %v, Got: %v", i, tt.title, val, capturedVal)
			}
		}
	}
}
