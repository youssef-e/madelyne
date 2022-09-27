package unittester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/madelyne-io/madelyne/comparator"
	"github.com/madelyne-io/madelyne/tester/testerclient"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"github.com/madelyne-io/madelyne/tester/testerfile"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	ErrWrongStatus      = fmt.Errorf("Wrong status found")
	ErrWrongContentType = fmt.Errorf("Wrong ContentType found")
	ErrRawBodyDontMatch = fmt.Errorf("Wrong raw body found")
	ErrPcreNoResult     = fmt.Errorf("No result found")
)

type UnitTesterError struct {
	Ut     testerconfig.UnitTest
	Result []byte
	Err    error
}

func ErrorIn(ut testerconfig.UnitTest, r []byte, err error) *UnitTesterError {
	return &UnitTesterError{
		Ut:     ut,
		Result: r,
		Err:    err,
	}
}

func (e *UnitTesterError) Error() string {
	if e.Result == nil {
		return fmt.Sprintf("in test : '%s:%s'\nIn: %s\nOut: %s\nCtOut: %s\nStatus: %d\nHeaders: %s\nErr: %s", e.Ut.Action, e.Ut.Url, e.Ut.InName, e.Ut.OutName, e.Ut.CtOut, e.Ut.Status, e.Ut.Headers, e.Err.Error())
	}
	return fmt.Sprintf("in test : '%s:%s'\nIn: %s\nOut: %s\nCtOut: %s\nStatus: %d\nHeaders: %s\nErr: %s\ngot : \n%s", e.Ut.Action, e.Ut.Url, e.Ut.InName, e.Ut.OutName, e.Ut.CtOut, e.Ut.Status, e.Ut.Headers, e.Err.Error(), e.Result)
}
func (e *UnitTesterError) Unwrap() error { return e.Err }

type UnitTester struct {
	client      testerclient.Requester
	comparator  comparator.Comparator
	fileOpener  testerfile.FileOpener
	Environment map[string]string
}

func New(r testerclient.Requester, c comparator.Comparator, f testerfile.FileOpener) *UnitTester {
	return &UnitTester{
		client:      r,
		comparator:  c,
		fileOpener:  f,
		Environment: map[string]string{},
	}
}

func (t *UnitTester) Env() map[string]string {
	return t.Environment
}

func (t *UnitTester) RunSingle(ut testerconfig.UnitTest) error {
	if ut.In != nil {
		ut.In = ReplaceWithEnvValue(ut.In, t.Environment)
	}

	var err error
	switch ut.Action {
	case "FILE":
		err = t.runFile(ut)
	default:
		err = t.runApi(ut)
	}
	if err != nil {
		return err
	}

	for k, v := range t.comparator.GetCaptured() {
		t.Environment[k] = fmt.Sprintf("%v", v)
	}

	return nil
}

func (t *UnitTester) runApi(ut testerconfig.UnitTest) error {
	var sendedBody io.Reader
	if ut.In != nil {
		sendedBody = bytes.NewReader(ut.In)
	}
	request := testerclient.Request{
		Method:  ut.Action,
		Url:     ReplaceStringWithEnvValue(ut.Url, t.Environment),
		Body:    sendedBody,
		Headers: map[string]string{"Content-Type": ut.CtIn},
	}

	for key, value := range ut.Headers {
		request.Headers[key] = ReplaceStringWithEnvValue(value, t.Environment)
	}

	r, err := t.client.Make(request)
	if err != nil {
		return ErrorIn(ut, nil, fmt.Errorf("Error while requesting : %w", err))
	}

	if r.StatusCode != ut.Status {
		return ErrorIn(ut, nil, fmt.Errorf("%w: got %d expected %d.\nRsp: \n%s", ErrWrongStatus, r.StatusCode, ut.Status, getResponseBody(r)))
	}

	if ut.CtOut != "" && !strings.HasPrefix(r.ContentType, ut.CtOut) {
		return ErrorIn(ut, nil, fmt.Errorf("%w: %s expected %s.\nRsp: \n%s", ErrWrongContentType, r.ContentType, ut.CtOut, getResponseBody(r)))
	}

	if ut.Out != nil {
		ctOut := ut.CtOut
		if ctOut == "" {
			ctOut = r.ContentType
		}

		utErr := t.compareBody(r.Body, ut.Out, ctOut, ut.Pcre)
		if utErr != nil {
			utErr.Ut = ut
			return utErr
		}
	}

	return nil
}

func (t *UnitTester) runFile(ut testerconfig.UnitTest) error {
	ctOut := ut.CtOut
	if ut.CtOut == "" && strings.Contains(ut.InName, ".json") {
		ctOut = "application/json"
	}

	filename := filepath.Clean(ut.InName)
	f, err := t.fileOpener.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	utErr := t.compareBody(f, ut.Out, ctOut, ut.Pcre)
	if utErr != nil {
		utErr.Ut = ut
		return utErr
	}

	return nil
}

func (t *UnitTester) compareBody(left io.Reader, right []byte, expectedContentType, pattern string) *UnitTesterError {
	ut := testerconfig.UnitTest{}
	if left == nil {
		return ErrorIn(ut, nil, ErrRawBodyDontMatch)
	}
	leftBytes, err := ioutil.ReadAll(left)
	if err != nil {
		return ErrorIn(ut, nil, err)
	}

	if expectedContentType == "application/json" {
		t.comparator.Reset()
		var leftData interface{}
		err := json.Unmarshal(leftBytes, &leftData)
		if err != nil {
			return ErrorIn(ut, leftBytes, err)
		}
		var rightData interface{}
		err = json.Unmarshal(right, &rightData)
		if err != nil {
			return ErrorIn(ut, right, err)
		}
		err = t.comparator.Compare(leftData, rightData)
		if err != nil {
			return ErrorIn(ut, leftBytes, err)
		}
	} else {
		ok := bytes.Equal(leftBytes, right)
		if !ok {
			return ErrorIn(ut, leftBytes, ErrRawBodyDontMatch)
		}
	}

	if pattern != "" {
		err = t.comparator.Capture(leftBytes, pattern)
		if err != nil {
			return ErrorIn(ut, leftBytes, ErrPcreNoResult)
		}
	}

	return nil
}

func getResponseBody(r testerclient.Response) string {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)

		if r.ContentType == "application/json" {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, bodyBytes, "", "\t")
			if err == nil {
				bodyBytes = prettyJSON.Bytes()
			}
		}
	}

	return string(bodyBytes)
}
