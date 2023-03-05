package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("root."))
		if err != nil {
			return
		}
	})

	r.Get("/api/search", searchEmail)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}

func searchEmail(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("q")
	result := doZincSearch(queryParam)
	_, err := w.Write([]byte(result))
	if err != nil {
		return
	}
}

func doZincSearch(queryParam string) string {
	query := fmt.Sprintf(`{
	   "search_type": "match",
	   "query": {"term": "%s"}
	}`, queryParam)
	fmt.Println(query)
	req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
