package testerconfig

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Url    string
	Groups map[string]TestGroup
}

type TestGroup struct {
	GroupName             string
	GlobalSetupCommand    string
	GlobalTearDownCommand string
	SetupCommand          string
	TeardownCommand       string
	Environment           map[string]string
	UnitTests             []UnitTest
	Scenarios             map[string][]UnitTest
}

type UnitTest struct {
	Action  string
	Url     string
	Status  int
	Headers map[string]string
	In      []byte
	Out     []byte
	CtIn    string
	CtOut   string
}

type ConfigLoader struct {
	fileOpener func(string) (io.ReadCloser, error)
}

func New() ConfigLoader {
	return ConfigLoader{
		fileOpener: func(path string) (io.ReadCloser, error) {
			return os.Open(path)
		},
	}
}

type ymlConfig struct {
	Url    string                  `yaml:"url"`
	Groups map[string]ymlTestGroup `yaml:"groups"`
}

type ymlTestGroup struct {
	GlobalSetupCommand    string   `yaml:"globalSetupCommand"`
	GlobalTearDownCommand string   `yaml:"globalTearDownCommand"`
	SetupCommand          string   `yaml:"setupCommand"`
	TeardownCommand       string   `yaml:"teardownCommand"`
	Environment           string   `yaml:"environment"`
	Tests                 []string `yaml:"tests"`
}

func (cl ConfigLoader) loadFile(filename string) ([]byte, error) {
	filename = filepath.Clean(filename)
	f, err := cl.fileOpener(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func (cl ConfigLoader) Load(filename string) (Config, error) {
	data, err := cl.loadFile(filename)
	if err != nil {
		return Config{}, err
	}
	yc := ymlConfig{}
	err = yaml.Unmarshal([]byte(data), &yc)
	if err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal file %s : %w", filename, err)
	}

	config := Config{
		Url:    yc.Url,
		Groups: map[string]TestGroup{},
	}
	for k, v := range yc.Groups {
		env, err := cl.loadEnvFile(k, v.Environment)
		if err != nil {
			return Config{}, fmt.Errorf("while loading env of group %s : %w", k, err)
		}
		units, scenarios, err := cl.loadTests(k, v.Tests)
		if err != nil {
			return Config{}, fmt.Errorf("while loading tests of group %s : %w", k, err)
		}
		config.Groups[k] = TestGroup{
			GroupName:             k,
			GlobalSetupCommand:    v.GlobalSetupCommand,
			GlobalTearDownCommand: v.GlobalTearDownCommand,
			SetupCommand:          v.SetupCommand,
			TeardownCommand:       v.TeardownCommand,
			Environment:           env,
			UnitTests:             units,
			Scenarios:             scenarios,
		}
	}
	return config, nil
}

func (cl ConfigLoader) loadEnvFile(group string, filename string) (map[string]string, error) {
	if len(filename) == 0 {
		return map[string]string{}, nil
	}
	data, err := cl.loadFile(group + "/" + filename)
	if err != nil {
		return nil, err
	}
	env := map[string]string{}
	err = json.Unmarshal([]byte(data), &env)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal file %s : %w", filename, err)
	}
	return env, nil
}

type ymlUnitTest struct {
	Action  string `yaml:"action"`
	Url     string `yaml:"url"`
	Status  int    `yaml:"status"`
	Headers string `yaml:"headers"`
	In      string `yaml:"in"`
	Out     string `yaml:"out"`
	CtIn    string `yaml:"ct_in"`
	CtOut   string `yaml:"ct_out"`
}

func (yut *ymlUnitTest) toUnitTest() (UnitTest, error) {

	h, err := parseHeader(yut.Headers)
	if err != nil {
		return UnitTest{}, err
	}
	out := UnitTest{
		Action:  yut.Action,
		Url:     yut.Url,
		Status:  yut.Status,
		Headers: h,
		CtIn:    yut.CtIn,
		CtOut:   yut.CtOut,
	}

	if len(out.CtIn) == 0 {
		out.CtIn = "application/json"
	}
	if len(out.CtOut) == 0 {
		out.CtOut = "application/json"
	}
	if out.Status == 0 {
		out.Status = 200
	}

	return out, nil
}

func parseHeader(headers string) (map[string]string, error) {
	out := map[string]string{}
	headerList := strings.Split(headers, ";")

	for _, h := range headerList {
		if len(h) == 0 {
			continue
		}
		part := strings.Split(h, ":")
		if len(part) != 2 {
			return nil, fmt.Errorf("header `%s must have only on `:`", h)
		}
		out[strings.TrimSpace(part[0])] = strings.TrimSpace(part[1])
	}
	return out, nil
}

type ymlTestConfig struct {
	UnitTests map[string][]ymlUnitTest `yaml:"unit_tests"`
	Scenarios map[string][]ymlUnitTest `yaml:"scenario"`
}

func (cl ConfigLoader) loadTests(group string, filenames []string) ([]UnitTest, map[string][]UnitTest, error) {
	if len(filenames) == 0 {
		return []UnitTest{}, map[string][]UnitTest{}, nil
	}
	uts := []UnitTest{}
	scenarios := map[string][]UnitTest{}

	for _, filename := range filenames {
		data, err := cl.loadFile(group + "/configs/" + filename)
		if err != nil {
			return nil, nil, fmt.Errorf("while loading %s : %w", filename, err)
		}
		config := ymlTestConfig{}
		err = yaml.Unmarshal([]byte(data), &config)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot unmarshal file %s : %w", filename, err)
		}

		actions := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
		for _, action := range actions {
			tests, ok := config.UnitTests[action]
			if !ok {
				continue
			}
			for _, v := range tests {
				v.Action = action
				u, err := v.toUnitTest()
				if err != nil {
					return nil, nil, err
				}
				if len(v.In) > 0 {
					ext := ".json"
					if u.CtIn != "application/json" {
						ext = ""
					}
					in, err := cl.loadFile(group + "/payloads/" + v.In + ext)
					if err != nil {
						return nil, nil, err
					}
					u.In = in
				}
				if len(v.Out) > 0 {
					ext := ".json"
					if u.CtOut != "application/json" {
						ext = ""
					}
					out, err := cl.loadFile(group + "/responses/" + v.Out + ext)
					if err != nil {
						return nil, nil, err
					}
					u.Out = out
				}
				uts = append(uts, u)
			}
		}
		for name, steps := range config.Scenarios {
			for _, v := range steps {
				u, err := v.toUnitTest()
				if err != nil {
					return nil, nil, err
				}
				if len(v.In) > 0 {
					ext := ".json"
					if u.CtIn != "application/json" {
						ext = ""
					}
					in, err := cl.loadFile(group + "/payloads/" + v.In + ext)
					if err != nil {
						return nil, nil, err
					}
					u.In = in
				}
				if len(v.Out) > 0 {
					ext := ".json"
					if u.CtOut != "application/json" {
						ext = ""
					}
					out, err := cl.loadFile(group + "/responses/" + v.Out + ext)
					if err != nil {
						return nil, nil, err
					}
					u.Out = out
				}
				scenarios[filename+":"+name] = append(scenarios[filename+":"+name], u)
			}
		}
	}

	return uts, scenarios, nil
}
