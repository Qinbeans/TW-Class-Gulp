package main

import (
	"backend/logger"
	"backend/template"
	"log"
	"net/http"
	"os"

	dotenv "github.com/joho/godotenv"
)

func getEnv() (address string, mode string) {
	dotenv.Load()

	mode = "dev"
	if os.Getenv("MODE") != "" {
		mode = os.Getenv("MODE")
	}
	log.Printf("%sRunning in %s mode", logger.INFO, mode)
	if mode == "dev" {

		port := "3000"
		address = ""

		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}

		if os.Getenv("ADDRESS") != "" {
			address = os.Getenv("ADDRESS")
		}

		address += ":" + port
		return
	}

	address = ":80"
	return
}

func main() {
	mux := http.NewServeMux()
	address, mode := getEnv()
	tmpl := template.NewEngine(mode)

	mux.HandleFunc("/", tmpl.Home)
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		if mode == "dev" {
			r.URL.Path = "tmp" + r.URL.Path
		} else {
			r.URL.Path = "dist" + r.URL.Path
		}
		http.ServeFile(w, r, r.URL.Path)
	})

	log.Printf("%sServer running at %s", logger.INFO, address)
	if err := http.ListenAndServe(address, logger.Logger(mux)); err != nil {
		log.Fatal(logger.ERROR, err)
	}
}
