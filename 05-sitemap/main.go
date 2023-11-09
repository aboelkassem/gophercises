package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aboelkassem/gophercises/sitemap/link"
)

func main() {
	// go run main.go -url https://blog.aboelkassem.tech -depth 1
	flagURL := flag.String("url", "", "The URL to create a sitemap for")
	flagDepth := flag.Int("depth", 2, "The depth of the links tree")
	flagXMLFileName := flag.String("xml", "sitemap.xml", "The sitemap file location to be saved")
	flag.Parse()

	if *flagURL == "" {
		log.Fatal("-url is required")
	}

	sitemapUrls, err := buildSitemap(*flagURL, *flagDepth)
	if err != nil {
		log.Fatalf("Failed to create sitemap: %s", err)
	}

	// // print urls
	// for _, url := range sitemapUrls {
	// 	log.Println(url)
	// }

	if err := generateSitemap(sitemapUrls, *flagXMLFileName); err != nil {
		log.Fatalf("Failed to generate sitemap in %s: %v", *flagXMLFileName, err)
	}

	log.Printf("Generate sitemap successfully with %d link(s) for %s in %s", len(sitemapUrls), *flagURL, *flagXMLFileName)
}

func buildSitemap(baseURL string, depth int) ([]string, error) {
	// just to ensure no duplication in O(1)
	urlsMap := map[string]bool{}

	urls := []string{baseURL}
	// instead of recursive, we will store var for links to be visited and after finish remove it
	for d := 0; d < depth; d++ {
		var newUlrs []string
		for _, url := range urls {
			newUrls, err := getUrls(url)
			if err != nil {
				return nil, err
			}

			var uniqueSubUrls []string
			for _, subUrl := range newUrls {
				// remove duplication
				if !urlsMap[subUrl] {
					uniqueSubUrls = append(uniqueSubUrls, subUrl)
					urlsMap[subUrl] = true
				}
			}

			newUlrs = append(newUlrs, uniqueSubUrls...) // uniqueSubUrls... = append uniqueSubUrls element by element
		}
		urls = newUlrs
	}

	return urls, nil
}

func getUrls(pageUrl string) ([]string, error) {
	// fetch the html page
	pageUrl = strings.TrimSuffix(pageUrl, "/")
	res, err := http.Get(pageUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	// res.Body is reader

	// parse page and get all <a>
	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var urls []string

	// get only href
	for _, link := range links {
		urls = append(urls, link.Href)
	}

	// filter out non-domain links

	var domainUrls []string

	for _, url := range urls {
		// skip http://google.com
		if strings.HasPrefix(url, "http") && !strings.HasSuffix(url, pageUrl) {
			continue
		}

		// http://domain.com/path/to/page
		if strings.HasSuffix(url, pageUrl) {
			domainUrls = append(domainUrls, url)
			continue
		}

		// mailto:email@example.com
		if strings.Contains(url, "@") {
			continue
		}

		// remove # portion of url
		// http://domain.com/path/to/page#anchor
		// take only = http://domain.com/path/to/page
		if i := strings.Index(url, "#"); i != -1 {
			url = url[:i]
		}

		// convert path to absolute path
		// firstCharacter := url[0]
		if url == "" || url[0] != '/' {
			url = "/" + url
		}
		url = pageUrl + url

		log.Println(url)
		domainUrls = append(domainUrls, url)
	}

	return domainUrls, nil
}

// Urlset was generated 2023-10-30 19:13:12 by https://xml-to-go.github.io/ in Ukraine.
// xml annotations, field tags
type SitemapXML struct {
	XMLName xml.Name        `xml:"urlset"`
	Xmlns   string          `xml:"xmlns,attr"`
	URLs    []SitemapXMLURL `xml:"url"`
}

type SitemapXMLURL struct {
	Loc string `xml:"loc"`
}

func generateSitemap(urls []string, pathToXML string) error {
	var sitemap SitemapXML
	sitemap.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

	for _, url := range urls {
		sitemap.URLs = append(sitemap.URLs, SitemapXMLURL{
			Loc: url,
		})
	}

	// f, error := os.Create(pathToXML)
	// if error != nil {
	// 	return error
	// }

	// defer f.Close()
	// return xml.NewEncoder(f).Encode(sitemap)

	// alternative way

	sitemapBytes, err := xml.MarshalIndent(&sitemap, "", "\t")
	if err != nil {
		return err
	}

	xmlData := []byte(xml.Header + string(sitemapBytes))

	// write to file
	return ioutil.WriteFile(pathToXML, xmlData, os.ModePerm)
}
