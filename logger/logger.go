package logger

import (
	"log"
	"net/http"
)

const (
	WARN  = "[\033[33mWARN \033[0m] "
	INFO  = "[\033[34mINFO \033[0m] "
	ERROR = "[\033[31mERROR\033[0m] "
)

var (
	color_map = map[string]string{
		"GET":    "\033[32m",
		"POST":   "\033[33m",
		"PUT":    "\033[34m",
		"DELETE": "\033[31m",
	}
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r_type := r.Method
		r_path := r.URL.Path
		log.Printf("[%s%s%s] - %s", color_map[r_type], r_type, "\033[0m", r_path)
		next.ServeHTTP(w, r)
	})
}
