package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shubhrad1/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// feed, err := urlToFeed("https://feeds.bbci.co.uk/news/world/europe/rss.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("[ERROR] PORT not found in the env.")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("[ERROR] Database URL not found in the env.")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("[ERROR] Cannot connect to database:	", err)
	}
	db := database.New(conn)
	apicfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", healthCheck)
	v1Router.Get("/err", errorHandler)
	v1Router.Post("/users", apicfg.createUserHandler)
	v1Router.Get("/users", apicfg.middlewareAuth(apicfg.getUserHandler))
	v1Router.Post("/feeds", apicfg.middlewareAuth(apicfg.createFeedHandler))
	v1Router.Get("/feeds", apicfg.getFeedHandler)
	v1Router.Post("/feed_follows", apicfg.middlewareAuth(apicfg.createFeedFollowHandler))
	v1Router.Get("/feed_follows", apicfg.middlewareAuth(apicfg.getFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{feedFollowID}", apicfg.middlewareAuth(apicfg.deleteFeedFollowsHandler))
	v1Router.Get("/posts", apicfg.middlewareAuth(apicfg.getPostsHandler))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting at PORT=%v", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error:	", err)
	}
}
