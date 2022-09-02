package main

import (
	"database/sql"
	"fmt"
	artDelivery "github.com/angelRaynov/clean-architecture/article/delivery/http"
	artMIddleware "github.com/angelRaynov/clean-architecture/article/delivery/http/middleware"
	artRepo "github.com/angelRaynov/clean-architecture/article/repository/db"
	artUsecase "github.com/angelRaynov/clean-architecture/article/usecase"
	authRepo "github.com/angelRaynov/clean-architecture/author/repository/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	fmt.Println(dsn)
	conn, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected")

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middleware := artMIddleware.InitMiddleware()
	e.Use(middleware.CORS)
	authorRepo := authRepo.NewAuthorRepository(conn)
	articleRepo := artRepo.NewArticleRepository(conn)

	to, err := strconv.Atoi(os.Getenv("CTX_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}
	timoutContext := time.Duration(to) * time.Second

	articleUsecase := artUsecase.NewArticleUseCase(articleRepo, authorRepo, timoutContext)
	artDelivery.NewArticleHandler(e, articleUsecase)

	log.Fatal(e.Start(os.Getenv("SERVER_ADDRESS")))
}
