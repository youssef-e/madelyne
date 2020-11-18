package unittester

import (
	"errors"
	"fmt"
	"github.com/madelyne-io/madelyne/tester/testerclient"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

type fakeClient struct {
	nexResponse testerclient.Response
	nextError   error
	lastRequest testerclient.Request
}

func (fc *fakeClient) Make(r testerclient.Request) (testerclient.Response, error) {
	fc.lastRequest = r
	return fc.nexResponse, fc.nextError
}

type fakeComparator struct {
	nexCapturedEnv map[string]interface{}
	nextError      error
}

func (fc *fakeComparator) Compare(actual interface{}, expected interface{}) error {
	return fc.nextError
}
func (fc *fakeComparator) GetCaptured() map[string]interface{} {
	return fc.nexCapturedEnv
}
func (*fakeComparator) Reset() {}

func TestRunSingle(t *testing.T) {
	FakeMakeError := fmt.Errorf("FakeMakeError")
	FakeComparatorError := fmt.Errorf("FakeComparatorError")

	tests := []struct {
		input             testerconfig.UnitTest
		simulatedResponse testerclient.Response
		simulatedError    error
		comparatorResult  error
		comparatorCapture map[string]interface{}
		endEnv            map[string]string
		expected          error
	}{
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{},
				In:      nil,
				Out:     nil,
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        nil,
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          nil,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     nil,
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        nil,
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          nil,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     nil,
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        nil,
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    FakeMakeError,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          FakeMakeError,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     nil,
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  201,
				Body:        nil,
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          ErrWrongStatus,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     nil,
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        nil,
				ContentType: "test/plain",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          ErrWrongContentType,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     []byte("test"),
				CtIn:    "application/json",
				CtOut:   "test/plain",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        nil,
				ContentType: "test/plain",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          ErrRawBodyDontMatch,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     []byte("test1"),
				CtIn:    "application/json",
				CtOut:   "test/plain",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        ioutil.NopCloser(strings.NewReader("test2")),
				ContentType: "test/plain",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          ErrRawBodyDontMatch,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     []byte("test"),
				CtIn:    "application/json",
				CtOut:   "test/plain",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        ioutil.NopCloser(strings.NewReader("test")),
				ContentType: "test/plain",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          nil,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     nil,
				CtIn:    "application/json",
				CtOut:   "",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        nil,
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          nil,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     []byte("{\"something\":\"1\"}"),
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        ioutil.NopCloser(strings.NewReader("{\"something\":\"1\"}")),
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          nil,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     []byte("{\"something\":\"1\"}"),
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        ioutil.NopCloser(strings.NewReader("{\"something\":\"1\"}")),
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  FakeComparatorError,
			comparatorCapture: map[string]interface{}{},
			endEnv:            map[string]string{},
			expected:          FakeComparatorError,
		},
		{
			input: testerconfig.UnitTest{
				Action:  "GET",
				Url:     "/test",
				Status:  200,
				Headers: map[string]string{"test": "value"},
				In:      nil,
				Out:     []byte("{\"somethingOnlyTheCOmparatorWillMatch\":\"1\"}"),
				CtIn:    "application/json",
				CtOut:   "application/json",
			},
			simulatedResponse: testerclient.Response{
				StatusCode:  200,
				Body:        ioutil.NopCloser(strings.NewReader("{\"something\":\"1\"}")),
				ContentType: "application/json",
				Headers:     map[string][]string{},
			},
			simulatedError:    nil,
			comparatorResult:  nil,
			comparatorCapture: map[string]interface{}{"test": "value"},
			endEnv:            map[string]string{"test": "value"},
			expected:          nil,
		},
	}

	for i, tt := range tests {
		fakeClient := &fakeClient{
			nexResponse: tt.simulatedResponse,
			nextError:   tt.simulatedError,
		}
		fakeComparator := &fakeComparator{
			nexCapturedEnv: tt.comparatorCapture,
			nextError:      tt.comparatorResult,
		}

		unittester := New(fakeClient, fakeComparator)

		err := unittester.RunSingle(tt.input)
		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d failed got %v, exp %v", i, err, tt.expected)
		}

		if fakeClient.lastRequest.Method != tt.input.Action {
			t.Fatalf("%d failed got %v exp %v ", i, fakeClient.lastRequest.Method, tt.input.Action)
		}
		if fakeClient.lastRequest.Url != tt.input.Url {
			t.Fatalf("%d failed got %v exp %v ", i, fakeClient.lastRequest.Url, tt.input.Url)
		}
		if tt.input.In == nil {
			if fakeClient.lastRequest.Body != nil {
				t.Fatalf("%d lastRequest.Body should be nil", i)
			}
		} else {
			body, err := ioutil.ReadAll(fakeClient.lastRequest.Body)
			if err != nil {
				t.Fatalf("%d cannot read body of request %v", i, err)
			}
			if !reflect.DeepEqual(body, tt.input.In) {
				t.Fatalf("%d failed got %v exp %v ", i, fakeClient.lastRequest.Body, tt.input.In)
			}

		}
		if fakeClient.lastRequest.Headers["Content-Type"] != tt.input.CtIn {
			t.Fatalf("%d failed got %v exp %v ", i, fakeClient.lastRequest.Headers["Content-Type"], tt.input.CtIn)
		}

		for k, v := range tt.input.Headers {
			_, ok := fakeClient.lastRequest.Headers[k]
			if !ok {
				t.Fatalf("%d should found %s ", i, k)
			}
			if fakeClient.lastRequest.Headers[k] != v {
				t.Fatalf("%d failed got %v exp %v ", i, fakeClient.lastRequest.Headers[k], v)
			}
		}
		unitTestEnv := unittester.Env()

		if len(tt.endEnv) != len(unitTestEnv) {
			t.Fatalf("%d failed got %d expected %d", i, len(unitTestEnv), len(tt.endEnv))
		}

		for k, v := range tt.endEnv {
			_, ok := unitTestEnv[k]
			if !ok {
				t.Fatalf("%d should found %s ", i, k)
			}
			if unitTestEnv[k] != v {
				t.Fatalf("%d failed got %v exp %v ", i, unitTestEnv[k], v)
			}
		}
	}

}
