package tester

import (
	"github.com/madelyne-io/madelyne/comparator"
	"github.com/madelyne-io/madelyne/tester/scenariotester"
	"github.com/madelyne-io/madelyne/tester/suitetester"
	"github.com/madelyne-io/madelyne/tester/testerclient"
	"github.com/madelyne-io/madelyne/tester/testercommand"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"github.com/madelyne-io/madelyne/tester/testerfile"
	"github.com/madelyne-io/madelyne/tester/testerprogress"
	"github.com/madelyne-io/madelyne/tester/unittester"
	"os"
)

type Tester struct {
	Suite       suitetester.SuiteTester
	GroupsOrder []string
	Groups      map[string]testerconfig.TestGroup
}

func Load(confFile string) (*Tester, error) {
	config, err := testerconfig.New().Load(confFile)
	if err != nil {
		return nil, err
	}
	progress := testerprogress.New(os.Stdout, countSteps(config.Groups))
	tester := Build(config, testercommand.Run)
	tester.Suite.ProgressLogger = func() { progress.Step() }
	return tester, nil
}

func Build(config testerconfig.Config, cmdLauncher func(cmd string) error) *Tester {
	return &Tester{
		Groups:      config.Groups,
		GroupsOrder: config.GroupsOrder,
		Suite: suitetester.SuiteTester{
			CommandLauncher: cmdLauncher,
			UnitTesterBuilder: func(groupName string, env map[string]string) suitetester.UnitTester {
				ut := unittester.New(
					testerclient.New(config.Url),
					comparator.New(groupName),
					testerfile.New(),
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
						testerfile.New(),
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
	return t.Suite.RunSuite(t.GroupsOrder, t.Groups)
}

func countSteps(groups map[string]testerconfig.TestGroup) int {
	count := 0
	for _, g := range groups {
		count += len(g.UnitTests)
		count += len(g.Scenarios)
	}
	return count
}
