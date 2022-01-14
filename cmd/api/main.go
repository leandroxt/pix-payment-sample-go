package main

import (
	"flag"
	"log"
	"net/http"
)

type config struct {
	mpAccessToken string
}

type application struct {
	config config
}

func main() {
	var config config

	flag.StringVar(&config.mpAccessToken, "mp-access-token", "", "Mercado pago access token")
	flag.Parse()

	app := &application{
		config: config,
	}

	router := http.NewServeMux()
	router.HandleFunc("/process_payment", app.processPayment)

	fileServer := http.FileServer(http.Dir("./public"))
	router.Handle("/", http.StripPrefix("/", fileServer))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
