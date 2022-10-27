package testerconfig

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func getTestFileOpener(fs map[string]string) func(string) (io.ReadCloser, error) {
	return func(path string) (io.ReadCloser, error) {
		file, ok := fs[path]
		if !ok {
			return nil, fmt.Errorf("file not found %s", path)
		}
		return ioutil.NopCloser(strings.NewReader(file)), nil
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		filesystem          map[string]string
		url                 string
		expected            map[string]TestGroup
		expectedGroupsOrder []string
	}{
		{
			filesystem: map[string]string{
				"conf.yml": `url: https://localhost:8000
groups:
  group1:
    globalSetupCommand: test1.sh
    globalTearDownCommand: test2.sh
    setupCommand: test3.sh
    teardownCommand: test4.sh
  group2:
    globalSetupCommand: test11.sh
    globalTearDownCommand: test12.sh
    setupCommand: test13.sh
    teardownCommand: test14.sh
    environment: env.json
  group3:
    globalSetupCommand: test21.sh
    globalTearDownCommand: test22.sh
    setupCommand: test23.sh
    teardownCommand: test24.sh
    environment: env.json
    tests:
      - access.yml
      - error.yml`,
				"group1/env.json": `{
	"key1":"value1",
	"key2":"value2"
}`,
				"group2/env.json": `{
	"key1":"value1",
	"key2":"value2"
}`,
				"group3/env.json": `{
	"key1":"value1",
	"key2":"value2"
}`,
				"group3/configs/access.yml": `unit_tests:
  GET:
    - { url: "/articles/all", out: "allArticles", headers: "Authorization : Bearer abc ; Test : value"}
  POST:
    - { url: "/articles", status: 201, in: "article", ct_in: "text/plain", out: "postedArticle" , ct_out: "text/plain"}
  PUT:
    - { url: "/articles/1", status: 200, in: "article", out: "updatedArticle" }
  PATCH:
    - { url: "/articles/1", status: 200, in: "article", out: "updatedArticle" }
  DELETE:
    - { url: "/articles/1", status: 204 }
scenario:
  createAndDeleteArticle:
    - { action: "POST", url: "/articles", status: 201, in: "postArticle" }
    - { action: "DELETE", url: "/articles/1", status: 204 }
  addArticle:
    - { action: "POST", url: "/articles", status: 201, in: "postArticle" }
`,
				"group3/configs/error.yml": `unit_tests:
  GET:
    - { url: "/articles/all", status: 404, out: "notfound" }
`,
				"group3/responses/allArticles.json":    "1",
				"group3/payloads/article.json":         "2",
				"group3/responses/postedArticle":       "3",
				"group3/responses/updatedArticle.json": "4",
				"group3/payloads/postArticle.json":     "5",
				"group3/payloads/article":              "6",
				"group3/responses/notfound.json":       "7",
			},
			url:                 "https://localhost:8000",
			expectedGroupsOrder: []string{"group1", "group2", "group3"},
			expected: map[string]TestGroup{
				"group1": TestGroup{
					GroupName:             "group1",
					GlobalSetupCommand:    "test1.sh",
					GlobalTearDownCommand: "test2.sh",
					SetupCommand:          "test3.sh",
					TeardownCommand:       "test4.sh",
					Environment:           map[string]string{},
					UnitTests:             []UnitTest{},
					ScenarioOrder:         []string{},
					Scenarios:             map[string][]UnitTest{},
				},
				"group2": TestGroup{
					GroupName:             "group2",
					GlobalSetupCommand:    "test11.sh",
					GlobalTearDownCommand: "test12.sh",
					SetupCommand:          "test13.sh",
					TeardownCommand:       "test14.sh",
					Environment: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
					UnitTests:     []UnitTest{},
					ScenarioOrder: []string{},
					Scenarios:     map[string][]UnitTest{},
				},
				"group3": TestGroup{
					GroupName:             "group3",
					GlobalSetupCommand:    "test21.sh",
					GlobalTearDownCommand: "test22.sh",
					SetupCommand:          "test23.sh",
					TeardownCommand:       "test24.sh",
					Environment: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
					UnitTests: []UnitTest{
						UnitTest{
							File:    "group3/configs/access.yml:GET",
							Action:  "GET",
							Url:     "/articles/all",
							Status:  200,
							In:      nil,
							CtIn:    "application/json",
							Out:     []byte("1"),
							CtOut:   "",
							OutName: "allArticles",
							Headers: map[string]string{
								"Authorization": "Bearer abc",
								"Test":          "value",
							},
						},
						UnitTest{
							File:    "group3/configs/access.yml:POST",
							Action:  "POST",
							Url:     "/articles",
							Status:  201,
							In:      []byte("6"),
							InName:  "article",
							CtIn:    "text/plain",
							Out:     []byte("3"),
							OutName: "postedArticle",
							CtOut:   "text/plain",
							Headers: map[string]string{},
						},
						UnitTest{
							File:    "group3/configs/access.yml:PUT",
							Action:  "PUT",
							Url:     "/articles/1",
							Status:  200,
							In:      []byte("2"),
							InName:  "article",
							CtIn:    "application/json",
							Out:     []byte("4"),
							OutName: "updatedArticle",
							CtOut:   "",
							Headers: map[string]string{},
						},
						UnitTest{
							File:    "group3/configs/access.yml:PATCH",
							Action:  "PATCH",
							Url:     "/articles/1",
							Status:  200,
							In:      []byte("2"),
							InName:  "article",
							CtIn:    "application/json",
							Out:     []byte("4"),
							OutName: "updatedArticle",
							CtOut:   "",
							Headers: map[string]string{},
						},
						UnitTest{
							File:    "group3/configs/access.yml:DELETE",
							Action:  "DELETE",
							Url:     "/articles/1",
							Status:  204,
							In:      nil,
							CtIn:    "application/json",
							Out:     nil,
							CtOut:   "",
							Headers: map[string]string{},
						},
						UnitTest{
							File:    "group3/configs/error.yml:GET",
							Action:  "GET",
							Url:     "/articles/all",
							Status:  404,
							In:      nil,
							CtIn:    "application/json",
							Out:     []byte("7"),
							OutName: "notfound",
							CtOut:   "",
							Headers: map[string]string{},
						},
					},
					ScenarioOrder: []string{
						"group3/configs/access.yml:addArticle",
						"group3/configs/access.yml:createAndDeleteArticle",
					},
					Scenarios: map[string][]UnitTest{
						"group3/configs/access.yml:createAndDeleteArticle": []UnitTest{
							UnitTest{
								File:    "group3/configs/access.yml:createAndDeleteArticle:POST:0",
								Action:  "POST",
								Url:     "/articles",
								Status:  201,
								In:      []byte("5"),
								InName:  "postArticle",
								CtIn:    "application/json",
								Out:     nil,
								CtOut:   "",
								Headers: map[string]string{},
							},
							UnitTest{
								File:    "group3/configs/access.yml:createAndDeleteArticle:DELETE:1",
								Action:  "DELETE",
								Url:     "/articles/1",
								Status:  204,
								In:      nil,
								CtIn:    "application/json",
								Out:     nil,
								CtOut:   "",
								Headers: map[string]string{},
							},
						},
						"group3/configs/access.yml:addArticle": []UnitTest{
							UnitTest{
								File:    "group3/configs/access.yml:addArticle:POST:0",
								Action:  "POST",
								Url:     "/articles",
								Status:  201,
								In:      []byte("5"),
								InName:  "postArticle",
								CtIn:    "application/json",
								Out:     nil,
								CtOut:   "",
								Headers: map[string]string{},
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		loader := New()
		if loader.fileOpener == nil {
			t.Fatalf("failed loader should have a fileOpener func defined")
		}
		loader.fileOpener = getTestFileOpener(tt.filesystem)

		result, err := loader.Load("conf.yml")
		if err != nil {
			t.Fatalf("%d failed %v", i, err)
		}
		if result.Url != tt.url {
			t.Fatalf("%d failed \n exp %#v \n got %#v", i, tt.url, result.Url)
		}
		if !reflect.DeepEqual(result.GroupsOrder, tt.expectedGroupsOrder) {
			t.Fatalf("%d group order failed \n exp %#v \n got %#v", i, tt.expectedGroupsOrder, result.GroupsOrder)
		}
		if !reflect.DeepEqual(result.Groups, tt.expected) {
			t.Fatalf("%d failed \n exp %#v \n got %#v", i, tt.expected, result.Groups)
		}
	}
}
