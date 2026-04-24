package routes

import (
	"main/internal/database/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(router *gin.Engine, conn *pgx.Conn) {
	// Users are the base entity. Roles are assigned as nested sub-resources.
	users := router.Group("/users")
	{
		users.POST("", handlers.CreateUser(conn))
		users.GET("/:id", handlers.GetUser(conn))

		user := users.Group("/:id")
		{
			user.PUT("", handlers.UpdateUser(conn))
			user.DELETE("", handlers.DeleteUser(conn))

			// User-owned profile/sub-resources
			user.GET("/profile", handlers.GetUserProfile(conn))
			user.POST("/profile", handlers.CreateUserProfile(conn))
			user.PUT("/profile", handlers.UpdateUserProfile(conn))
			user.DELETE("/profile", handlers.DeleteUserProfile(conn))

			// Role assignment under the user resource
			analystRole := user.Group("/analyst")
			{
				analystRole.POST("", handlers.CreateAnalyst(conn))
				analystRole.GET("", handlers.GetAnalyst(conn))
				analystRole.PUT("", handlers.UpdateAnalyst(conn))
				analystRole.DELETE("", handlers.DeleteAnalyst(conn))
				// Analyst-specific sub-resources

				analystRole.GET("/skills", handlers.GetAnalystSkills(conn))
				analystRole.POST("/skills", handlers.CreateAnalystSkill(conn))
				analystRole.DELETE("/skills", handlers.DeleteAnalystSkill(conn))
				analystRole.GET("/profile", handlers.GetAnalystProfile(conn))
				analystRole.POST("/profile", handlers.CreateAnalystProfile(conn))
				analystRole.PUT("/profile", handlers.UpdateAnalystProfile(conn))
				analystRole.DELETE("/profile", handlers.DeleteAnalystProfile(conn))

				analystRole.GET("/certifications", handlers.GetAnalystCertifications(conn))
				analystRole.POST("/certifications", handlers.CreateAnalystCertification(conn))
				analystRole.DELETE("/certifications", handlers.DeleteAnalystCertification(conn))
			}

			clientRole := user.Group("/client")
			{
				clientRole.POST("", handlers.CreateClient(conn))
				clientRole.GET("", handlers.GetClient(conn))
				clientRole.PUT("", handlers.UpdateClient(conn))
				clientRole.DELETE("", handlers.DeleteClient(conn))

				// Client-specific sub-resources
				clientRole.GET("/profile", handlers.GetClientProfile(conn))
				clientRole.POST("/profile", handlers.CreateClientProfile(conn))
				clientRole.PUT("/profile", handlers.UpdateClientProfile(conn))
				clientRole.DELETE("/profile", handlers.DeleteClientProfile(conn))
			}
		}
	}

	// Top-level collections for search/list endpoints
	router.GET("/analysts", handlers.GetAnalysts(conn))
	router.GET("/clients", handlers.GetClients(conn))

	reviews := router.Group("/reviews")
	{
		reviews.GET("", handlers.GetReviews(conn))
		reviews.GET("/:id", handlers.GetReview(conn))
		reviews.POST("", handlers.CreateReview(conn))
		reviews.PUT("/:id", handlers.UpdateReview(conn))
		reviews.DELETE("/:id", handlers.DeleteReview(conn))
	}

	proposals := router.Group("/proposals")
	{
		proposals.GET("", handlers.GetProposalLetters(conn))
		proposals.GET("/:id", handlers.GetProposalLetter(conn))
		proposals.POST("", handlers.CreateProposalLetter(conn))
		proposals.PUT("/:id", handlers.UpdateProposalLetter(conn))
		proposals.DELETE("/:id", handlers.DeleteProposalLetter(conn))
	}

	services := router.Group("/services")
	{
		services.GET("", handlers.GetServices(conn))
		services.GET("/:id", handlers.GetService(conn))
		services.POST("", handlers.CreateService(conn))
		services.PUT("/:id", handlers.UpdateService(conn))
		services.DELETE("/:id", handlers.DeleteService(conn))
	}

	skills := router.Group("/skills")
	{
		skills.GET("", handlers.GetSkills(conn))
		skills.POST("", handlers.CreateSkill(conn))
		skills.PUT("/:id", handlers.UpdateSkill(conn))
		skills.DELETE("/:id", handlers.DeleteSkill(conn))
	}

	certifications := router.Group("/certifications")
	{
		certifications.GET("", handlers.GetCertifications(conn))
		certifications.POST("", handlers.CreateCertification(conn))
		certifications.PUT("/:id", handlers.UpdateCertification(conn))
		certifications.DELETE("/:id", handlers.DeleteCertification(conn))
	}
}
