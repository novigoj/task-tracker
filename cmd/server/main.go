package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"task-tracker/internal/api/rest/handlers"
	"task-tracker/internal/config"
	"task-tracker/internal/connection/initialize"
	"task-tracker/internal/domain/middleware"
	"task-tracker/internal/domain/repository"
	"task-tracker/internal/domain/service"
	"task-tracker/internal/domain/types"
)

func sanitizeDBURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "<invalid DATABASE_URL>"
	}
	if u.User != nil {
		user := u.User.Username()
		u.User = url.UserPassword(user, "***")
	}
	return u.String()
}

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ERROR config: %v", err)
	}
	log.Printf("INFO  config loaded port=%s db=%s", cfg.Port, sanitizeDBURL(cfg.DatabaseURL))

	gormDB, err := initialize.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("ERROR database connect: %v", err)
	}
	log.Printf("INFO  database connected")

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("db sql error: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("db ping error: %v", err)
	}
	log.Printf("INFO  db ping ok")

	if err := gormDB.AutoMigrate(&types.User{}, &types.Task{}); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}
	log.Printf("INFO  db migrated")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.RecoveryJSON())
	router.Use(middleware.ErrorLogger())
	
	// CORS для фронтенда
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	
	// Статические файлы
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/")
	})	

	taskRepo := repository.NewTaskGormRepository(gormDB)
	taskService := service.NewTaskService(taskRepo)
	versionHandler := handlers.NewVersionHandler(taskService)	
	taskHandler := handlers.NewTaskHandler(taskService)

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		api.GET("/version", versionHandler.GetVersion)
		api.POST("/tasks", taskHandler.Create)
		api.GET("/tasks", taskHandler.List)
		api.GET("/tasks/:id", taskHandler.GetByID)
		api.PATCH("/tasks/:id", taskHandler.Update)
		api.DELETE("/tasks/:id", taskHandler.Delete)
	}

	addr := ":" + cfg.Port
	log.Printf("INFO  starting server addr=%s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("ERROR server: %v", err)
	}
}

