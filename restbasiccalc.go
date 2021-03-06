package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/namsral/flag"
	"github.com/stanlyzoolo/basiccalc"
	"go.uber.org/zap"
)

type returnData struct {
	Result int    `json:"result"`
	Error  error  `json:"error"`
	Expr   string `json:"expr"`
}

func handleExpr(w http.ResponseWriter, r *http.Request) {
	var expr = r.URL.Query().Get("expr")

	w.Header().Set("Content-Type", "application/json")

	logger, _ := zap.NewDevelopment()

	result, err := basiccalc.Eval(expr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("failed evaluating an expression",
			zap.String("400", "Bad Request"),
			zap.String("package", "restbasiccalc"),
			zap.Error(err))
	}

	rd := returnData{
		Result: result,
		Error:  err,
		Expr:   expr,
	}

	data, errW := rd.marshalJSON()
	if errW != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("failed evaluating an expression",
			zap.String("500", "Internal Server Error"),
			zap.String("package", "restbasiccalc"),
			zap.Error(errW))
	}
	w.Write(data) //nolint
}

func (rd *returnData) marshalJSON() ([]byte, error) {
	type Alias returnData

	return json.Marshal(&struct { //nolint
		Error string `json:"error"`
		*Alias
	}{
		Error: fmt.Sprint(rd.Error),
		Alias: (*Alias)(rd),
	})
}

func main() {
	var port = 8080

	flag.IntVar(&port, "port", port, "Port number")
	flag.Parse()
	// nolint
	var srv = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(handleExpr),
	}

	go func() {
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			logger, _ := zap.NewDevelopment()
			logger.Fatal("failed server start",
				zap.String("package", "restbasiccalc"),
				zap.String("func", "ListenAndServe"),
				zap.Error(err),
			)
		}
	}()

	logger, _ := zap.NewDevelopment()
	logger.Info("server started")

	gracefulShDw := make(chan os.Signal, 1)
	signal.Notify(gracefulShDw, os.Interrupt)

	<-gracefulShDw
	logger.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger, _ := zap.NewDevelopment()
		logger.Fatal("failed shutdown server",
			zap.String("package", "restbasiccalc"),
			zap.String("func", "Shutdown"),
			zap.Error(err),
		)
	}
}
