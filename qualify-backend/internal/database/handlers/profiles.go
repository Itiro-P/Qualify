package handlers

import (
	"errors"
	"fmt"
	"main/pkg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/profile [get]
func GetUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "ID inválido"))
			return
		}
		var profile pkg.UserProfile
		err = conn.QueryRow(c.Request.Context(),
			`SELECT user_id, biography, COALESCE(picture, 'default_picture.png') FROM user_profile WHERE user_id = $1`,
			userID).Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Perfil não encontrado"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/profile [post]
func CreateUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := strconv.Atoi(id)

		var profile pkg.UserProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Verificando se o perfil existe
		var exists bool
		err := conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM user_profile WHERE id = $1)`,
			id,
		).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Profile already exists"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`INSERT INTO user_profile (user_id, biography) VALUES ($1, $2)
             RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')`,
			userID, profile.Biography).
			Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		c.JSON(http.StatusCreated, pkg.UserProfileResponse{User_profile: profile})
	}
}

// UploadProfilePicture godoc
// @Summary Upload da foto de perfil
// @Description Realiza o upload de uma imagem e retorna o perfil atualizado
// @Tags Perfis
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID do usuário"
// @Param picture formData file true "Arquivo de imagem"
// @Success 200 {object} pkg.UserProfileResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/profile/picture [post]
func UploadProfilePicture(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := strconv.Atoi(id)

		file, err := c.FormFile("picture")
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "File not sent"))
			return
		}

		var oldFileName *string
		_ = conn.QueryRow(c.Request.Context(), "SELECT picture FROM user_profile WHERE user_id = $1", userID).Scan(&oldFileName)

		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
		savePath := filepath.Join("/app/uploads", newFileName)
		dbPath := "/uploads/" + newFileName

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Fail when saving image file"))
			return
		}

		var profile pkg.UserProfile
		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET picture = $1 WHERE user_id = $2 
			RETURNING user_id, biography, picture`,
			dbPath, userID).Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Erro ao atualizar banco"))
			return
		}

		if oldFileName != nil && *oldFileName != "" && *oldFileName != "/uploads/default_picture.png" {
			_ = os.Remove("." + *oldFileName)
		}

		c.JSON(http.StatusOK, pkg.UserProfileResponse{User_profile: profile})
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/profile [put]
func UpdateUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}
		var profile pkg.UserProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Validando parâmetros obrigatórios
		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty biography"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			 WHERE user_id = $2
			 RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')`,
			profile.Biography, userID).
			Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/profile [delete]
func DeleteUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var fileName string
		_ = conn.QueryRow(c.Request.Context(),
			"SELECT COALESCE(picture, '') FROM user_profile WHERE user_id = $1", id).Scan(&fileName)

		if fileName != "" && fileName != "default_picture.png" {
			_ = os.Remove(filepath.Join("uploads", fileName))
		}

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM user_profile WHERE user_id = $1`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Perfil não encontrado"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Perfil removido"})
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/profile [get]
func GetAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var profile pkg.AnalystProfile

		err := conn.QueryRow(c.Request.Context(),
			`SELECT u.user_id, u.biography, COALESCE(u.picture, 'default_picture.png') 
			 FROM user_profile u
			 INNER JOIN analyst_profile a ON u.user_id = a.analyst_id 
			 WHERE u.user_id = $1`,
			id).Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analista não encontrado"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/profile [post]
func CreateAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := strconv.Atoi(id)

		var profile pkg.AnalystProfile
		_ = c.BindJSON(&profile)

		tx, err := conn.Begin(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Transaction error"))
			return
		}
		defer tx.Rollback(c.Request.Context())

		err = tx.QueryRow(c.Request.Context(),
			`INSERT INTO user_profile (user_id, biography) VALUES ($1, $2)
            RETURNING user_id, biography, picture`,
			userID, profile.Biography).Scan(&profile.User_id, &profile.Biography, &profile.Picture)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case "23505": // unique_violation
					c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Profile already exists for this user"))
					return
				case "23503": // foreign_key_violation (user_id não existe)
					c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "User not found"))
					return
				}
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Error when creating base profile"))
			return
		}

		_, err = tx.Exec(c.Request.Context(),
			`INSERT INTO analyst_profile (analyst_id) VALUES ($1)`, userID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Analyst profile already exists for this user"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Error when creating analyst profile"))
			return
		}

		tx.Commit(c.Request.Context())
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/profile [put]
func UpdateAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}
		var profile pkg.AnalystProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Validando parâmetros obrigatórios
		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty biography"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			 WHERE user_id = $2
			 RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')`,
			profile.Biography, userID).
			Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/profile [delete]
func DeleteAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var fileName string
		_ = conn.QueryRow(c.Request.Context(),
			"SELECT COALESCE(picture, '') FROM user_profile WHERE user_id = $1", id).Scan(&fileName)

		if fileName != "" && fileName != "default_picture.png" {
			_ = os.Remove(filepath.Join("uploads", fileName))
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM user_profile WHERE user_id = $1`,
			userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst profile not found"))
			return
		}

		c.Status(http.StatusNoContent)
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/client/profile [get]
func GetClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var profile pkg.ClientProfile

		err := conn.QueryRow(c.Request.Context(),
			`SELECT u.user_id, u.biography, COALESCE(u.picture, 'default_picture.png') 
			 FROM user_profile u
			 INNER JOIN client_profile c_tab ON u.user_id = c_tab.client_id 
			 WHERE u.user_id = $1`,
			id).Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Cliente não encontrado"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client/profile [post]
func CreateClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := strconv.Atoi(id)

		var profile pkg.AnalystProfile
		_ = c.BindJSON(&profile)

		tx, err := conn.Begin(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Transaction error"))
			return
		}
		defer tx.Rollback(c.Request.Context())

		err = tx.QueryRow(c.Request.Context(),
			`INSERT INTO user_profile (user_id, biography) VALUES ($1, $2)
            RETURNING user_id, biography, picture`,
			userID, profile.Biography).Scan(&profile.User_id, &profile.Biography, &profile.Picture)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case "23505": // unique_violation
					c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Profile already exists for this user"))
					return
				case "23503": // foreign_key_violation (user_id não existe)
					c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "User not found"))
					return
				}
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Error when creating base profile"))
			return
		}

		_, err = tx.Exec(c.Request.Context(),
			`INSERT INTO client_profile (client_id) VALUES ($1)`, userID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Client profile already exists for this user"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Error when creating analyst profile"))
			return
		}

		tx.Commit(c.Request.Context())
		c.JSON(http.StatusCreated, pkg.AnalystProfileResponse{Analyst_profile: profile})
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client/profile [put]
func UpdateClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}
		var profile pkg.ClientProfile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			WHERE user_id = $2
			RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')`,
			profile.Biography, userID).
			Scan(&profile.User_id, &profile.Biography, &profile.Picture)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client/profile [delete]
func DeleteClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var fileName string
		_ = conn.QueryRow(c.Request.Context(),
			"SELECT COALESCE(picture, '') FROM user_profile WHERE user_id = $1", id).Scan(&fileName)

		if fileName != "" && fileName != "default_picture.png" {
			_ = os.Remove(filepath.Join("uploads", fileName))
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM user_profile WHERE user_id = $1`,
			userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client profile not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
