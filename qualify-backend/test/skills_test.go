package test

import (
	"encoding/json"
	"fmt"
	"main/internal/routes"
	"main/pkg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSkill(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Reginaldo Rézinho ain",
		Email:         "reginald0o@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999956799",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	analystToken, analystUserID := registerAndLogin(t, router, analystUser)

	// ID para persistência entre subtestes
	var skillID int

	t.Run("Criando perfil e skill base", func(t *testing.T) {
		// 1. Criar Analista
		analystData := pkg.Analyst{
			Hourly_rate:   100.0,
			Total_reviews: 10,
			Mean_rating:   ToPtrFloat64(4.5),
		}
		body, _ := json.Marshal(analystData)
		req := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst", analystUserID), body, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var aResp pkg.AnalystResponse
		json.Unmarshal(w.Body.Bytes(), &aResp)

		// 2. Criar Skill Global (C++)
		skillData := pkg.Skill{Name: "C++"}
		body, _ = json.Marshal(skillData)
		req = authRequest(http.MethodPost, "/skills", body, analystToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var sResp pkg.SkillResponse
		json.Unmarshal(w.Body.Bytes(), &sResp)
		skillID = sResp.Skill.Id
		assert.Equal(t, "C++", sResp.Skill.Name)
	})

	t.Run("Alterando skill", func(t *testing.T) {
		newName := "aaaa"
		patchBody := map[string]interface{}{"name": newName}
		body, _ := json.Marshal(patchBody)

		req := authRequest(http.MethodPut, fmt.Sprintf("/skills/%d", skillID), body, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var sResp pkg.SkillResponse
		json.Unmarshal(w.Body.Bytes(), &sResp)
		assert.Equal(t, newName, sResp.Skill.Name)
	})

	t.Run("Colocando skill no analista", func(t *testing.T) {
		assocBody := pkg.Skill{Name: "aaaa"} // nome após o update
		body, _ := json.Marshal(assocBody)
		req := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst/skills", analystUserID), body, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Removendo skill do analista", func(t *testing.T) {
		t.Logf("skillID ao remover = %d", skillID)
		targetURL := fmt.Sprintf("/users/%d/analyst/skills?skill_id=%d", analystUserID, skillID)
		req := authRequest(http.MethodDelete, targetURL, nil, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("resposta remoção: %s", w.Body.String())
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Cleanup final", func(t *testing.T) {
		// Remove a Skill
		req := authRequest(http.MethodDelete, fmt.Sprintf("/skills/%d", skillID), nil, analystToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Remove o Analista
		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d/analyst", analystUserID), nil, analystToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Remove o Usuário
		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d", analystUserID), nil, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
