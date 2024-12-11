package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hakaton/internal/handlers"
	"hakaton/internal/middleware"
	"hakaton/internal/repository"
	"log"
	"net/http"
)

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgresql://hakaton_db_user:VhcZZoDRUymYsiCPH2rbwGwSPnDDqsaT@dpg-ctbguh3tq21c73e8g030-a.oregon-postgres.render.com/hakaton_db")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Создаем репозиторий
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	compRepo := repository.NewCompanyRepository(db)
	imageRepo := repository.NewImageRepository(db)

	// Создаем обработчики
	handler := handlers.NewHandler(userRepo, gameRepo, compRepo, imageRepo)

	// Настраиваем маршруты
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// Сервируем статические файлы из папки uploads
	r.Static("/uploads", "../../internal/uploads")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/registration", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	r.POST("/create_or_update_game", handler.CreateOrUpdateGame)
	r.POST("/upload-image", handler.UploadImageHandler)
	r.GET("/game_data", handler.GetGame)

	//r.POST("create_update_game", handler.CreateGameHandler)

	r.POST("/create_company", handler.CreateCompany)

	//router.POST("/images", handler.SaveImageHandler)
	//router.POST("/users", handler.CreateUserHandler)
	//router.GET("/users/:email", handler.GetUserByEmailHandler)
	//router.POST("/games", handler.CreateGameHandler)
	//router.PUT("/games/:game_id", handler.UpdateGameHandler)
	//router.GET("/games/:company_id/:name", handler.GetGameByNameHandler)

	// Запуск сервера
	log.Fatal(r.Run(":8080"))
}
