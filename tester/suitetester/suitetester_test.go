package suitetester

import (
	"errors"
	"fmt"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"testing"
)

func nopCommand(string) error {
	return nil
}

type fakeTester struct {
	nextMultipleError error
	nextSingleError   error
	t                 *testing.T
	i                 int
}

func (ft *fakeTester) RunMultiple(ut []testerconfig.UnitTest) error {
	return ft.nextMultipleError
}

func (ft *fakeTester) RunSingle(ut testerconfig.UnitTest) error {
	return ft.nextSingleError
}

func NextFakeGroupScenarioTesterBuilder(ft []fakeTester, t *testing.T, i int) func(string, map[string]string) ScenarioTester {
	scenarioIndex := 0
	return func(string, map[string]string) ScenarioTester {
		scenarioIndex += 1
		ft[scenarioIndex-1].t = t
		ft[scenarioIndex-1].i = i
		return &ft[scenarioIndex-1]
	}
}

func NextFakeGroupUnitTesterBuilder(ft []fakeTester, t *testing.T, i int) func(string, map[string]string) UnitTester {
	utIndex := 0
	return func(string, map[string]string) UnitTester {
		utIndex += 1
		ft[utIndex-1].t = t
		ft[utIndex-1].i = i
		return &ft[utIndex-1]
	}
}

func TestRunGroupScenario(t *testing.T) {
	ErrFakeTest := fmt.Errorf("ErrFakeTest")

	tests := []struct {
		in                  map[string]testerconfig.TestGroup
		fakeScenarioTesters []fakeTester
		expected            error
	}{
		{
			in: map[string]testerconfig.TestGroup{
				"fakename": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario": []testerconfig.UnitTest{},
					},
				},
			},
			fakeScenarioTesters: []fakeTester{
				fakeTester{
					nextMultipleError: nil,
				},
			},
			expected: nil,
		},
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario": []testerconfig.UnitTest{},
					},
				},
				"fakename2": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario1": []testerconfig.UnitTest{},
						"fakescenario2": []testerconfig.UnitTest{},
					},
				},
			},
			fakeScenarioTesters: []fakeTester{
				fakeTester{
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: ErrFakeTest,
				},
			},
			expected: ErrFakeTest,
		},
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario": []testerconfig.UnitTest{},
					},
				},
			},
			fakeScenarioTesters: []fakeTester{
				fakeTester{
					nextMultipleError: ErrFakeTest,
				},
			},
			expected: ErrFakeTest,
		},
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario1": []testerconfig.UnitTest{},
						"fakescenario2": []testerconfig.UnitTest{},
					},
				},
			},
			fakeScenarioTesters: []fakeTester{
				fakeTester{
					nextMultipleError: ErrFakeTest,
				},
				fakeTester{
					nextMultipleError: ErrFakeTest,
				},
			},
			expected: ErrFakeTest,
		},
	}

	for i, tt := range tests {
		tester := &SuiteTester{
			ScenarioTesterBuilder: NextFakeGroupScenarioTesterBuilder(tt.fakeScenarioTesters, t, i),
			CommandLauncher:       nopCommand,
			ProgressLogger:        func() {},
		}
		err := tester.RunSuite(tt.in)

		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d failed got %v, exp %v", i, err, tt.expected)
		}

	}
}

func TestRunGroupUnitTest(t *testing.T) {
	ErrFakeTest := fmt.Errorf("ErrFakeTest")

	tests := []struct {
		in              map[string]testerconfig.TestGroup
		fakeUnitTesters []fakeTester
		expected        error
	}{
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{
						testerconfig.UnitTest{},
					},
					Scenarios: map[string][]testerconfig.UnitTest{},
				},
			},
			fakeUnitTesters: []fakeTester{
				fakeTester{
					nextSingleError: nil,
				},
			},
			expected: nil,
		},
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{
						testerconfig.UnitTest{},
					},
					Scenarios: map[string][]testerconfig.UnitTest{},
				},
			},
			fakeUnitTesters: []fakeTester{
				fakeTester{
					nextSingleError: ErrFakeTest,
				},
			},
			expected: ErrFakeTest,
		},
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{
						testerconfig.UnitTest{},
						testerconfig.UnitTest{},
					},
					Scenarios: map[string][]testerconfig.UnitTest{},
				},
			},
			fakeUnitTesters: []fakeTester{
				fakeTester{
					nextSingleError: nil,
				},
				fakeTester{
					nextSingleError: ErrFakeTest,
				},
			},
			expected: ErrFakeTest,
		},
		{
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{
						testerconfig.UnitTest{},
						testerconfig.UnitTest{},
					},
					Scenarios: map[string][]testerconfig.UnitTest{},
				},
				"fakename2": testerconfig.TestGroup{
					GroupName: "fakename",
					UnitTests: []testerconfig.UnitTest{
						testerconfig.UnitTest{},
						testerconfig.UnitTest{},
					},
					Scenarios: map[string][]testerconfig.UnitTest{},
				},
			},
			fakeUnitTesters: []fakeTester{
				fakeTester{
					nextSingleError: nil,
				},
				fakeTester{
					nextSingleError: nil,
				},
				fakeTester{
					nextSingleError: nil,
				},
				fakeTester{
					nextSingleError: ErrFakeTest,
				},
			},
			expected: ErrFakeTest,
		},
	}

	for i, tt := range tests {
		tester := &SuiteTester{
			UnitTesterBuilder: NextFakeGroupUnitTesterBuilder(tt.fakeUnitTesters, t, i),
			CommandLauncher:   nopCommand,
			ProgressLogger:    func() {},
		}
		err := tester.RunSuite(tt.in)

		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d failed got %v, exp %v", i, err, tt.expected)
		}

	}
}
