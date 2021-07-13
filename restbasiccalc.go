package main

import (
	"encoding/json"
	"net/http"

	. "github.com/stanlyzoolo/basiccalc"
)

type ReturnData struct {
	Result int    `json:"result"`
	Error  error  `json:"error"`
	Expr   string `json:"expr"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	var expr = r.URL.Query().Get("expr")
	// var expr = r.URL.Query()["expr"][1]

	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprint(w, "Hello: ", expr)

	result, err := Eval(expr)
	data, _ := json.Marshal(ReturnData{result, err, expr})
	w.Write(data)
	// fmt.Fprintf(w, "Result: %v. Error: %s", result, err)
}

func main() {

	http.HandleFunc("/", handler)

	http.HandleFunc("/test", handler)

	http.ListenAndServe(":8080", nil)
}

// Gracefully shutdown - читать
// os.signals - читать
