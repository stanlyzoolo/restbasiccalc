package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// для теста можно подключить пакет exprgen

	testCases := []struct {
		Result int
		Error  error
		Expr   string
	}{
		{
			Result: 5,
			Error:  nil,
			Expr:   "1%2B4",
		},
		{
			Result: 9,
			Error:  nil,
			Expr:   "5%2B4",
		},
		{
			Result: 8,
			Error:  nil,
			Expr:   "2%2B6",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Expr, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/?expr=%s", tc.Expr), nil)
			handleExpr(w, r)
			rd := returnData{}
			err := json.Unmarshal(w.Body.Bytes(), &rd)
			if err != nil {
				t.Errorf("failed evaluating an expression: %w", err)
			}

			if rd.Result != tc.Result {
				t.Errorf("want: %v, got: %v", tc.Result, rd.Result)
			}

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("expected status to equal 200, got:%v", w.Result().StatusCode)
			}

		})
	}
}
