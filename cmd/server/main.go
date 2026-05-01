package main

import (
	"log"
	"net/http"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/database"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/db/generated"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/repository"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/router"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.ConnectPostgres(cfg)
	defer db.Close()

	rdb := database.ConnectRedis(cfg)
	defer rdb.Close()

	queries := generated.New(db)

	userRepo := repository.NewUserRepository(db, queries, rdb)

	userSvc := service.NewUserService(userRepo)

	mux := router.Register(userSvc, cfg)

	log.Printf("Server running on :%s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, mux))
}
