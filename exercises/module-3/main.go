package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// read flags
	verbose := flag.Bool("verbose", false, "verbose output")
	flag.Parse()

	if *verbose {
		log.Printf("Version: %s", os.Getenv("VERSION"))
	}

	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/healthz", loggerMiddleware(headerMiddleware(healthHandler)))

	if *verbose {
		log.Printf("Listening on port 8080")
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("HTTP server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %s", err)
		}
	}()

	// Listen for SIGTERM and SIGINT signals to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	// Shutdown the server
	if *verbose {
		log.Println("Shutting down HTTP server...")
	}
	if err := server.Shutdown(nil); err != nil {
		if *verbose {
			log.Fatalf("HTTP server shutdown failed: %s", err)
		}
	}
	if *verbose {
		log.Println("HTTP server stopped.")
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// write a statement to log response status code
		writer := &responseStatusWriter{ResponseWriter: w}

		next.ServeHTTP(writer, r)
		log.Printf("%s %d", r.RemoteAddr, writer.statusCode)

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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
