package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/services/capsule"
	"github.com/TenacityLabs/time-capsule-backend/services/file"
	"github.com/TenacityLabs/time-capsule-backend/services/questionAnswer"
	"github.com/TenacityLabs/time-capsule-backend/services/song"
	"github.com/TenacityLabs/time-capsule-backend/services/user"
	"github.com/TenacityLabs/time-capsule-backend/services/writing"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	capsuleStore := capsule.NewCapsuleStore(server.db)
	fileStore := file.NewFileStore()

	songStore := song.NewSongStore(server.db)
	questionAnswerStore := questionAnswer.NewQuestionAnswerStore(server.db)
	writingStore := writing.NewWritingStore(server.db)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	capsuleHandler := capsule.NewHandler(capsuleStore, userStore, songStore, questionAnswerStore, writingStore)
	capsuleHandler.RegisterRoutes(subrouter)
	fileHandler := file.NewHandler(userStore, fileStore)
	fileHandler.RegisterRoutes(subrouter)

	songHandler := song.NewHandler(capsuleStore, userStore, songStore)
	songHandler.RegisterRoutes(subrouter)
	questionAnswerHanlder := questionAnswer.NewHandler(capsuleStore, userStore, questionAnswerStore)
	questionAnswerHanlder.RegisterRoutes(subrouter)
	writingHandler := writing.NewHandler(capsuleStore, userStore, writingStore)
	writingHandler.RegisterRoutes(subrouter)

	// TODO: limit origins for prod
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization"},
	})
	handler := c.Handler(router)

	log.Println("Listening on", server.addr)
	return http.ListenAndServe(server.addr, handler)
}
