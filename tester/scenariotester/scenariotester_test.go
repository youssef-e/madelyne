package scenariotester

import (
	"errors"
	"fmt"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"testing"
)

type fakeUnitTester struct {
	nextError error
	nextEnv   map[string]string
}

func (fut *fakeUnitTester) RunSingle(ut testerconfig.UnitTest) error {
	return fut.nextError
}

func (fut *fakeUnitTester) Env() map[string]string {
	return fut.nextEnv
}

func NextFakeUnitTesterBuilder(fut []fakeUnitTester) func() UnitTester {
	i := 0
	return func() UnitTester {
		i += 1
		return &fut[i-1]
	}
}

func TestRunMultiple(t *testing.T) {
	ErrUnitTest := fmt.Errorf("ErrUnitTest")

	tests := []struct {
		fakeUnitTesters []fakeUnitTester
		uts             []testerconfig.UnitTest
		startingEnv     map[string]string
		endingEnd       map[string]string
		expected        error
	}{
		{
			fakeUnitTesters: []fakeUnitTester{},
			uts:             []testerconfig.UnitTest{},
			startingEnv:     map[string]string{},
			endingEnd:       map[string]string{},
			expected:        nil,
		},
		{
			fakeUnitTesters: []fakeUnitTester{
				fakeUnitTester{
					nextError: nil,
					nextEnv:   map[string]string{},
				},
			},
			uts: []testerconfig.UnitTest{
				testerconfig.UnitTest{},
			},
			startingEnv: map[string]string{},
			endingEnd:   map[string]string{},
			expected:    nil,
		},
		{
			fakeUnitTesters: []fakeUnitTester{
				fakeUnitTester{
					nextError: nil,
					nextEnv:   map[string]string{"test": "test"},
				},
			},
			uts: []testerconfig.UnitTest{
				testerconfig.UnitTest{},
			},
			startingEnv: map[string]string{},
			endingEnd:   map[string]string{"test": "test"},
			expected:    nil,
		},
		{
			fakeUnitTesters: []fakeUnitTester{
				fakeUnitTester{
					nextError: nil,
					nextEnv:   map[string]string{"test": "test"},
				},
			},
			uts: []testerconfig.UnitTest{
				testerconfig.UnitTest{},
			},
			startingEnv: map[string]string{"init": "init"},
			endingEnd:   map[string]string{"init": "init", "test": "test"},
			expected:    nil,
		},
		{
			fakeUnitTesters: []fakeUnitTester{
				fakeUnitTester{
					nextError: ErrUnitTest,
					nextEnv:   map[string]string{"test": "test"},
				},
			},
			uts: []testerconfig.UnitTest{
				testerconfig.UnitTest{},
			},
			startingEnv: map[string]string{"init": "init"},
			endingEnd:   map[string]string{"init": "init"},
			expected:    ErrUnitTest,
		},
		{
			fakeUnitTesters: []fakeUnitTester{
				fakeUnitTester{
					nextError: nil,
					nextEnv:   map[string]string{"test1": "test1"},
				},
				fakeUnitTester{
					nextError: ErrUnitTest,
					nextEnv:   map[string]string{"test2": "test2"},
				},
			},
			uts: []testerconfig.UnitTest{
				testerconfig.UnitTest{},
				testerconfig.UnitTest{},
			},
			startingEnv: map[string]string{"init": "init"},
			endingEnd:   map[string]string{"init": "init", "test1": "test1"},
			expected:    ErrUnitTest,
		},
	}

	for i, tt := range tests {
		scenariotester := New(NextFakeUnitTesterBuilder(tt.fakeUnitTesters))
		for k, v := range tt.startingEnv {
			scenariotester.Env()[k] = v
		}

		err := scenariotester.RunMultiple(tt.uts)

		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d failed got %v, exp %v", i, err, tt.expected)
		}
		if len(tt.endingEnd) != len(scenariotester.Env()) {
			t.Fatalf("%d failed got %d expected %d", i, len(scenariotester.Env()), len(tt.endingEnd))
		}
		for k, v := range tt.endingEnd {
			_, ok := scenariotester.Env()[k]
			if !ok {
				t.Fatalf("%d should found %s ", i, k)
			}
			if scenariotester.Env()[k] != v {
				t.Fatalf("%d failed got %v exp %v ", i, scenariotester.Env()[k], v)
			}
		}
	}
}
