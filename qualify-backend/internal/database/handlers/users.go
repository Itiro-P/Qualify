package handlers

import (
	"fmt"
	"main/internal/database/services"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// GetUser godoc
// @Summary Obter usuário
// @Description Retorna um usuário pelo ID
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func GetUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var user pkg.User
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone 
             FROM "user" WHERE id = $1`, userID).Scan(
			&user.Id, &user.Name, &user.Email, &user.Phone, &user.Time_created,
			&user.Country_code, &user.Country_name, &user.Country_state, &user.City, &user.Timezone,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /register [post]
func CreateUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reg pkg.UserRegister
		if err := c.ShouldBindJSON(&reg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obrigatórios ausentes ou inválidos"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar segurança da senha"})
			return
		}

		// 2. Mapear do registro para a struct de banco de dados
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
		err = services.CreateUser(c.Request.Context(), conn, &user)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				c.JSON(http.StatusConflict, gin.H{"error": "Este e-mail já está cadastrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar usuário no banco"})
			return
		}

		// Retorna o usuário criado (sem o hash da senha graças ao tag `json:"-"`)
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /users/{id} [put]
func UpdateUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var user pkg.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE "user" SET name = $1, email = $2, phone = $3, country_code = $4, 
             country_name = $5, country_state = $6, city = $7, timezone = $8
             WHERE id = $9
             RETURNING id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone`,
			user.Name, user.Email, user.Phone, user.Country_code, user.Country_name,
			user.Country_state, user.City, user.Timezone, userID).
			Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Time_created,
				&user.Country_code, &user.Country_name, &user.Country_state, &user.City, &user.Timezone)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /users/{id} [patch]
func UpdateUserPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var req pkg.UserUpdateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		set := []string{}
		args := []interface{}{}
		i := 1

		// Helper para montar a query dinâmica
		addToSet := func(field string, val interface{}) {
			if val != nil {
				set = append(set, fmt.Sprintf("%s = $%d", field, i))
				args = append(args, val)
				i++
			}
		}

		addToSet("name", req.Name)
		addToSet("email", req.Email)
		addToSet("phone", req.Phone)
		addToSet("country_code", req.Country_code)
		addToSet("country_name", req.Country_name)
		addToSet("country_state", req.Country_state)
		addToSet("city", req.City)
		addToSet("timezone", req.Timezone)

		if len(set) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		args = append(args, userID)
		query := fmt.Sprintf(
			`UPDATE "user" SET %s WHERE id = $%d 
             RETURNING id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone`,
			strings.Join(set, ", "), i,
		)

		var user pkg.User
		err = conn.QueryRow(c.Request.Context(), query, args...).Scan(
			&user.Id, &user.Name, &user.Email, &user.Phone, &user.Time_created,
			&user.Country_code, &user.Country_name, &user.Country_state, &user.City, &user.Timezone,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /users/{id} [delete]
func DeleteUser(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM "user" WHERE id = $1`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	}
}
