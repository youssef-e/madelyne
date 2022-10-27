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
	lastScenarionFile string
}

func (ft *fakeTester) RunMultiple(ut []testerconfig.UnitTest) error {
	for _, u := range ut {
		ft.lastScenarionFile = u.File
	}
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
		order               []string
		in                  map[string]testerconfig.TestGroup
		fakeScenarioTesters []fakeTester
		expected            error
		expectedOrder       []string
	}{
		{
			order: []string{"fakename"},
			in: map[string]testerconfig.TestGroup{
				"fakename": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario"},
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
			order: []string{"fakename1", "fakename2"},
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario": []testerconfig.UnitTest{},
					},
				},
				"fakename2": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario1", "fakescenario2"},
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
			order: []string{"fakename1"},
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario"},
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
			order: []string{"fakename1"},
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario1", "fakename2"},
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
		{
			order: []string{"fakename1", "fakename2", "fakename3", "fakename4", "fakename5", "fakename6"},
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario1"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario1": []testerconfig.UnitTest{
							testerconfig.UnitTest{
								File: "file1",
							},
						},
					},
				},
				"fakename2": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario2"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario2": []testerconfig.UnitTest{
							testerconfig.UnitTest{
								File: "file2"},
						},
					},
				},
				"fakename3": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario3"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario3": []testerconfig.UnitTest{
							testerconfig.UnitTest{
								File: "file3"},
						},
					},
				},
				"fakename4": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario4"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario4": []testerconfig.UnitTest{
							testerconfig.UnitTest{
								File: "file4"},
						},
					},
				},
				"fakename5": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario5"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario5": []testerconfig.UnitTest{
							testerconfig.UnitTest{
								File: "file5"},
						},
					},
				},
				"fakename6": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario6"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario6": []testerconfig.UnitTest{
							testerconfig.UnitTest{
								File: "file6"},
						},
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
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
			},
			expected:      nil,
			expectedOrder: []string{"file1", "file2", "file3", "file4", "file5", "file6"},
		},
		{
			order: []string{"fakename1"},
			in: map[string]testerconfig.TestGroup{
				"fakename1": testerconfig.TestGroup{
					GroupName:     "fakename",
					UnitTests:     []testerconfig.UnitTest{},
					ScenarioOrder: []string{"fakescenario1", "fakescenario2", "fakescenario3", "fakescenario4", "fakescenario5", "fakescenario6"},
					Scenarios: map[string][]testerconfig.UnitTest{
						"fakescenario1": []testerconfig.UnitTest{
							testerconfig.UnitTest{File: "file1"},
						},
						"fakescenario2": []testerconfig.UnitTest{
							testerconfig.UnitTest{File: "file2"},
						},
						"fakescenario3": []testerconfig.UnitTest{
							testerconfig.UnitTest{File: "file3"},
						},
						"fakescenario4": []testerconfig.UnitTest{
							testerconfig.UnitTest{File: "file4"},
						},
						"fakescenario5": []testerconfig.UnitTest{
							testerconfig.UnitTest{File: "file5"},
						},
						"fakescenario6": []testerconfig.UnitTest{
							testerconfig.UnitTest{File: "file6"},
						},
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
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
				fakeTester{
					nextMultipleError: nil,
				},
			},
			expected:      nil,
			expectedOrder: []string{"file1", "file2", "file3", "file4", "file5", "file6"},
		},
	}

	for i, tt := range tests {
		tester := &SuiteTester{
			ScenarioTesterBuilder: NextFakeGroupScenarioTesterBuilder(tt.fakeScenarioTesters, t, i),
			CommandLauncher:       nopCommand,
			ProgressLogger:        func() {},
		}
		err := tester.RunSuite(tt.order, tt.in)

		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d failed got %v, exp %v", i, err, tt.expected)
		}

		for j, name := range tt.expectedOrder {
			if tt.fakeScenarioTesters[j].lastScenarionFile != name {
				t.Fatalf("%d failed on expectedOrder %d want %s got %s", i, j, name, tt.fakeScenarioTesters[j].lastScenarionFile)
			}
		}
	}
}

func TestRunGroupUnitTest(t *testing.T) {
	ErrFakeTest := fmt.Errorf("ErrFakeTest")

	tests := []struct {
		order           []string
		in              map[string]testerconfig.TestGroup
		fakeUnitTesters []fakeTester
		expected        error
	}{
		{
			order: []string{"fakename1"},
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
			order: []string{"fakename1"},
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
			order: []string{"fakename1"},
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
			order: []string{"fakename1", "fakename2"},
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
		err := tester.RunSuite(tt.order, tt.in)

		if !errors.Is(err, tt.expected) {
			t.Fatalf("%d failed got %v, exp %v", i, err, tt.expected)
		}

	}
}
