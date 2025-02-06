package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type apiConfig struct{
	DB *database.Queries 

}

func main(){
	fmt.Println("hello world")
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if portString == ""{
		log.Fatal("DB is not found in the environment")
	}

	conn,err := sql.open("postgres",dbURL)
	if err != nil {
		log.Fatal("Cant connect to database ",err)
	}


	apiCfg := apiConfig{
		DB:database.New(conn)
	}
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, or specify specific domains like "https://example.com"
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,  // Allow credentials (cookies, authorization headers, etc.)
		MaxAge:           300,   // Cache preflight request for 5 minutes
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handleErr)
	v1router.post("/users",apiCfg.handlerCreateUser)
	v1router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	
	v1router.Get("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
		
	v1router.Get("/feed_follows",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	router.Mount("/v1",v1Router)
	
	

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	err:= srv.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Port:",portString)

}

