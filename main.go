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

type arrayFlags []string

func (af *arrayFlags) String() string {
	return ""
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

type resource struct {
	method string
	path   string
}

var (
	blocked     []resource
	bind        *string
	backend     *string
	paths       arrayFlags
	dumpRequest *bool
)

func init() {
	bind = flag.String("bind", ":8080", "Front end address")
	backend = flag.String("backend", "localhost:3000", "Backend address")
	dumpRequest = flag.Bool("dump", false, "Dump HTTP request data")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nUsage:")
		fmt.Fprintf(os.Stderr, "    %v\n", os.Args[0])
	}
}

// isBlocked returns true if the path is in the list of blocked paths
func isBlocked(r *http.Request) bool {
	for _, b := range blocked {
		if r.URL.Path == b.path && (r.Method == b.method || b.method == "*") {
			log.Printf("blocked %v %v matching rule %+v",
				r.Method,
				r.URL.Path,
				b)
			return true
		}

	}

	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	if *dumpRequest == true {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Printf("DEUMP FAIL\n")
			log.Fatal(err)
		}

		fmt.Printf("%s\n", dump)
	}

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
	flag.Var(&paths, "block", "method:path to block. e.g. \"GET:/foo *:/bar\"")
	flag.Parse()

	// grab the blocked method:path pairs
	for _, path := range paths {
		parts := strings.Split(path, ":")
		block := resource{
			method: parts[0],
			path:   parts[1],
		}
		blocked = append(blocked, block)
	}

	log.Printf("blocked resources : %v", blocked)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
