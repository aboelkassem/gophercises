package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href, Text string
}

func main() {
	flagHTMLFileName := flag.String("html", "ex1.html", "The Path to the HTML file to parse")
	flag.Parse()
	// Open the HTML file
	file, err := os.Open(*flagHTMLFileName)
	if err != nil {
		// print error and exit
		log.Fatalf("Failed to open HTML file: %s", err)
	}

	defer file.Close()

	// Parse the HTML document into tree node
	root, err := html.Parse(file)
	if err != nil {
		log.Fatalf("Failed to parse HTML file: %s", err)
	}

	// Traverse the root tree using BFS

	// Find all the anchors in the HTML document
	// using goroutines in parallel and listen if find any anchors print in
	aChan := make(chan *html.Node)
	go findAnchors(root, aChan)

	// listen for the channel
	// this loop will ended once the channel is closed

	for a := range aChan {
		log.Println(Link{
			Text: extractText(a),
			Href: extractHref(a),
		})
	}
}

func findAnchors(node *html.Node, aChan chan *html.Node) {
	// check if the node is an anchor tag
	if node.Type == html.ElementNode && node.Data == "a" {
		aChan <- node
		return
	}

	// traverse the children of the node using DFS
	// for i = 0 (init); i<= size (condition); i++ (increment)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		findAnchors(c, aChan)
	}

	// close the channel if the node is the root
	if node.Parent == nil {
		close(aChan)
	}
}

func extractHref(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func extractText(node *html.Node) string {
	var text string
	// loop through the children of the node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode { // TextNode = is a text node like <p>Text text</p>
			text += c.Data
		} else { // if not text like text<strong>ssss</strong>
			text += extractText(c) // call the function recursively to extract text from the child node
		}
	}
	return strings.TrimSpace(text)
}
