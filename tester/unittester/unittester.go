package unittester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/madelyne-io/madelyne/comparator"
	"github.com/madelyne-io/madelyne/tester/testerclient"
	"github.com/madelyne-io/madelyne/tester/testerconfig"
	"io"
	"io/ioutil"
	"strings"
)

var (
	ErrWrongStatus      = fmt.Errorf("Wrong status found")
	ErrWrongContentType = fmt.Errorf("Wrong ContentType found")
	ErrRawBodyDontMatch = fmt.Errorf("Wrong raw body found")
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
		return fmt.Sprintf("in test : '%s', CtOut : %s \nErr : %s", e.Ut.Url, e.Ut.CtOut, e.Err.Error())
	}
	return fmt.Sprintf("in test : '%s', CtOut : %s \nErr : %s \ngot : \n%s", e.Ut.Url, e.Ut.CtOut, e.Err.Error(), e.Result)
}
func (e *UnitTesterError) Unwrap() error { return e.Err }

type UnitTester struct {
	client      testerclient.Requester
	comparator  comparator.Comparator
	Environment map[string]string
}

func New(r testerclient.Requester, c comparator.Comparator) *UnitTester {
	return &UnitTester{
		client:      r,
		comparator:  c,
		Environment: map[string]string{},
	}
}

func (t *UnitTester) Env() map[string]string {
	return t.Environment
}

func (t *UnitTester) RunSingle(ut testerconfig.UnitTest) error {
	var sendedBody io.Reader
	if ut.In != nil {
		sendedBody = bytes.NewReader(ut.In)
	}
	request := testerclient.Request{
		Method:  ut.Action,
		Url:     ReplaceWithEnvValue(ut.Url, t.Environment),
		Body:    sendedBody,
		Headers: map[string]string{"Content-Type": ut.CtIn},
	}

	for key, value := range ut.Headers {
		request.Headers[key] = ReplaceWithEnvValue(value, t.Environment)
	}

	r, err := t.client.Make(request)
	if err != nil {
		return ErrorIn(ut, nil, fmt.Errorf("Error while requesting : %w", err))
	}

	if r.StatusCode != ut.Status {
		return ErrorIn(ut, nil, fmt.Errorf("%w: got %d expected %d", ErrWrongStatus, r.StatusCode, ut.Status))
	}

	ctOut := ut.CtOut
	if ctOut == "" {
		ctOut = "application/json"
	}

	if !strings.HasPrefix(r.ContentType, ut.CtOut) {
		return ErrorIn(ut, nil, fmt.Errorf("%w: %s expected %s.", ErrWrongContentType, r.ContentType, ut.CtOut))
	}

	if ut.Out == nil {
		return nil
	}

	utErr := t.compareBody(r.Body, ut.Out, ut.CtOut)
	if utErr != nil {
		utErr.Ut = ut
		return utErr
	}

	for k, v := range t.comparator.GetCaptured() {
		t.Environment[k] = fmt.Sprintf("%v", v)
	}

	return nil
}

func (t *UnitTester) compareBody(left io.Reader, right []byte, expectedContentType string) *UnitTesterError {
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
		return nil
	}

	ok := bytes.Equal(leftBytes, right)
	if !ok {
		return ErrorIn(ut, leftBytes, ErrRawBodyDontMatch)
	}
	return nil
}
