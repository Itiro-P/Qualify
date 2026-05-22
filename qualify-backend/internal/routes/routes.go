package routes

import (
	"main/internal/database/handlers"
	"main/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, conn *pgxpool.Pool) {
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},

		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},

		AllowCredentials: true,

		MaxAge: 12 * 3600,
	}))

	// Fotos de perfil
	router.Static("/uploads", "./uploads")

	// Rate limiters
	authRateLimiter := middleware.NewRateLimiter(5, 15*time.Minute)
	loginRateLimiter := middleware.NewRateLimiter(5, 5*time.Minute)

	// --- ROTAS PÚBLICAS ---
	router.POST("/register", middleware.RateLimitMiddleware(authRateLimiter), handlers.CreateUser(conn))

	auth := router.Group("/auth")
	{
		auth.POST("/login", middleware.RateLimitMiddleware(loginRateLimiter), handlers.Login(conn))
		auth.POST("/refresh", handlers.RefreshToken(conn))
		auth.POST("/reset-password", middleware.RateLimitMiddleware(authRateLimiter), handlers.ResetPassword(conn))
		auth.POST("/reset-password/confirm", handlers.ConfirmPasswordReset(conn))
	}

	// Leituras Públicas (Users, Analysts, Clients)
	router.GET("/analysts", handlers.GetAnalysts(conn))
	router.GET("/clients", handlers.GetClients(conn))

	usersPublic := router.Group("/users")
	{
		usersPublic.GET("/:id", handlers.GetUser(conn))
		usersPublic.GET("/:id/profile", handlers.GetUserProfile(conn))
		usersPublic.GET("/:id/analyst", handlers.GetAnalyst(conn))
		usersPublic.GET("/:id/analyst/skills", handlers.GetAnalystSkills(conn))
		usersPublic.GET("/:id/analyst/certifications", handlers.GetAnalystCertifications(conn))
		usersPublic.GET("/:id/analyst/profile", handlers.GetAnalystProfile(conn))
		usersPublic.GET("/:id/client", handlers.GetClient(conn))
		usersPublic.GET("/:id/client/profile", handlers.GetClientProfile(conn))
	}

	// Outras Leituras Públicas
	router.GET("/reviews", handlers.GetReviews(conn))
	router.GET("/reviews/:id", handlers.GetReview(conn))
	router.GET("/proposals", handlers.GetProposalLetters(conn))
	router.GET("/proposals/:id", handlers.GetProposalLetter(conn))
	router.GET("/services", handlers.GetServices(conn))
	router.GET("/services/:id", handlers.GetService(conn))
	router.GET("/skills", handlers.GetSkills(conn))
	router.GET("/certifications", handlers.GetCertifications(conn))
	router.GET("/certifications/:id", handlers.GetCertification(conn))

	// --- ROTAS AUTENTICADAS (Escrita e Sensíveis) ---
	authenticated := router.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// Auth Privado
		authenticated.POST("/auth/logout", handlers.Logout(conn))
		authenticated.POST("/auth/change-password", handlers.ChangePassword(conn))
		authenticated.GET("/auth/me", handlers.GetCurrentUser(conn))

		users := authenticated.Group("/users")
		{
			users.PUT("/:id", handlers.UpdateUser(conn))
			users.PATCH("/:id", handlers.UpdateUserPartial(conn))
			users.DELETE("/:id", handlers.DeleteUser(conn))

			profile := users.Group("/:id/profile")
			{
				profile.POST("", handlers.CreateUserProfile(conn))
				profile.PUT("", handlers.UpdateUserProfile(conn))
				profile.DELETE("", handlers.DeleteUserProfile(conn))
				profile.POST("/picture", handlers.UploadProfilePicture(conn))
			}

			analyst := users.Group("/:id/analyst")
			{
				analyst.POST("", handlers.CreateAnalyst(conn))
				analyst.PUT("", handlers.UpdateAnalyst(conn))
				analyst.PATCH("", handlers.UpdateAnalystPartial(conn))
				analyst.DELETE("", handlers.DeleteAnalyst(conn))
				analyst.POST("/skills", handlers.CreateAnalystSkill(conn))
				analyst.DELETE("/skills", handlers.DeleteAnalystSkill(conn))
				analyst.POST("/certifications", handlers.CreateAnalystCertification(conn))
				analyst.DELETE("/certifications", handlers.DeleteAnalystCertification(conn))

				analystProfile := analyst.Group("/profile")
				{
					analystProfile.POST("", handlers.CreateAnalystProfile(conn))
					analystProfile.PUT("", handlers.UpdateAnalystProfile(conn))
					analystProfile.DELETE("", handlers.DeleteAnalystProfile(conn))
				}
			}

			client := users.Group("/:id/client")
			{
				client.POST("", handlers.CreateClient(conn))
				client.PUT("", handlers.UpdateClient(conn))
				client.PATCH("", handlers.UpdateClientPartial(conn))
				client.DELETE("", handlers.DeleteClient(conn))

				clientProfile := client.Group("/profile")
				{
					clientProfile.POST("", handlers.CreateClientProfile(conn))
					clientProfile.PUT("", handlers.UpdateClientProfile(conn))
					clientProfile.DELETE("", handlers.DeleteClientProfile(conn))
				}
			}
		}

		// Escrita em Coleções
		authenticated.POST("/reviews", handlers.CreateReview(conn))
		authenticated.PUT("/reviews/:id", handlers.UpdateReview(conn))
		authenticated.PATCH("/reviews/:id", handlers.UpdateReviewPartial(conn))
		authenticated.DELETE("/reviews/:id", handlers.DeleteReview(conn))

		authenticated.POST("/proposals", handlers.CreateProposalLetter(conn))
		authenticated.PUT("/proposals/:id", handlers.UpdateProposalLetter(conn))
		authenticated.PATCH("/proposals/:id", handlers.UpdateProposalLetterPartial(conn))
		authenticated.DELETE("/proposals/:id", handlers.DeleteProposalLetter(conn))

		authenticated.POST("/services", handlers.CreateService(conn))
		authenticated.PUT("/services/:id", handlers.UpdateService(conn))
		authenticated.PATCH("/services/:id", handlers.UpdateServicePartial(conn))
		authenticated.DELETE("/services/:id", handlers.DeleteService(conn))

		authenticated.POST("/skills", handlers.CreateSkill(conn))
		authenticated.PUT("/skills/:id", handlers.UpdateSkill(conn))
		authenticated.DELETE("/skills/:id", handlers.DeleteSkill(conn))

		authenticated.POST("/certifications", handlers.CreateCertification(conn))
		authenticated.PUT("/certifications/:id", handlers.UpdateCertification(conn))
		authenticated.PATCH("/certifications/:id", handlers.UpdateCertificationPartial(conn))
		authenticated.DELETE("/certifications/:id", handlers.DeleteCertification(conn))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
