package tester

import (
	"github.com/madelyne-io/madelyne/comparator"
	"github.com/madelyne-io/madelyne/tester/scenariotester"
	"github.com/madelyne-io/madelyne/tester/suitetester"
	"github.com/madelyne-io/madelyne/tester/testerclient"
	"github.com/madelyne-io/madelyne/tester/testercommand"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"github.com/madelyne-io/madelyne/tester/testerprogress"
	"github.com/madelyne-io/madelyne/tester/unittester"
	"os"
)

type Tester struct {
	suite  suitetester.SuiteTester
	groups map[string]testerconfig.TestGroup
}

func Load(confFile string) (*Tester, error) {
	config, err := testerconfig.New().Load(confFile)
	if err != nil {
		return nil, err
	}
	progress := testerprogress.New(os.Stdout, countSteps(config.Groups))
	tester := Build(config, testercommand.Run)
	tester.suite.ProgressLogger = func() { progress.Step() }
	return tester, nil
}

func Build(config testerconfig.Config, cmdLauncher func(cmd string) error) *Tester {
	return &Tester{
		groups: config.Groups,
		suite: suitetester.SuiteTester{
			CommandLauncher: cmdLauncher,
			UnitTesterBuilder: func(groupName string, env map[string]string) suitetester.UnitTester {
				ut := unittester.New(
					testerclient.New(config.Url),
					comparator.New(groupName),
				)
				for k, v := range env {
					ut.Env()[k] = v
				}
				return ut
			},
			ScenarioTesterBuilder: func(groupName string, env map[string]string) suitetester.ScenarioTester {
				st := scenariotester.New(func() scenariotester.UnitTester {
					return unittester.New(
						testerclient.New(config.Url),
						comparator.New(groupName),
					)
				})
				for k, v := range env {
					st.Env()[k] = v
				}
				return st
			},
		},
	}
}

func (t *Tester) Run() error {
	return t.suite.RunSuite(t.groups)
}

func countSteps(groups map[string]testerconfig.TestGroup) int {
	count := 0
	for _, g := range groups {
		count += len(g.UnitTests)
		count += len(g.Scenarios)
	}
	return count
}
