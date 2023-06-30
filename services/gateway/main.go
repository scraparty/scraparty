package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(50, 1*time.Minute))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		Request(w, r)
	})

	http.ListenAndServe(":3000", r)
}

func Request(w http.ResponseWriter, r *http.Request) {
	services := strings.Split(r.URL.String()[1:], "/")
	service := services[0]

	if service == "api" {
		service = services[1]
	}

	dns := fmt.Sprintf("station-%s.default.src.cluster.local", service)

	ips, err := net.LookupIP(dns)

	if err != nil {
		http.Error(w, "There was an error forwarding your request. Please try again later.", http.StatusInternalServerError)

		return
	}

	ip := ips[0].String()

	address := fmt.Sprintf("http://%s:3000%s", ip, r.URL.String())

	url, err := url.Parse(address)

	if err != nil {
		http.Error(w, "There was an error forwarding your request. Please try again later.", http.StatusInternalServerError)

		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)	

	proxy.Director = func(r *http.Request) {
		r.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.URL.Host = url.Host
		r.URL.Path = url.Path
	}

	proxy.ServeHTTP(w, r)
}
