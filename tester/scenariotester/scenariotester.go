package scenariotester

import (
	"fmt"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
)

type UnitTester interface {
	RunSingle(ut testerconfig.UnitTest) error
	Env() map[string]string
}

type ScenarioTester struct {
	Environment       map[string]string
	UnitTesterBuilder func() UnitTester
}

func New(buildUnitTester func() UnitTester) *ScenarioTester {
	return &ScenarioTester{
		Environment:       map[string]string{},
		UnitTesterBuilder: buildUnitTester,
	}
}

func (t *ScenarioTester) Env() map[string]string {
	return t.Environment
}

func (t *ScenarioTester) RunMultiple(uts []testerconfig.UnitTest) error {

	for i, ut := range uts {
		unittester := t.UnitTesterBuilder()
		for k, v := range t.Environment {
			unittester.Env()[k] = v
		}
		err := unittester.RunSingle(ut)
		if err != nil {
			return fmt.Errorf("In test %d : %w", i, err)
		}
		for k, v := range unittester.Env() {
			t.Environment[k] = v
		}
	}
	return nil
}
