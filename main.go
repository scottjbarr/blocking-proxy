package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	blocked []string
)

// isBlocked returns true if the path is in the list of blocked paths
func isBlocked(path string) bool {
	for _, blockedPath := range blocked {
		if path == blockedPath {
			return true
		}

	}

	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 404 any blacklisted urls
	if isBlocked(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// proxy the connection the backend (hard coded local rails server)
	director := func(req *http.Request) {
		req = r
		req.URL.Scheme = "http"
		req.URL.Host = "localhost:3000"
	}

	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage : %v /blocked-url ...\n", os.Args[0])
		os.Exit(1)
	}

	blocked = os.Args[1:]
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
