package main

import (
	"context"
	"net/http"

	"strconv"
	"time"

	"errors"

	"github.com/wojciech12/talk_monitoring_with_prometheus/middleware"
	"github.com/wojciech12/talk_monitoring_with_prometheus/service"
)

func main() {
	s := service.New("order-mgmt", "0.0.0.0:8080")
	m := middleware.NewMetricCollector(s.Name)
	s.MetricCollector = m
	s.SetupBasicRoutes()

	app := App{m}

	s.Router.Get("/hello", app.helloHandler)
	s.Router.Get("/world", app.worldHandler)
	s.Router.Get("/complex", app.complexHandler)

	s.Start()
}

type App struct {
	m service.MetricCollector
}

func (app *App) helloHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello!"))
}

func (app *App) worldHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("world!"))
}

func (app *App) complexHandler(writer http.ResponseWriter, request *http.Request) {
	vars := request.URL.Query()
	ctx := request.Context()
	sleepDb := float64(0)
	sleepAudit := float64(0)

	doDbError := false
	doAuditError := false

	v, ok := vars["db_sleep"]
	if ok {
		sleepDb, _ = strconv.ParseFloat(v[0], 64)
	}

	v, ok = vars["srv_sleep"]
	if ok {
		sleepAudit, _ = strconv.ParseFloat(v[0], 64)
	}

	v, ok = vars["is_db_error"]
	if ok {
		doDbError, _ = strconv.ParseBool(v[0])
	}

	v, ok = vars["is_srv_error"]
	if ok {
		doAuditError, _ = strconv.ParseBool(v[0])
	}

	err := app.callDatabase(ctx, sleepDb, doDbError)
	if err == nil {
		err = app.callAudit(ctx, sleepAudit, doAuditError)
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Success!"))
			return
		}
	}
	writer.Write([]byte("Failed " + err.Error() + "!"))
	writer.WriteHeader(http.StatusServiceUnavailable)
}

func (app *App) callDatabase(_ context.Context, sleepDuration float64, doError bool) error {
	// do sth with db
	if doError {
		app.m.ObserveDatabase("SQL_1", 1001, "HY000", 0.01)
		return errors.New("database")
	}
	if sleepDuration != 0 {
		simulateWork(sleepDuration)
		app.m.ObserveDatabase("SQL_1", 0, "0", sleepDuration)
	} else {
		v := float64(0.1)
		simulateWork(v)
		app.m.ObserveDatabase("SQL_1", 0, "0", v)
	}
	return nil
}

func (app *App) callAudit(_ context.Context, sleepDuration float64, doError bool) error {
	// do sth with db
	if doError {
		app.m.ObserveAudit("POST", "/audit", 500, 0)
		return errors.New("service")
	}
	if sleepDuration != 0 {
		simulateWork(sleepDuration)
		app.m.ObserveAudit("POST", "/audit", 200, sleepDuration)
	} else {
		v := float64(0.3)
		simulateWork(v)
		app.m.ObserveAudit("POST", "/audit", 200, v)
	}
	return nil
}

func simulateWork(seconds float64) {
	time.Sleep(time.Duration(seconds*1000) * time.Millisecond)
}
