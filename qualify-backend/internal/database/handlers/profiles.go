package handlers

import (
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// User Profile Handlers

func GetUserProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, profile)
	}
}

func CreateUserProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusCreated, profile)
	}
}

func UpdateUserProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, profile)
	}
}

func DeleteUserProfile(conn *pgx.Conn) gin.HandlerFunc {
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

func GetAnalystProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, profile)
	}
}

func CreateAnalystProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusCreated, profile)
	}
}

func UpdateAnalystProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, profile)
	}
}

func DeleteAnalystProfile(conn *pgx.Conn) gin.HandlerFunc {
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

func GetClientProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, profile)
	}
}

func CreateClientProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusCreated, profile)
	}
}

func UpdateClientProfile(conn *pgx.Conn) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, profile)
	}
}

func DeleteClientProfile(conn *pgx.Conn) gin.HandlerFunc {
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
