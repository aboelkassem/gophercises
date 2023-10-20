package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github/aboelkassem/gophercises-solutions/urlshort/urlshort"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/lib/pq" // load postgres driver, _ for not unused in the code
)

func main() {
	flagYamlFileName := flag.String("yaml", "urls.yaml", "Path to YAML file containing urls")
	flagJsonFileName := flag.String("json", "urls.json", "Path to JSON file containing urls")
	// start reading/parsing the above defined flags
	flag.Parse()

	// mux = web router to map
	mux := defaultMux()

	// map[key]value is hash map data structure
	// Build the MapHandler using the mux as the fallback
	// pathsToUrls := map[string]string{
	// 	"/urlshort":   "https://godoc.org/github.com/gophercises/urlshort",
	// 	"/yaml-godoc": "https://godoc.org/gopkg.in/yaml.v2",
	// }

	pathsToUrls, err := loadFromDB()

	if err != nil {
		fmt.Printf("Failed to load from db: %v", err)
		return
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlFile, err := os.Open(*flagYamlFileName)
	if err != nil {
		fmt.Printf("Failed to load file %v due to %v", flagYamlFileName, err)
		return
	}

	defer yamlFile.Close()

	yaml, err := ioutil.ReadAll(yamlFile)

	// yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `

	if err != nil {
		fmt.Printf("Failed to load file %v due to %v", flagYamlFileName, err)
		return
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Json file
	jsonFile, err := os.Open(*flagJsonFileName)
	if err != nil {
		fmt.Printf("Failed to load file %v due to %v", flagJsonFileName, err)
		return
	}

	defer jsonFile.Close()

	json, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Printf("Failed to load file %v due to %v", flagJsonFileName, err)
		return
	}

	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func loadFromDB() (map[string]string, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234567890 dbname=urls sslmode=disable")

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT path, url FROM urls`)

	if err != nil {
		return nil, err
	}

	// iterate for all rows
	var urls []urlshort.LinkData

	for rows.Next() {
		var url urlshort.LinkData
		// map
		if err := rows.Scan(&url.Path, &url.URL); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	pathsToUrls := map[string]string{}
	for _, url := range urls {
		pathsToUrls[url.Path] = url.URL
	}

	fmt.Println("urls = %v", pathsToUrls)

	return pathsToUrls, nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	// default route to handle all routes that not mapped/addressed
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
