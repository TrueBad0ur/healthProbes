package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port              int
	t                 time.Time
	waitStartupTime   time.Duration
	waitLivenessTime  time.Duration
	waitReadinessTime time.Duration
}

func main() {
	var s Server
	s.port = 8080

	err := s.getEnvValues()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/startupProbe", s.startupProbe)
	http.HandleFunc("/livenessProbe", s.livenessProbe)
	http.HandleFunc("/readinessProbe", s.readinessProbe)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start time
	s.t = time.Now()

	fmt.Printf("Starting server. Listening on port: %d", s.port)
	log.Fatal(srv.ListenAndServe())
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvToDuration(e string, defaultValue string) (d time.Duration, err error) {
	var envValue int
	envValue, err = strconv.Atoi(getEnv(e, defaultValue))
	d = time.Duration(envValue) * time.Second
	return
}

func (s *Server) getEnvValues() (err error) {
	s.waitStartupTime, err = getEnvToDuration("WAIT_STARTUP_TIME", "15")
	if err != nil {
		return
	}
	s.waitLivenessTime, err = getEnvToDuration("WAIT_LIVENESS_TIME", "20")
	if err != nil {
		return
	}
	s.waitReadinessTime, err = getEnvToDuration("WAIT_READINESS_TIME", "20")
	if err != nil {
		return
	}
	return
}

// Was the start of app OK?
// livenessProbe and readinessProbe will start just after startupProbe ends true
func (s *Server) startupProbe(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.t) > s.waitStartupTime {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(503)
	}
}

// Does the app actually works/live?
func (s *Server) livenessProbe(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.t) > s.waitLivenessTime {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(503)
	}
}

// Is the app ready to serve the traffic?
func (s *Server) readinessProbe(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.t) > s.waitReadinessTime {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(503)
	}
}
