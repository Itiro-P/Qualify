package handlers

import (
	"main/internal/database/services"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
	"golang.org/x/crypto/bcrypt"
)

const userSelect = `id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone`
const userFrom = `"user"`

// GetUser godoc
// @Summary Obter usuário
// @Description Retorna um usuário pelo ID
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.UserResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id} [get]
func GetUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(userSelect).
			From(userFrom).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		user, err := pkg.ScanUser(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.UserResponse{User: user})
	}
}

// GetCurrentUser godoc
// @Summary Obter usuário atual
// @Description Retorna os dados do usuário autenticado a partir do token JWT
// @Tags Usuários
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} pkg.UserResponse
// @Failure 401 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /auth/me [get]
func GetCurrentUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, pkg.Unauthorized(c.FullPath(), "user_id not found in token"))
			return
		}

		userID, ok := userIDValue.(int)
		if !ok {
			c.JSON(http.StatusUnauthorized, pkg.Unauthorized(c.FullPath(), "invalid user_id in token"))
			return
		}

		query, args, err := squirrel.Select(userSelect).
			From(userFrom).Where(squirrel.Eq{"id": userID}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		user, err := pkg.ScanUser(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.UserResponse{User: user})
	}
}

// CreateUser godoc
// @Summary Criar usuário
// @Description Registra e cria um novo usuário
// @Tags Usuários
// @Accept json
// @Produce json
// @Param user body pkg.UserRegister true "Objeto do usuário"
// @Success 201 {object} pkg.UserResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /register [post]
func CreateUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reg pkg.UserRegister
		if err := c.ShouldBindJSON(&reg); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
		if pkg.HandleErr(c, err) {
			return
		}

		user := pkg.User{
			Name:          reg.Name,
			Email:         reg.Email,
			Phone:         reg.Phone,
			Password_hash: string(hashedPassword),
			Country_code:  reg.Country_code,
			Country_name:  reg.Country_name,
			Country_state: reg.Country_state,
			City:          reg.City,
			Timezone:      reg.Timezone,
		}
		if err = services.CreateUser(c.Request.Context(), conn, &user); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.UserResponse{User: user})
	}
}

// UpdateUser godoc
// @Summary Atualizar usuário
// @Description Atualiza um usuário pelo ID
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param user body pkg.User true "Objeto do usuário"
// @Success 200 {object} pkg.UserResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id} [put]
func UpdateUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var user pkg.User
		if err := c.BindJSON(&user); pkg.HandleErr(c, err) {
			return
		}

		query, args, err := squirrel.Update(userFrom).
			Set("name", user.Name).
			Set("email", user.Email).
			Set("phone", user.Phone).
			Set("country_code", user.Country_code).
			Set("country_name", user.Country_name).
			Set("country_state", user.Country_state).
			Set("city", user.City).
			Set("timezone", user.Timezone).
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + userSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		user, err = pkg.ScanUser(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.UserResponse{User: user})
	}
}

// UpdateUserPartial godoc
// @Summary Atualizar parcialmente um ou mais dados do usuário
// @Description Atualiza um ou mais dados do usuário pelo ID
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param user body pkg.UserUpdateRequest true "Objeto do usuário"
// @Success 200 {object} pkg.UserResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id} [patch]
func UpdateUserPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var req pkg.UserUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		builder := squirrel.Update(userFrom).PlaceholderFormat(squirrel.Dollar)
		hasFields := false

		if req.Name != nil {
			builder = builder.Set("name", *req.Name)
			hasFields = true
		}
		if req.Email != nil {
			builder = builder.Set("email", *req.Email)
			hasFields = true
		}
		if req.Phone != nil {
			builder = builder.Set("phone", *req.Phone)
			hasFields = true
		}
		if req.Country_code != nil {
			builder = builder.Set("country_code", *req.Country_code)
			hasFields = true
		}
		if req.Country_name != nil {
			builder = builder.Set("country_name", *req.Country_name)
			hasFields = true
		}
		if req.Country_state != nil {
			builder = builder.Set("country_state", *req.Country_state)
			hasFields = true
		}
		if req.City != nil {
			builder = builder.Set("city", *req.City)
			hasFields = true
		}
		if req.Timezone != nil {
			builder = builder.Set("timezone", *req.Timezone)
			hasFields = true
		}

		if !hasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		query, args, err := builder.
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + userSelect).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		user, err := pkg.ScanUser(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.UserResponse{User: user})
	}
}

// DeleteUser godoc
// @Summary Excluir usuário
// @Description Remove um usuário pelo ID
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id} [delete]
func DeleteUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete(userFrom).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "User not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
