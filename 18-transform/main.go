package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/", displayUploadForm)
	http.HandleFunc("/upload", handleUploadForm)
	// static files handler
	// proxy server for /temp as file server
	// by default add /temp to check within it, so we told him remove existing /temp
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("tmp"))))
	log.Print("Listening on port 3000....")
	http.ListenAndServe(":3000", nil)
}

func displayUploadForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "upload.html")
}

func handleUploadForm(w http.ResponseWriter, r *http.Request) {
	// read image from request
	f, _, err := r.FormFile("img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	img, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read other request data
	_ = img
	mode := r.FormValue("mode")
	num, err := strconv.Atoi(r.FormValue("num"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// use temp folder to store images
	cwd := filepath.Join("tmp", uuid.New().String())
	if err := os.MkdirAll(cwd, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store img
	if err := os.WriteFile(filepath.Join(cwd, "img.png"), img, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := make([]string, num)
	for n := 1; n <= num; n++ {
		path, err := run(filepath.Join(cwd, "img.png"), mode, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(path)
		// results[n-1] = path
		// replace \ to /
		results[n-1] = strings.Replace(path, string(filepath.Separator), "/", -1)
	}

	tmpl := template.Must(template.ParseFiles("results.html"))
	if err := tmpl.Execute(w, results); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func run(inPath string, mode string, n int) (string, error) {
	outDir := filepath.Dir(inPath)
	// get last element name in the path and remove extension  => image.png => image
	outName := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	outName += fmt.Sprintf("%05d", n)
	outName += filepath.Ext(inPath)
	outPath := filepath.Join(outDir, outName)
	// run cmd command
	cmd := exec.Command("primitive",
		"-i", inPath,
		"-o", outPath,
		"-n", strconv.Itoa(n),
		"-m", mode)
	log.Printf("Run: %s %v", cmd.Path, strings.Join(cmd.Args[1:], " "))
	out, err := cmd.Output()
	if err != nil {
		// cmd.Output has exitErr that has exitErr.Stderr
		exitErr := err.(*exec.ExitError)
		return "", fmt.Errorf("%v: %s", err, exitErr.Stderr)
	}
	log.Print(string(out))
	return outPath, nil
}
