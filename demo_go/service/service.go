package service

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/go-chi/chi/v5"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricCollector interface {
	ObserveMyself(path string, method string, statusCode int, duration float64)
	ObserveDatabase(queryId string, statusCode int, sqlState string, duration float64)
	ObserveAudit(path string, method string, statusCode int, duration float64)
}

type Service struct {
	Name            string
	Host            string
	Port            int
	MetricCollector MetricCollector
	Router          *chi.Mux
	Stop            chan os.Signal
}

func New(name string, host string) *Service {
	s := Service{}
	s.Name = name
	s.Host = host
	s.Router = chi.NewRouter()
	return &s
}

// add the standard handlers
// we have omitted, e.g., /health and /info
func (s *Service) SetupBasicRoutes() {
	ph := promhttp.Handler()
	s.Router.Use(s.requestTimerMiddleware)
	s.Router.Get("/metrics", ph.ServeHTTP)
}

func (s *Service) Start() {
	s.Stop = make(chan os.Signal, 1)
	signal.Notify(s.Stop, os.Interrupt)

	h := &http.Server{Addr: s.Host, Handler: s.Router}
	// run server in background
	go func() {
		err := h.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting Service due to %v", err)
		}
	}()
	//wait for SIGTERM
	<-s.Stop
	log.Print("Service shutting down...")
}

// HTTP middleware setting a value on the request context
func (s *Service) requestTimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(orgWriter http.ResponseWriter, r *http.Request) {
		timer := NewTimer()
		w := &ResponseWriter{ResponseWriter: orgWriter, status: 200}
		defer func() {
			s.finalizeRequest(w, r, timer)
		}()
		next.ServeHTTP(w, r)
	})
}

func (s Service) finalizeRequest(w *ResponseWriter, r *http.Request, timer *Timer) {
	if err := recover(); err != nil {
		stack := make([]byte, 1024*8)
		stack = stack[:runtime.Stack(stack, false)]
		log.Fatalf("PANIC: %s\n%s", err, stack)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		//log.Info("%s \"%s %s %s\" %d", r.RemoteAddr, r.Method, r.URL, r.Proto, w.status)
	}
	s.MetricCollector.ObserveMyself(r.URL.Path, r.Method, w.status, timer.durationSeconds())
	// access log
	log.Printf("%s \"%s %s %s\" %d - %s", r.RemoteAddr, r.URL, r.Method, r.Proto, w.status, r.UserAgent())
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Shutdown allows to stop the HTTP Server gracefully
func (s *Service) Shutdown() {
	s.Stop <- os.Signal(os.Interrupt)
}
