package handlers

import (
	"main/internal/database/services"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetUser(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var user pkg.User
		err = conn.QueryRow(c.Request.Context(), "SELECT id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone FROM \"user\" WHERE id = $1", userID).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Time_created, &user.Country_code, &user.Country_name, &user.Country_state, &user.City, &user.Timezone)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func CreateUser(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user pkg.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate required fields
		if user.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name is required"})
			return
		}
		if user.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user email is required"})
			return
		}
		if len(user.Country_code) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "country_code must be exactly 2 characters"})
			return
		}

		err := services.CreateUser(c.Request.Context(), conn, &user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func UpdateUser(conn *pgx.Conn) gin.HandlerFunc {
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

		// Validate required fields
		if user.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name is required"})
			return
		}
		if user.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user email is required"})
			return
		}
		if len(user.Country_code) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "country_code must be exactly 2 characters"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE "user" SET name = $1, email = $2, phone = $3, country_code = $4, 
			 country_name = $5, country_state = $6, city = $7, timezone = $8
			 WHERE id = $9
			 RETURNING id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone`,
			user.Name, user.Email, user.Phone, user.Country_code, user.Country_name, user.Country_state, user.City, user.Timezone, userID).
			Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Time_created, &user.Country_code, &user.Country_name, &user.Country_state, &user.City, &user.Timezone)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(conn *pgx.Conn) gin.HandlerFunc {
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
