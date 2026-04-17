package routes

import (
	"main/internal/database/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(router *gin.Engine, conn *pgx.Conn) {
	// User routes
	router.GET("/user/:id", handlers.GetUser(conn))
	router.POST("/user", handlers.CreateUser(conn))
	router.PUT("/user/:id", handlers.UpdateUser(conn))
	router.DELETE("/user/:id", handlers.DeleteUser(conn))

	// Analyst routes
	router.GET("/analyst/:id", handlers.GetAnalyst(conn))
	router.GET("/analysts", handlers.GetAnalysts(conn))
	router.POST("/analyst", handlers.CreateAnalyst(conn))
	router.PUT("/analyst/:id", handlers.UpdateAnalyst(conn))
	router.DELETE("/analyst/:id", handlers.DeleteAnalyst(conn))

	// Client routes
	router.GET("/client/:id", handlers.GetClient(conn))
	router.GET("/clients", handlers.GetClients(conn))
	router.POST("/client", handlers.CreateClient(conn))
	router.PUT("/client/:id", handlers.UpdateClient(conn))
	router.DELETE("/client/:id", handlers.DeleteClient(conn))

	// Review routes
	router.GET("/review/:id", handlers.GetReview(conn))
	router.GET("/reviews", handlers.GetReviews(conn))
	router.POST("/review", handlers.CreateReview(conn))
	router.PUT("/review/:id", handlers.UpdateReview(conn))
	router.DELETE("/review/:id", handlers.DeleteReview(conn))

	// Profile routes
	router.GET("/user/:id/profile", handlers.GetUserProfile(conn))
	router.POST("/user/profile", handlers.CreateUserProfile(conn))
	router.PUT("/user/:id/profile", handlers.UpdateUserProfile(conn))
	router.DELETE("/user/:id/profile", handlers.DeleteUserProfile(conn))

	router.GET("/analyst/:id/profile", handlers.GetAnalystProfile(conn))
	router.POST("/analyst/profile", handlers.CreateAnalystProfile(conn))
	router.PUT("/analyst/:id/profile", handlers.UpdateAnalystProfile(conn))
	router.DELETE("/analyst/:id/profile", handlers.DeleteAnalystProfile(conn))

	router.GET("/client/:id/profile", handlers.GetClientProfile(conn))
	router.POST("/client/profile", handlers.CreateClientProfile(conn))
	router.PUT("/client/:id/profile", handlers.UpdateClientProfile(conn))
	router.DELETE("/client/:id/profile", handlers.DeleteClientProfile(conn))

	// Skill routes
	router.GET("/skill/:id", handlers.GetSkill(conn))
	router.GET("/skills", handlers.GetSkills(conn))
	router.POST("/skill", handlers.CreateSkill(conn))
	router.PUT("/skill/:id", handlers.UpdateSkill(conn))
	router.DELETE("/skill/:id", handlers.DeleteSkill(conn))

	// Analyst skill routes
	router.GET("/analyst/:id/skills", handlers.GetAnalystSkills(conn))
	router.POST("/analyst/skill", handlers.CreateAnalystSkill(conn))
	router.DELETE("/analyst/:id/skill", handlers.DeleteAnalystSkill(conn))

	// User skill routes
	router.GET("/user/:id/skills", handlers.GetUserSkills(conn))
	router.POST("/user/skill", handlers.CreateUserSkill(conn))
	router.DELETE("/user/:id/skill", handlers.DeleteUserSkill(conn))

	// Proposal letter routes
	router.GET("/proposal/:id", handlers.GetProposalLetter(conn))
	router.GET("/proposals", handlers.GetProposalLetters(conn))
	router.POST("/proposal", handlers.CreateProposalLetter(conn))
	router.PUT("/proposal/:id", handlers.UpdateProposalLetter(conn))
	router.DELETE("/proposal/:id", handlers.DeleteProposalLetter(conn))

	// Service routes
	router.GET("/service/:id", handlers.GetService(conn))
	router.GET("/services", handlers.GetServices(conn))
	router.POST("/service", handlers.CreateService(conn))
	router.PUT("/service/:id", handlers.UpdateService(conn))
	router.DELETE("/service/:id", handlers.DeleteService(conn))

	// Certification routes
	router.GET("/certification/:id", handlers.GetCertification(conn))
	router.GET("/certifications", handlers.GetCertifications(conn))
	router.POST("/certification", handlers.CreateCertification(conn))
	router.PUT("/certification/:id", handlers.UpdateCertification(conn))
	router.DELETE("/certification/:id", handlers.DeleteCertification(conn))

	// Analyst certification routes
	router.GET("/analyst/:id/certifications", handlers.GetAnalystCertifications(conn))
	router.POST("/analyst/certification", handlers.CreateAnalystCertification(conn))
	router.DELETE("/analyst/:id/certification", handlers.DeleteAnalystCertification(conn))
}
