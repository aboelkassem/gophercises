package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	ph := panicHandler{mux}
	log.Fatal(http.ListenAndServe(":3000", ph))
}

type panicHandler struct {
	fallback http.Handler
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// defer will be executed whatever happened to function
	defer func() {
		// catch if panic happened
		if err := recover(); err != nil {

			// 1- Logs the error, as well as the stack trace.
			log.Printf("panic  happened: %v\n", err)
			debug.PrintStack()

			// 2- Sets the status code to http.StatusInternalServerError (500) whenever a panic occurs.
			w.WriteHeader(http.StatusInternalServerError)

			// 3- Write a "Something went wrong" message when a panic occurs.
			fmt.Fprint(w, "Something went wrong")

			// 5- If the environment is set to be development, print the stack trace and
			// the error to the webpage as well as to the logs.
			// Otherwise default to the "Something went wrong" message described in (3).

			// run ENV=dev go run main.go
			if env, ok := os.LookupEnv("ENV"); ok && env == "dev" {
				fmt.Fprint(w, "\n\n")
				fmt.Fprintf(w, "panic: %v\n", err)
				fmt.Fprint(w, string(debug.Stack()))
			}
		}
	}()

	// 4- Ensure that partial writes and 200 headers aren't set
	// even if the handler started writing to the http.ResponseWriter
	// BEFORE the panic occurred (this one may be trickier) /panic-after/
	c := httptest.NewRecorder()
	h.fallback.ServeHTTP(c, r)

	// this will not be executed if panic happened
	fmt.Println("after fallback")

	// set the original headers if no panics occurred
	w.WriteHeader(c.Code)
	for k, vs := range c.HeaderMap {
		for _, v := range vs {
			w.Header().Set(k, v)
		}
	}
	c.Body.WriteTo(w)
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
