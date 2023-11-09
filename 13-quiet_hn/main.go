package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aboelkassem/gophercises/quiet_hn/hn"
)

var (
	cache     = map[int]hn.Item{}
	cacheLock sync.RWMutex
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		log.Printf("Got %d stories", len(ids))

		storyChan := make(chan orderedItem)

		// only get the first 25% of the stories
		ids = ids[:int(float64(numStories)*1.25)]

		//
		var wg sync.WaitGroup
		for i, id := range ids {
			// add a goroutine to the group to wait for
			wg.Add(1)

			// go routines for concurrency calls
			go func(itemId int, idx int) {
				// notify the wg group that this goroutine is done
				// will executed at the end of current context ended
				defer wg.Done()

				log.Printf("Fetching item %d at %d", itemId, idx)

				// check if exists in cache
				// race condition will be done here, since concurrency threads try to access the same memory at the same time
				// to solve it, use lock of sync package sync.RWMutex
				// check cacheRead and cacheWrite
				if _, ok := cacheRead(itemId); !ok {
					hnItem, err := client.GetItem(itemId)
					if err != nil {
						return
					}
					cacheWrite(itemId, hnItem)
				}

				// hnItem = cache[itemId]
				hnItem, _ := cacheRead(itemId)
				item := parseHNItem(hnItem)
				if isStoryLink(item) {
					storyChan <- orderedItem{item, idx}
				}
			}(id, i)
		}

		// start a monitoring goroutine to wait for all the group to be done
		go func() {
			wg.Wait()
			close(storyChan)
		}()

		var stories []orderedItem

		for orderedItem := range storyChan {
			stories = append(stories, orderedItem)
		}

		// sort by idx
		sort.Slice(stories, func(i, j int) bool {
			return stories[i].idx < stories[j].idx
		})

		data := templateData{
			Stories: stories[:numStories],
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type orderedItem struct {
	item
	idx int
}

type templateData struct {
	Stories []orderedItem
	Time    time.Duration
}

func cacheRead(id int) (hn.Item, bool) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	hnItem, ok := cache[id]
	return hnItem, ok
}

func cacheWrite(id int, hnItem hn.Item) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[id] = hnItem
}
