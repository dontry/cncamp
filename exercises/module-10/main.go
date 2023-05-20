package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dontry/cncamp/exercises/module-10/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// read flags
	verbose := flag.Bool("verbose", false, "verbose output")
	flag.Parse()

	if *verbose {
		log.Info().Msgf("Version: %s", os.Getenv("VERSION"))
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal().Msg("PORT environment variable is not set")
		panic("PORT environment variable is not set")
	}

	metrics.Register()

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", loggerMiddleware(metricsMiddleware(headerMiddleware(rootHandler))))
	mux.HandleFunc("/", loggerMiddleware((headerMiddleware(notFoundHandler))))
	mux.HandleFunc("/healthz", loggerMiddleware(metricsMiddleware(headerMiddleware(healthHandler))))

	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	if *verbose {
		log.Info().Msgf("Listening on port %s", os.Getenv("PORT"))
	}

	// Start the server in a goroutine
	go func() {
		log.Info().Msgf("HTTP server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("HTTP server failed: %s", err)
		}
	}()

	// Listen for SIGTERM and SIGINT signals to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	// Shutdown the server
	if *verbose {
		log.Info().Msg("Shutting down HTTP server...")
	}
	if err := server.Shutdown(nil); err != nil {
		if *verbose {
			log.Fatal().Msgf("HTTP server shutdown failed: %s", err)
		}
	}
	if *verbose {
		log.Info().Msg("HTTP server stopped.")
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// write a statement to log response status code
		writer := &responseStatusWriter{ResponseWriter: w}

		next.ServeHTTP(writer, r)
		log.Info().Msgf("%s %d", r.RemoteAddr, writer.statusCode)

	})
}

type responseStatusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseStatusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := metrics.NewTimer()
		defer timer.ObserveTotal()
		delay := randInt(10, 2000)
		time.Sleep(time.Millisecond * time.Duration(delay))
		next(w, r)
	})
}

func headerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// write request header to response
		for k, v := range r.Header {
			w.Header().Add(k, v[0])
		}
		// read environment variable and write it to response
		w.Header().Set("VERSION", os.Getenv("VERSION"))
		next(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	w.WriteHeader(http.StatusOK)
	if user != "" {
		w.Write([]byte("Hello World!" + user + "\n"))
	} else {
		w.Write([]byte("Hello [stranger]\n"))
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not Found\n"))
}

func initLogger() {
	// initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
