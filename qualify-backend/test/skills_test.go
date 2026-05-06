package test

import (
	"bytes"
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
	var analyst = pkg.Analyst{
		User: pkg.User{
			Name:          "Reginaldo Rézinho ain",
			Email:         "reginald0o@utfpr.edu.br",
			Phone:         "41999956799",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		Hourly_rate:   100.0,
		Total_reviews: 10,
		Mean_rating:   ToPtrFloat64(4.5),
	}

	var skill = pkg.Skill{
		Id:   0,
		Name: "C++",
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

	var postSkillReponse = pkg.Skill{}
	var postAnalystResponse = pkg.Analyst{}

	t.Run("Criando perfil de analista", func(t *testing.T) {
		// Primeiro, criamos o usuário associado ao analista
		body, _ := json.Marshal(analyst.User)
		targetURL := "/users"

		req := httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o usuário foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var userResponse pkg.UserResponse
		json.Unmarshal(w.Body.Bytes(), &userResponse)

		if userResponse.User.Id == 0 {
			t.Error("O ID do analista não deveria ser zero")
		}

		// Agora criamos o analista usando o ID do usuário criado
		analyst.User = userResponse.User
		body, _ = json.Marshal(analyst)
		targetURL = fmt.Sprintf("/users/%d/analyst", analyst.User.Id)

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o analista foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var analystResponse pkg.AnalystResponse
		json.Unmarshal(w.Body.Bytes(), &analystResponse)

		if analystResponse.Analyst.User.Id == 0 {
			t.Error("O ID do analista não deveria ser zero")
		}

		postAnalystResponse = analystResponse.Analyst

		// Agora sim criamos sua competencia
		body, _ = json.Marshal(skill)
		targetURL = "/skills"

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var skillResponse = pkg.SkillResponse{}
		json.Unmarshal(w.Body.Bytes(), &skillResponse)
		postSkillReponse = skillResponse.Skill

		assert.Equal(t, postSkillReponse.Name, skill.Name)
	})

	t.Run("Alterando skill", func(t *testing.T) {
		newName := "aaaa"
		patchBody := map[string]interface{}{
			"name": newName,
		}
		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/skills/%d", postSkillReponse.Id)

		req := httptest.NewRequest(http.MethodPut, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var skillResponse pkg.SkillResponse
		json.Unmarshal(w.Body.Bytes(), &skillResponse)

		if skillResponse.Skill.Name != newName {
			t.Errorf("Valor do campo 'name' diferente do esperado")
		}
	})

	t.Run("Colocando skill no analista", func(t *testing.T) {
		targetURL := fmt.Sprintf("/users/%d/analyst/skills", postAnalystResponse.User.Id)

		assocBody := pkg.AnalystSkill{
			Skill_id: postSkillReponse.Id,
		}
		body, _ := json.Marshal(assocBody)

		req := httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp pkg.AnalystSkillResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
	})

	t.Run("Removendo skill do analista", func(t *testing.T) {
		targetURL := fmt.Sprintf("/users/%d/analyst/skills?skill_id=%d",
			postAnalystResponse.User.Id, postSkillReponse.Id)

		req := httptest.NewRequest(http.MethodDelete, targetURL, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Removendo skill", func(t *testing.T) {
		targetURL := fmt.Sprintf("/skills/%d", postSkillReponse.Id)

		req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/users/%d/analyst", postAnalystResponse.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/users/%d", postAnalystResponse.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})
}
