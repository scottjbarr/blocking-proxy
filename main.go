package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var (
	blocked []string
	bind    *string
	backend *string
	paths   *string
)

func init() {
	bind = flag.String("bind", ":8080", "Front end address")
	backend = flag.String("backend", "localhost:3000", "Backend address")
	paths = flag.String("blocked",
		"",
		"method:path to block. e.g. \"GET:/foo *:/bar\"")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nUsage:")
		fmt.Fprintf(os.Stderr, "    %v\n", os.Args[0])
	}
}

// isBlocked returns true if the path is in the list of blocked paths
//
// Each string in the []string blocked will be in the format
// METHOD:/path. If the METHOD is "*" then all HTTP verbs are blocked
// for that path.
func isBlocked(r *http.Request) bool {
	for _, b := range blocked {
		parts := strings.Split(b, ":")
		method := parts[0]
		path := parts[1]

		if r.URL.Path == path && (r.Method == method || method == "*") {
			log.Printf("blocked %v %v matching rule %v %v",
				r.Method,
				r.URL.Path,
				method,
				path)
			return true
		}

	}

	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 404 any blocked resources
	if isBlocked(r) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// proxy the connection the backend
	director := func(req *http.Request) {
		req = r
		req.URL.Scheme = "http"
		req.URL.Host = *backend
	}

	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	// grab the blocked method:path pairs
	blocked = strings.Split(*paths, " ")

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
