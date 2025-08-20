package main

import (
	"pea-blog-backend/internal/config"
	"pea-blog-backend/internal/frontend"
	"pea-blog-backend/internal/handler"
	"pea-blog-backend/internal/middleware"
	"pea-blog-backend/internal/repository"
	"pea-blog-backend/internal/scheduler"
	"pea-blog-backend/internal/service"
	"pea-blog-backend/pkg/database"
	"pea-blog-backend/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	log := logger.New(cfg.Environment)

	// 前端自动构建
	buildService := frontend.NewBuildService(&cfg.Frontend, log)
	if err := buildService.CheckAndBuild(); err != nil {
		log.Error("Frontend build failed", err)
		// 不退出程序，继续启动服务
	}
	
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	err = database.Migrate(db)
	if err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	repos := repository.New(db)
	services := service.New(repos, log)
	handlers := handler.New(services, log)
	handlers.Image = handler.NewImageHandler(log)
	
	// 设置系统处理器
	handlers.System = handler.NewSystemHandler(buildService)

	// Start the scheduler
	sched := scheduler.New(services.Article, log)
	go sched.Start()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())

	// 静态文件服务 - 前端项目
	r.Static("/assets", cfg.Frontend.DistPath+"/assets")
	r.StaticFile("/", cfg.Frontend.DistPath+"/index.html")
	r.StaticFile("/favicon.ico", cfg.Frontend.DistPath+"/favicon.ico")
	
	// 对于前端路由，返回index.html让前端路由处理
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		// 否则返回前端index.html
		c.File(cfg.Frontend.DistPath + "/index.html")
	})

	api := r.Group("/api")
	
	auth := api.Group("/auth")
	{
		auth.POST("/login", handlers.Auth.Login)
		auth.POST("/logout", middleware.Auth(), handlers.Auth.Logout)
		auth.GET("/me", middleware.Auth(), handlers.Auth.GetCurrentUser)
		auth.POST("/refresh", handlers.Auth.RefreshToken)
	}

	articles := api.Group("/articles")
	{
		articles.GET("", middleware.Auth(), handlers.Article.GetArticles)
		articles.GET("/published", handlers.Article.GetPublishedArticles)
		articles.GET("/:id", handlers.Article.GetArticleByID)
		articles.GET("/search", handlers.Article.SearchArticles)
		articles.POST("", middleware.Auth(), middleware.AdminOnly(), handlers.Article.CreateArticle)
		articles.PUT("/:id", middleware.Auth(), middleware.AdminOnly(), handlers.Article.UpdateArticle)
		articles.DELETE("/:id", middleware.Auth(), middleware.AdminOnly(), handlers.Article.DeleteArticle)
		articles.POST("/:id/like", handlers.Article.LikeArticle)
		articles.DELETE("/:id/like", handlers.Article.UnlikeArticle)
		articles.POST("/:id/unpublish", middleware.Auth(), middleware.AdminOnly(), handlers.Article.UnpublishArticle)
		// articles.GET("/export", middleware.Auth(), middleware.AdminOnly(), handlers.Article.ExportArticles)
		// articles.POST("/import", middleware.Auth(), middleware.AdminOnly(), handlers.Article.ImportArticles)
		articles.GET("/:id/comments", handlers.Comment.GetCommentsByArticleID)
	}

	comments := api.Group("/comments")
	{
		comments.POST("", handlers.Comment.CreateComment)
		comments.DELETE("/:id", handlers.Comment.DeleteComment)
		comments.GET("/:id/replies", handlers.Comment.GetRepliesByCommentID)
	}

	system := api.Group("/system")
	{
		system.POST("/rebuild-frontend", middleware.Auth(), middleware.AdminOnly(), handlers.System.RebuildFrontend)
		system.GET("/build-status", handlers.System.GetBuildStatus)
	}

	log.Info("Server starting on port " + cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server", err)
	}
}