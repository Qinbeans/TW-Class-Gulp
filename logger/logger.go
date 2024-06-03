package logger

import (
	"log"
	"net/http"
)

const (
	WARN  = "[\033[33mWARN \033[0m] "
	INFO  = "[\033[34mINFO \033[0m] "
	ERROR = "[\033[31mERROR\033[0m] "
	LOAD  = "[\033[32mFULL \033[0m] "
	BOOST = "[\033[3m\033[33mBOOST\033[0m] "
)

var (
	color_map = map[string]string{
		"GET":    "[\033[32mGET \033[0m] ",
		"POST":   "[\033[33mPOST\033[0m] ",
		"PUT":    "[\033[34mPUT \033[0m] ",
		"DELETE": "[\033[31mDEL \033[0m] ",
	}
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r_type := r.Method
		r_path := r.URL.Path
		boosted := LOAD
		if _, ok := r.Header["Hx-Boosted"]; ok {
			boosted = BOOST
		}
		log.Printf("%s%s- %s", color_map[r_type], boosted, r_path)
		next.ServeHTTP(w, r)
	})
}
