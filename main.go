package main

import (
	"log"
	"net/http"
	"os"

	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func loggerMiddleware(next http.Handler) http.Handler {
	// Open a log file for writing.
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new logger instance that writes to the log file.
	logger := log.New(logFile, "", log.LstdFlags)

	// Return a new http.Handler that logs the request and calls the next middleware in the chain.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		if username != "myusername" || password != "mypassword" {
			w.Header().Set("WWW-Authenticate", `Basic realm="my realm"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}
		logger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// handleRequest is the main request handler function that handles the HTTP request and sends a response back to the client.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Handle the request...
	w.Write([]byte("Hello"))
}

func main() {
	// Create a middleware chain using alice that includes our loggerMiddleware, cors middleware and the main handler.
	chain := alice.New(
		loggerMiddleware,
		cors.Default().Handler,
	).ThenFunc(handleRequest)

	// Create a new HTTP server with our middleware chain and handleRequest function.
	server := http.Server{
		Addr:    ":3030",
		Handler: chain,
	}

	// Start the server and log any errors that occur.
	log.Fatal(server.ListenAndServe())
}
