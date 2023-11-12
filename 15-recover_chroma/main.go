package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/debug/", debugDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// log.Println(err)
				stack := debug.Stack()
				// log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, addLinksToFiles(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func debugDemo(w http.ResponseWriter, r *http.Request) {
	filepath := strings.TrimPrefix(r.URL.Path, "/debug/")
	lineNumber, err := strconv.Atoi(r.URL.Query().Get("line"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read files

	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "%s:%d\n", filepath, lineNumber)
	// fmt.Fprintln(w)
	// err = quick.Highlight(w, string(f), "go", "html", "xcode-dark")
	err = syntaxHighlight(w, string(f), lineNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func addLinksToFiles(s string) string {
	// for windows only
	// for linux remove ([A-Za-z]:)?
	return regexp.
		MustCompile("([A-Za-z]:)?(\\/.+?):(\\d+)").
		ReplaceAllString(s, `
	<a href="/debug/$1$2?line=$3">$1$2$3</a>
	`)
}

// copy of file with little modifications
// C:\Users\Mohamed\go\pkg\mod\github.com\alecthomas\chroma@v0.10.0\quick
// 	err = quick.Highlight(w, string(f), "go", "html", "xcode-dark")

func syntaxHighlight(w io.Writer, source string, line int) error {
	var (
		lexer     = "go"
		formatter = "html"
		style     = "xcode-dark"
	)

	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	// parse formatter interface to struct

	f, ok := formatters.Get(formatter).(*html.Formatter)
	if !ok {
		return fmt.Errorf("undefined formatter: %s", formatter)
	}

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	// create set of options and apply it to the formatter
	for _, option := range []html.Option{
		html.TabWidth(4),
		html.WithLineNumbers(true),
		html.LinkableLineNumbers(true, "line-"),
		html.HighlightLines([][2]int{{line, line}}),
	} {
		option(f)
	}

	return f.Format(w, s, it)
}
