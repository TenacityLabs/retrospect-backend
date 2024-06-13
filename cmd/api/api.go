package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/cmd/service/capsule"
	"github.com/TenacityLabs/time-capsule-backend/cmd/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (server *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewUserStore(server.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	capsuleStore := capsule.NewCapsuleStore(server.db)
	capsuleHandler := capsule.NewHandler(capsuleStore, userStore)
	capsuleHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", server.addr)
	return http.ListenAndServe(server.addr, router)
}
