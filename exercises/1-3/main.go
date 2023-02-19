// create a basic http server that listens on port 8080
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.Handle("/healthz", loggerMiddleware(headerMiddleware(http.HandlerFunc(healthHandler))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loggerMiddleware(next http.Handler) http.Handler {
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

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// write request header to response
		for k, v := range r.Header {
			w.Header().Add(k, v[0])
		}
		// read environment variable and write it to response
		w.Header().Set("VERSION", os.Getenv("VERSION"))
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
