package handlers

import (
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User Profile Handlers

// GetUserProfile godoc
// @Summary Obter perfil do usuário
// @Description Retorna o perfil do usuário pelo ID
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.UserProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/profile [get]
func GetUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var profile pkg.UserProfile
		err = conn.QueryRow(c.Request.Context(),
			`SELECT user_id, biography FROM user_profile WHERE user_id = $1`,
			userID).Scan(&profile.User_id, &profile.Biography)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user profile not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.UserProfileResponse{User_profile: profile})
	}
}

// CreateUserProfile godoc
// @Summary Criar perfil do usuário
// @Description Cria o perfil para o usuário especificado
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param profile body pkg.UserProfile true "Objeto perfil"
// @Success 201 {object} pkg.UserProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/profile [post]
func CreateUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var profile pkg.UserProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		profile.User_id = userID

		err = conn.QueryRow(c.Request.Context(),
			`INSERT INTO user_profile (user_id, biography)
			 VALUES ($1, $2)
			 RETURNING user_id`,
			profile.User_id, profile.Biography).
			Scan(&profile.User_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pkg.UserProfileResponse{User_profile: profile})
	}
}

// UpdateUserProfile godoc
// @Summary Atualizar perfil do usuário
// @Description Atualiza o perfil do usuário pelo ID
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param profile body pkg.UserProfile true "Objeto perfil"
// @Success 200 {object} pkg.UserProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/profile [put]
func UpdateUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var profile pkg.UserProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validando parâmetros obrigatórios
		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "biography is required"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			 WHERE user_id = $2
			 RETURNING user_id, biography`,
			profile.Biography, userID).
			Scan(&profile.User_id, &profile.Biography)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user profile not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.UserProfileResponse{User_profile: profile})
	}
}

// DeleteUserProfile godoc
// @Summary Excluir perfil do usuário
// @Description Remove o perfil do usuário
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/profile [delete]
func DeleteUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM user_profile WHERE user_id = $1`,
			userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user profile not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user profile deleted successfully"})
	}
}

// Analyst Profile Handlers

// GetAnalystProfile godoc
// @Summary Obter perfil do analista
// @Description Retorna o perfil do analista pelo ID do usuário
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Success 200 {object} pkg.AnalystProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/profile [get]
func GetAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var profile pkg.AnalystProfile
		err = conn.QueryRow(c.Request.Context(),
			`SELECT user_id, biography FROM user_profile WHERE user_id = $1`,
			userID).Scan(&profile.User_id, &profile.Biography)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "analyst profile not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystProfileResponse{Analyst_profile: profile})
	}
}

// CreateAnalystProfile godoc
// @Summary Criar perfil do analista
// @Description Cria o perfil do analista para o usuário especificado
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Param profile body pkg.AnalystProfile true "Objeto perfil"
// @Success 201 {object} pkg.AnalystProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/profile [post]
func CreateAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var profile pkg.AnalystProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		profile.User_id = userID

		err = conn.QueryRow(c.Request.Context(),
			`INSERT INTO user_profile (user_id, biography)
			 VALUES ($1, $2)
			 RETURNING user_id`,
			profile.User_id, profile.Biography).
			Scan(&profile.User_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pkg.AnalystProfileResponse{Analyst_profile: profile})
	}
}

// UpdateAnalystProfile godoc
// @Summary Atualizar perfil do analista
// @Description Atualiza o perfil do analista pelo ID do usuário
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Param profile body pkg.AnalystProfile true "Objeto perfil"
// @Success 200 {object} pkg.AnalystProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/profile [put]
func UpdateAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var profile pkg.AnalystProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validando parâmetros obrigatórios
		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "biography is required"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			 WHERE user_id = $2
			 RETURNING user_id, biography`,
			profile.Biography, userID).
			Scan(&profile.User_id, &profile.Biography)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "analyst profile not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystProfileResponse{Analyst_profile: profile})
	}
}

// DeleteAnalystProfile godoc
// @Summary Excluir perfil do analista
// @Description Remove o perfil do analista
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/profile [delete]
func DeleteAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM user_profile WHERE user_id = $1`,
			userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "analyst profile not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "analyst profile deleted successfully"})
	}
}

// Client Profile Handlers

// GetClientProfile godoc
// @Summary Obter perfil do cliente
// @Description Retorna o perfil do cliente pelo ID do usuário
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (cliente)"
// @Success 200 {object} pkg.ClientProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client/profile [get]
func GetClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var profile pkg.ClientProfile
		err = conn.QueryRow(c.Request.Context(),
			`SELECT user_id, biography FROM user_profile WHERE user_id = $1`,
			userID).Scan(&profile.User_id, &profile.Biography)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "client profile not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.ClientProfileResponse{Client_profile: profile})
	}
}

// CreateClientProfile godoc
// @Summary Criar perfil do cliente
// @Description Cria o perfil do cliente para o usuário especificado
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (cliente)"
// @Param profile body pkg.ClientProfile true "Objeto perfil"
// @Success 201 {object} pkg.ClientProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client/profile [post]
func CreateClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var profile pkg.ClientProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validando parâmetros obrigatórios
		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "biography is required"})
			return
		}
		profile.User_id = userID

		err = conn.QueryRow(c.Request.Context(),
			`INSERT INTO user_profile (user_id, biography)
			 VALUES ($1, $2)
			 RETURNING user_id`,
			profile.User_id, profile.Biography).
			Scan(&profile.User_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pkg.ClientProfileResponse{Client_profile: profile})
	}
}

// UpdateClientProfile godoc
// @Summary Atualizar perfil do cliente
// @Description Atualiza o perfil do cliente pelo ID do usuário
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (cliente)"
// @Param profile body pkg.ClientProfile true "Objeto perfil"
// @Success 200 {object} pkg.ClientProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client/profile [put]
func UpdateClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var profile pkg.ClientProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			 WHERE user_id = $2
			 RETURNING user_id, biography`,
			profile.Biography, userID).
			Scan(&profile.User_id, &profile.Biography)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "client profile not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.ClientProfileResponse{Client_profile: profile})
	}
}

// DeleteClientProfile godoc
// @Summary Excluir perfil do cliente
// @Description Remove o perfil do cliente
// @Tags Perfis
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (cliente)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client/profile [delete]
func DeleteClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM user_profile WHERE user_id = $1`,
			userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "client profile not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "client profile deleted successfully"})
	}
}
