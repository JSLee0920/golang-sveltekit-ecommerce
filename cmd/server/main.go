package main

import (
	"fmt"
	"net/http"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/database"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! This is my first go backend project!!")
}

func main() {
	cfg := config.Load()
	db := database.ConnectPostgres(cfg)
	rdb := database.ConnectRedis(cfg)

	defer db.Close()
	defer rdb.Close()

	http.HandleFunc("/", helloHandler)

	fmt.Println("Server starting on port 8080.....")
	http.ListenAndServe(":8080", nil)
}
