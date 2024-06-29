package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/TenacityLabs/retrospect-backend/config"
	"github.com/TenacityLabs/retrospect-backend/services/audio"
	"github.com/TenacityLabs/retrospect-backend/services/capsule"
	"github.com/TenacityLabs/retrospect-backend/services/doodle"
	"github.com/TenacityLabs/retrospect-backend/services/file"
	"github.com/TenacityLabs/retrospect-backend/services/miscFile"
	"github.com/TenacityLabs/retrospect-backend/services/photo"
	"github.com/TenacityLabs/retrospect-backend/services/questionAnswer"
	"github.com/TenacityLabs/retrospect-backend/services/song"
	"github.com/TenacityLabs/retrospect-backend/services/user"
	"github.com/TenacityLabs/retrospect-backend/services/writing"
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
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}
	defer client.Close()

	// Get a handle to your GCS bucket
	bucket := client.Bucket(config.Envs.GCSBucketName)

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewUserStore(server.db)
	capsuleStore := capsule.NewCapsuleStore(server.db)
	fileStore := file.NewFileStore(bucket)

	songStore := song.NewSongStore(server.db)
	questionAnswerStore := questionAnswer.NewQuestionAnswerStore(server.db)
	writingStore := writing.NewWritingStore(server.db)
	photoStore := photo.NewPhotoStore(server.db)
	audioStore := audio.NewAudioStore(server.db)
	doodleStore := doodle.NewDoodleStore(server.db)
	miscFileStore := miscFile.NewMiscFileStore(server.db)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	capsuleHandler := capsule.NewHandler(
		capsuleStore,
		userStore,
		fileStore,

		songStore,
		questionAnswerStore,
		writingStore,
		photoStore,
		audioStore,
		doodleStore,
		miscFileStore,
	)
	capsuleHandler.RegisterRoutes(subrouter)
	fileHandler := file.NewHandler(userStore, fileStore)
	fileHandler.RegisterRoutes(subrouter)

	songHandler := song.NewHandler(capsuleStore, userStore, songStore)
	songHandler.RegisterRoutes(subrouter)
	questionAnswerHanlder := questionAnswer.NewHandler(capsuleStore, userStore, questionAnswerStore)
	questionAnswerHanlder.RegisterRoutes(subrouter)
	writingHandler := writing.NewHandler(capsuleStore, userStore, writingStore)
	writingHandler.RegisterRoutes(subrouter)
	photoHandler := photo.NewHandler(capsuleStore, userStore, fileStore, photoStore)
	photoHandler.RegisterRoutes(subrouter)
	audioHandler := audio.NewHandler(capsuleStore, userStore, fileStore, audioStore)
	audioHandler.RegisterRoutes(subrouter)
	doodleHandler := doodle.NewHandler(capsuleStore, userStore, fileStore, doodleStore)
	doodleHandler.RegisterRoutes(subrouter)
	miscFileHandler := miscFile.NewHandler(capsuleStore, userStore, fileStore, miscFileStore)
	miscFileHandler.RegisterRoutes(subrouter)

	// TODO: limit origins for prod
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization"},
	})
	handler := c.Handler(router)

	log.Println("Listening on", server.addr)
	return http.ListenAndServe(server.addr, handler)
}
