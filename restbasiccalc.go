package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	. "github.com/stanlyzoolo/basiccalc"
)

type ReturnData struct {
	Result int    `json:"result"`
	Error  error  `json:"error"`
	Expr   string `json:"expr"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	var expr = r.URL.Query().Get("expr")

	w.Header().Set("Content-Type", "application/json")

	result, err := Eval(expr)
	data, _ := json.Marshal(ReturnData{result, err, expr})
	w.Write(data)
}

func main() {

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
	}

	go func() {

		gracefulShDw := make(chan os.Signal, 1)
		signal.Notify(gracefulShDw, os.Interrupt)
		<-gracefulShDw

		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Printf("HTTP server Shutdown: %v", err)
		}

		close(gracefulShDw)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("HTTP server ListenAndServe: %v", err)

	}

}

// Gracefully shutdown - читать
// os.signals - читать
