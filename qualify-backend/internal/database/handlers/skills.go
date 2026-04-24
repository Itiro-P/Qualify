package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Skill Handlers

// GetSkills godoc
// @Summary Listar habilidades
// @Description Retorna lista de habilidades
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param name query string false "Nome parcial"
// @Success 200 {object} pkg.SkillsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills [get]
func GetSkills(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT id, name FROM skill WHERE 1=1`
		args := []interface{}{}
		argCounter := 1

		if name := c.Query("name"); name != "" {
			query += fmt.Sprintf(" AND name ILIKE $%d", argCounter)
			args = append(args, "%"+name+"%")
			argCounter++
		}

		query += " ORDER BY name ASC"

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var skills []pkg.Skill
		for rows.Next() {
			var skill pkg.Skill
			if err := rows.Scan(&skill.Id, &skill.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			skills = append(skills, skill)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.SkillsResponse{Skills: skills, Count: len(skills)})
	}
}

// GetSkill godoc
// @Summary Obter habilidade
// @Description Retorna uma habilidade pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID da habilidade"
// @Success 200 {object} pkg.SkillResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [get]
func GetSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		skillID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}
		var skill pkg.Skill
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id, name FROM skill WHERE id = $1`, skillID).
			Scan(&skill.Id, &skill.Name)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.SkillResponse{Skill: skill})
	}
}

// CreateSkill godoc
// @Summary Criar habilidade
// @Description Cria uma nova habilidade
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param skill body pkg.Skill true "Objeto habilidade"
// @Success 201 {object} pkg.SkillResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills [post]
func CreateSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var skill pkg.Skill
		if err := c.BindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO skill (name)
			 VALUES ($1)
			 RETURNING id`,
			skill.Name).
			Scan(&skill.Id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pkg.SkillResponse{Skill: skill})
	}
}

// UpdateSkill godoc
// @Summary Atualizar habilidade
// @Description Atualiza uma habilidade pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID da habilidade"
// @Param skill body pkg.Skill true "Objeto habilidade"
// @Success 200 {object} pkg.SkillResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [put]
func UpdateSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		skillID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}
		var skill pkg.Skill
		if err := c.BindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE skill SET name = $1
			 WHERE id = $2
			 RETURNING id, name`,
			skill.Name, skillID).
			Scan(&skill.Id, &skill.Name)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.SkillResponse{Skill: skill})
	}
}

// DeleteSkill godoc
// @Summary Excluir habilidade
// @Description Remove uma habilidade pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID da habilidade"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [delete]
func DeleteSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		skillID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM skill WHERE id = $1`, skillID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "skill deleted successfully"})
	}
}

// Analyst Skill Handlers

// GetAnalystSkills godoc
// @Summary Obter habilidades do analista
// @Description Retorna as habilidades de um analista pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Success 200 {object} pkg.AnalystSkillsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/skills [get]
func GetAnalystSkills(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		analystID := c.Param("id")
		analystIDVal, err := strconv.Atoi(analystID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analyst id"})
			return
		}

		query := `SELECT analyst_id, skill_id FROM analyst_skill WHERE analyst_id = $1 ORDER BY skill_id`

		rows, err := conn.Query(c.Request.Context(), query, analystIDVal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var skills []pkg.AnalystSkill
		for rows.Next() {
			var skill pkg.AnalystSkill
			if err := rows.Scan(&skill.Analyst_id, &skill.Skill_id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			skills = append(skills, skill)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystSkillsResponse{Analyst_skills: skills, Count: len(skills)})
	}
}

// CreateAnalystSkill godoc
// @Summary Criar habilidade para o analista
// @Description Cria uma nova habilidade para um analista pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Success 200 {object} pkg.AnalystSkillResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/skills [post]
func CreateAnalystSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var skill pkg.AnalystSkill
		if err := c.BindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		_, err := conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_skill (analyst_id, skill_id)
			 VALUES ($1, $2)`,
			skill.Analyst_id, skill.Skill_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pkg.AnalystSkillResponse{Analyst_skill: skill})
	}
}

// DeleteAnalystSkill godoc
// @Summary Excluir habilidade do analista
// @Description Exclui uma habilidade de um analista pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/skills [delete]
func DeleteAnalystSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		analystID := c.Param("id")
		analystIDVal, err := strconv.Atoi(analystID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analyst id"})
			return
		}

		skillID := c.Query("skill_id")
		if skillID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "skill_id query parameter required"})
			return
		}
		skillIDVal, err := strconv.Atoi(skillID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM analyst_skill WHERE analyst_id = $1 AND skill_id = $2`,
			analystIDVal, skillIDVal)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "analyst skill not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "analyst skill deleted successfully"})
	}
}

// User Skill Handlers

// GetUserSkills godoc
// @Summary Obter habilidades do usuário
// @Description Retorna as habilidades de um usuário pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.UserSkillsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/skills [get]
func GetUserSkills(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		userIDVal, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		query := `SELECT user_id, skill_id FROM user_skill WHERE user_id = $1 ORDER BY skill_id`

		rows, err := conn.Query(c.Request.Context(), query, userIDVal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var skills []pkg.UserSkill
		for rows.Next() {
			var skill pkg.UserSkill
			if err := rows.Scan(&skill.User_id, &skill.Skill_id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			skills = append(skills, skill)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.UserSkillsResponse{User_skills: skills, Count: len(skills)})
	}
}

// CreateUserSkill godoc
// @Summary Criar habilidade do usuário
// @Description Cria uma habilidade para um usuário pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.UserSkillsResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/skills [post]
func CreateUserSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var skill pkg.UserSkill
		if err := c.BindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		skill.User_id = userID

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO user_skill (user_id, skill_id)
			 VALUES ($1, $2)`,
			skill.User_id, skill.Skill_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pkg.UserSkillResponse{User_skill: skill})
	}
}

// DeleteUserSkill godoc
// @Summary Excluir habilidade do usuário
// @Description Exclui uma habilidade de um usuário pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/skills [delete]
func DeleteUserSkill(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		userIDVal, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		skillID := c.Query("skill_id")
		if skillID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "skill_id query parameter required"})
			return
		}
		skillIDVal, err := strconv.Atoi(skillID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM user_skill WHERE user_id = $1 AND skill_id = $2`,
			userIDVal, skillIDVal)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user skill not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user skill deleted successfully"})
	}
}
