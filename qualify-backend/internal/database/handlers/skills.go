package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /skills [get]
func GetSkills(conn *pgxpool.Pool) gin.HandlerFunc {
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
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		var skills []pkg.Skill
		for rows.Next() {
			skill, err := pkg.ScanSkill(rows)
			if pkg.HandleErr(c, err) {
				return
			}
			skills = append(skills, skill)
		}

		if err = rows.Err(); pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /skills/{id} [get]
func GetSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		skill, err := pkg.ScanSkill(conn.QueryRow(c.Request.Context(), `SELECT id, name FROM skill WHERE id = $1`, id))

		if pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /skills [post]
func CreateSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var skill pkg.Skill
		if err := c.BindJSON(&skill); pkg.HandleErr(c, err) {
			return
		}

		var exists bool
		err := conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM skill WHERE name = $1)`, skill.Name).Scan(&exists)

		if pkg.HandleErr(c, err) {
			return
		}
		if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Skill already exists"))
			return
		}

		skill, err = pkg.ScanSkill(conn.QueryRow(c.Request.Context(),
			`INSERT INTO skill (name) VALUES ($1) RETURNING id, name`, skill.Name))

		if pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /skills/{id} [put]
func UpdateSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var skill pkg.Skill
		if err := c.BindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		if skill.Name == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty name"))
			return
		}

		skill, err = pkg.ScanSkill(conn.QueryRow(c.Request.Context(),
			`UPDATE skill SET name = $1 WHERE id = $2 RETURNING id, name`, skill.Name, id))

		if pkg.HandleErr(c, err) {
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /skills/{id} [delete]
func DeleteSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM skill WHERE id = $1`, id)

		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Skill not found"))
			return
		}

		c.Status(http.StatusNoContent)
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
// @Success 200 {object} pkg.SkillsResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/skills [get]
func GetAnalystSkills(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		rows, err := conn.Query(c.Request.Context(),
			`SELECT s.id, s.name
             FROM skill s
             JOIN analyst_skill ac ON s.id = ac.skill_id
             WHERE ac.analyst_id = $1
             ORDER BY s.name`,
			id)

		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		var skills []pkg.Skill
		for rows.Next() {
			skill, err := pkg.ScanSkill(rows)
			if pkg.HandleErr(c, err) {
				return
			}
			skills = append(skills, skill)
		}

		if err = rows.Err(); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.SkillsResponse{Skills: skills, Count: len(skills)})
	}
}

// AssociateAnalystSkill godoc
// @Summary Associar skill existente ao analista
// @Description Associa uma skill já existente a um analista pelo ID
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param skill_id path int true "ID da skill"
// @Success 201 {object} pkg.SkillResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/analyst/skills/{skill_id} [post]
func AssociateAnalystSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		skillID, err := pkg.ParsePathParam(c, "skill_id")
		if err != nil {
			return
		}

		skill, err := pkg.ScanSkill(conn.QueryRow(c.Request.Context(),
			`SELECT id, name FROM skill WHERE id = $1`, skillID))
		if pkg.HandleErr(c, err) {
			return
		}

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_skill (analyst_id, skill_id) VALUES ($1, $2)`,
			id, skillID)

		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.SkillResponse{Skill: skill})
	}
}

// CreateAnalystSkill godoc
// @Summary Adicionar skill ao analista
// @Description Cria uma skill (se não existir) e a associa a um analista
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param skill body pkg.Skill true "Objeto skill (envie apenas `name`)"
// @Success 201 {object} pkg.SkillResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/analyst/skills [post]
func CreateAnalystSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var skill pkg.Skill
		if err := c.BindJSON(&skill); pkg.HandleErr(c, err) {
			return
		}

		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, id).Scan(&analystExists)

		if pkg.HandleErr(c, err) {
			return
		}
		if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst not found"))
			return
		}

		// Busca ou cria a skill
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id FROM skill WHERE name = $1`, skill.Name).Scan(&skill.Id)

		if err == pgx.ErrNoRows {
			skill, err = pkg.ScanSkill(conn.QueryRow(c.Request.Context(),
				`INSERT INTO skill (name) VALUES ($1) RETURNING id, name`, skill.Name))
			if pkg.HandleErr(c, err) {
				return
			}
		} else if pkg.HandleErr(c, err) {
			return
		}

		// Associa a skill ao analyst
		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_skill (analyst_id, skill_id) VALUES ($1, $2)`,
			id, skill.Id)
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.SkillResponse{Skill: skill})
	}
}

// DeleteAnalystSkill godoc
// @Summary Excluir habilidade do analista
// @Description Exclui uma habilidade de um analista pelo ID e id da habilidade
// @Tags Habilidades
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param skill_id query int true "ID da skill"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/skills [delete]
func DeleteAnalystSkill(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		skillID, err := pkg.ParsePathQuery(c, "skill_id")

		if err != nil {
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM analyst_skill WHERE analyst_id = $1 AND skill_id = $2`, id, skillID)

		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst skill not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
