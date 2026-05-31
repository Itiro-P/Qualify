package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const userProfileSelect = `SELECT user_id, biography, COALESCE(picture, 'default_picture.png')`

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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(userProfileSelect).
			From("user_profile").Where(squirrel.Eq{"user_id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		profile, err := pkg.ScanProfile(conn.QueryRow(c.Request.Context(), query, args...))

		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var profile pkg.UserProfile
		if err := c.BindJSON(&profile); pkg.HandleErr(c, err) {
			return
		}

		// Verificando se o perfil existe
		var exists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM user_profile WHERE id = $1)`, id).Scan(&exists)

		if pkg.HandleErr(c, err) {
			return
		}

		if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Profile already exists"))
			return
		}

		query, args, err := squirrel.Insert("user_profile").Values(id, profile.Biography).
			Suffix(`RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')`).PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		profile, err = pkg.ScanProfile(conn.QueryRow(c.Request.Context(), query, args...))

		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		file, err := c.FormFile("picture")
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "File not sent"))
			return
		}

		var oldFileName *string
		_ = conn.QueryRow(c.Request.Context(),
			`SELECT picture FROM user_profile WHERE user_id = $1`, id).Scan(&oldFileName)

		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("user_%d_%d%s", id, time.Now().Unix(), ext)
		savePath := filepath.Join("/app/uploads", newFileName)
		dbPath := "/uploads/" + newFileName

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Fail when saving image file"))
			return
		}

		query, args, err := squirrel.Update("user_profile").
			Set("picture", dbPath).
			Where(squirrel.Eq{"user_id": id}).
			Suffix("RETURNING user_id, biography, picture").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err := pkg.ScanProfile(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var profile pkg.UserProfile
		if err := c.BindJSON(&profile); pkg.HandleErr(c, err) {
			return
		}

		// Validando parâmetros obrigatórios
		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty biography"))
			return
		}

		profile, err = pkg.ScanProfile(conn.QueryRow(c.Request.Context(),
			`UPDATE user_profile SET biography = $1
			 WHERE user_id = $2
			 RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')`,
			profile.Biography, id))

		if pkg.HandleErr(c, err) {
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/profile [delete]
func DeleteUserProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var fileName string
		_ = conn.QueryRow(c.Request.Context(),
			"SELECT COALESCE(picture, '') FROM user_profile WHERE user_id = $1", id).Scan(&fileName)

		if fileName != "" && fileName != "default_picture.png" {
			_ = os.Remove(filepath.Join("uploads", fileName))
		}

		query, args, err := squirrel.Delete("user_profile").
			Where(squirrel.Eq{"user_id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Profile not found"))
			return
		}

		c.Status(http.StatusNoContent)
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(
			"u.user_id",
			"u.biography",
			"COALESCE(u.picture, 'default_picture.png')").
			From("user_profile u").
			Join("analyst_profile a ON u.user_id = a.analyst_id").
			Where(squirrel.Eq{"u.user_id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err := pkg.ScanAnalystProfile(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var profile pkg.AnalystProfile
		if err := c.BindJSON(&profile); pkg.HandleErr(c, err) {
			return
		}

		var exists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst_profile WHERE analyst_id = $1)`, id).Scan(&exists)
		if pkg.HandleErr(c, err) {
			return
		}
		if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Analyst profile already exists for this user"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer func() { _ = tx.Rollback(c.Request.Context()) }()

		userProfileQuery, userProfileArgs, err := squirrel.Insert("user_profile").
			Columns("user_id", "biography").
			Values(id, profile.Biography).
			Suffix("RETURNING user_id, biography, picture").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err = pkg.ScanAnalystProfile(tx.QueryRow(c.Request.Context(), userProfileQuery, userProfileArgs...))
		if pkg.HandleErr(c, err) {
			return
		}

		analystProfileQuery, analystProfileArgs, err := squirrel.Insert("analyst_profile").
			Columns("analyst_id").
			Values(id).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		if _, err = tx.Exec(c.Request.Context(), analystProfileQuery, analystProfileArgs...); pkg.HandleErr(c, err) {
			return
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/profile [put]
func UpdateAnalystProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var profile pkg.AnalystProfile
		if err := c.BindJSON(&profile); pkg.HandleErr(c, err) {
			return
		}

		if profile.Biography == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty biography"))
			return
		}

		query, args, err := squirrel.Update("user_profile").
			Set("biography", profile.Biography).
			Where(squirrel.Eq{"user_id": id}).
			Suffix("RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err = pkg.ScanAnalystProfile(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var fileName string
		_ = conn.QueryRow(c.Request.Context(),
			"SELECT COALESCE(picture, '') FROM user_profile WHERE user_id = $1", id).Scan(&fileName)

		if fileName != "" && fileName != "default_picture.png" {
			_ = os.Remove(filepath.Join("uploads", fileName))
		}

		query, args, err := squirrel.Delete("user_profile").
			Where(squirrel.Eq{"user_id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(
			"u.user_id",
			"u.biography",
			"COALESCE(u.picture, 'default_picture.png')").
			From("user_profile u").
			Join("client_profile a ON u.user_id = c.client_id").
			Where(squirrel.Eq{"u.user_id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err := pkg.ScanClientProfile(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var profile pkg.ClientProfile
		if err := c.BindJSON(&profile); pkg.HandleErr(c, err) {
			return
		}

		var exists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM client_profile WHERE client_id = $1)`, id).Scan(&exists)
		if pkg.HandleErr(c, err) {
			return
		}
		if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Client profile already exists"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer func() { _ = tx.Rollback(c.Request.Context()) }()

		userProfileQuery, userProfileArgs, err := squirrel.Insert("user_profile").
			Columns("user_id", "biography").
			Values(id, profile.Biography).
			Suffix("RETURNING user_id, biography, picture").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err = pkg.ScanClientProfile(tx.QueryRow(c.Request.Context(), userProfileQuery, userProfileArgs...))
		if pkg.HandleErr(c, err) {
			return
		}

		clientProfileQuery, clientProfileArgs, err := squirrel.Insert("client_profile").
			Columns("client_id").
			Values(id).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		if _, err = tx.Exec(c.Request.Context(), clientProfileQuery, clientProfileArgs...); pkg.HandleErr(c, err) {
			return
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client/profile [put]
func UpdateClientProfile(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var profile pkg.ClientProfile
		if err := c.BindJSON(&profile); pkg.HandleErr(c, err) {
			return
		}

		query, args, err := squirrel.Update("user_profile").
			Set("biography", profile.Biography).
			Where(squirrel.Eq{"user_id": id}).
			Suffix("RETURNING user_id, biography, COALESCE(picture, 'default_picture.png')").
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		profile, err = pkg.ScanClientProfile(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var fileName string
		_ = conn.QueryRow(c.Request.Context(),
			"SELECT COALESCE(picture, '') FROM user_profile WHERE user_id = $1", id).Scan(&fileName)

		if fileName != "" && fileName != "default_picture.png" {
			_ = os.Remove(filepath.Join("uploads", fileName))
		}

		query, args, err := squirrel.Delete("user_profile").
			Where(squirrel.Eq{"user_id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client profile not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
