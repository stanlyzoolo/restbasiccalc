package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)

type testData struct {
	Result int    `json:"result"`
	Error  error  `json:"error"`
	Expr   string `json:"expr"`
}

func TestHandler(t *testing.T) {

	testCases := map[string]struct {
		Result int
		Error  error
		Expr   string
		Status int
	}{
		"valid expr": {
			Result: 5,
			Error:  nil,
			Expr:   "1%2B4",
			Status: 200,
		},
		"invalid expr": {
			Result: 0,
			Error:  errors.New("unexpected token in tokenFactory() at position 2"),
			Expr:   "2%2B*",
			Status: 400,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/?expr=%s", tc.Expr), nil)

			handleExpr(w, r)

			rd := new(testData)

			err := rd.unmarshalJSON(w.Body.Bytes())

			if err != nil {
				t.Errorf("handler returns invalid json: %s", err)
			}

			if rd.Result != tc.Result {
				t.Errorf("want: %v, got: %v", tc.Result, rd.Result)
			}

			if w.Result().StatusCode != tc.Status {
				t.Errorf("expected status to equal 200, got:%v", w.Result().StatusCode)
			}

			if errors.Is(tc.Error, rd.Error) {
				t.Errorf("unexpected error:\n want: %v \n  got: %v", tc.Error, rd.Error)
			}	

		})
	}
}

func (rd *testData) unmarshalJSON(b []byte) error {
	type Alias testData
	aux := &struct {
		Error string `json:"error"`
		*Alias
	}{
		Alias: (*Alias)(rd),
	}

	err := json.Unmarshal(b, &aux)

	if err != nil {
		return err
	}

	rd.Error = errors.New(aux.Error)

	return nil
}
