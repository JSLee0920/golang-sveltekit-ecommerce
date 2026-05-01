package router

import (
	"net/http"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/handler"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/middleware"
	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/service"
)

func Register(userSvc *service.UserService, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	userHandler := handler.NewUserHandler(userSvc, cfg)

	jwt := middleware.JWTProtected(cfg)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		middleware.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /api/v1/auth/register", userHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", userHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", userHandler.Logout)

	mux.Handle("GET /api/v1/users/me", jwt(http.HandlerFunc(userHandler.GetMe)))
	mux.Handle("PUT /api/v1/users/me", jwt(http.HandlerFunc(userHandler.UpdateMe)))

	return middleware.CORS(middleware.Logger(mux))
}
