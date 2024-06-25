package main

import (
	"fmt"
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()

	dbHost, exist := os.LookupEnv("DB_HOST")

	if !exist {
		log.Fatal("DB_HOST not set in .env")
	}

	dbPort, exist := os.LookupEnv("DB_PORT")

	if !exist {
		log.Fatal("DB_PORT not set in .env")
	}

	dbUser, exist := os.LookupEnv("DB_USER")

	if !exist {
		log.Fatal("DB_USER not set in .env")
	}

	dbPass, exist := os.LookupEnv("DB_PASS")

	if !exist {
		log.Fatal("DB_PASS not set in .env")
	}

	dbName, exist := os.LookupEnv("DB_NAME")

	if !exist {
		log.Fatal("DB_NAME not set in .env")
	}

	if err := database.ConnectMySQL(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
		log.Fatal(err)
	}

	router := routes.ApiRoutes()

	fmt.Println("Muzz now serving love on port 6666 🚀")

	if err := http.ListenAndServe(":6666", router); err != nil {
		log.Fatal(err)
	}
}
