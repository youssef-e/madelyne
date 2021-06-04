package madtesting

import (
	"github.com/madelyne-io/madelyne/tester"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MadelyneConf struct {
	Handler               http.HandlerFunc
	ConfFileName          string
	Env                   map[string]string
	GlobalSetupCommand    func() error
	GlobalTearDownCommand func() error
	SetupCommand          func() error
	TeardownCommand       func() error
}

func RunTestSuite(t *testing.T, mc MadelyneConf) {
	app := httptest.NewServer(http.HandlerFunc(mc.Handler))
	defer app.Close()
	conf, err := testerconfig.New().Load(mc.ConfFileName)
	if err != nil {
		t.Fatalf("can't load madelyne conf  %s", err.Error())
	}
	conf.Url = app.URL

	for name, group := range conf.Groups {
		group.Environment = mc.Env
		group.GlobalSetupCommand = "GlobalSetupCommand"
		group.GlobalTearDownCommand = "GlobalTearDownCommand"
		group.SetupCommand = "SetupCommand"
		group.TeardownCommand = "TeardownCommand"
		conf.Groups[name] = group
	}

	mt := tester.Build(conf, func(name string) error {
		if name == "GlobalSetupCommand" && mc.GlobalSetupCommand != nil {
			return mc.GlobalSetupCommand()
		}
		if name == "GlobalTearDownCommand" && mc.GlobalTearDownCommand != nil {
			return mc.GlobalTearDownCommand()
		}
		if name == "SetupCommand" && mc.SetupCommand != nil {
			return mc.SetupCommand()
		}
		if name == "TeardownCommand" && mc.TeardownCommand != nil {
			return mc.TeardownCommand()
		}
		return nil
	})
	mt.Suite.ProgressLogger = func() {}
	err = mt.Run()
	if err != nil {
		t.Fatalf("testsuite failed %s", err.Error())
	}
}
