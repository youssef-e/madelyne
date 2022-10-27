package suitetester

import (
	"github.com/madelyne-io/madelyne/tester/testerconfig"
)

type ScenarioTester interface {
	RunMultiple(uts []testerconfig.UnitTest) error
}

type UnitTester interface {
	RunSingle(ut testerconfig.UnitTest) error
}

type SuiteTester struct {
	ProgressLogger        func()
	CommandLauncher       func(cmd string) error
	UnitTesterBuilder     func(groupName string, env map[string]string) UnitTester
	ScenarioTesterBuilder func(groupName string, env map[string]string) ScenarioTester
}

func (t *SuiteTester) runTest(setup string, teardown string, test func() error) error {
	err := t.CommandLauncher(setup)
	if err != nil {
		t.CommandLauncher(teardown)
		return err
	}
	err = test()
	if err != nil {
		t.CommandLauncher(teardown)
		return err
	}
	t.CommandLauncher(teardown)
	t.ProgressLogger()
	return nil
}

func (t *SuiteTester) RunSuite(order []string, groups map[string]testerconfig.TestGroup) error {
	for _, name := range order {
		group := groups[name]
		err := t.runTest(group.GlobalSetupCommand, group.GlobalTearDownCommand, func() error {
			err := t.runGroupUnitTest(group)
			if err != nil {
				return err
			}
			return t.runGroupScenario(group)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *SuiteTester) runGroupUnitTest(group testerconfig.TestGroup) error {
	for _, ut := range group.UnitTests {
		err := t.runTest(group.SetupCommand, group.TeardownCommand, func() error {
			return t.UnitTesterBuilder(group.GroupName, group.Environment).RunSingle(ut)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *SuiteTester) runGroupScenario(group testerconfig.TestGroup) error {
	for _, name := range group.ScenarioOrder {
		scenario := group.Scenarios[name]
		err := t.runTest(group.SetupCommand, group.TeardownCommand, func() error {
			return t.ScenarioTesterBuilder(group.GroupName, group.Environment).RunMultiple(scenario)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
