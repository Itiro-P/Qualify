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

func TestProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Reginaldo Ré",
		Email:         "reginaldo@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999989",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	analystToken, analystUserID := registerAndLogin(t, router, analystUser)

	clientUser := pkg.UserRegister{
		Name:          "Marcos Calvaro",
		Email:         "markos@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999989",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	clientToken, clientUserID := registerAndLogin(t, router, clientUser)

	t.Run("Fluxo de Perfil do Analista", func(t *testing.T) {
		// Criar o Analista primeiro
		analyst := pkg.Analyst{Hourly_rate: 100.0}
		body, _ := json.Marshal(analyst)
		req := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst", analystUserID), body, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Criar Perfil
		profile := pkg.AnalystProfile{
			UserProfile: pkg.UserProfile{Biography: "Initial analyst biography"},
		}
		body, _ = json.Marshal(profile)
		req = authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst/profile", analystUserID), body, analystToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Alterar Perfil
		newBio := "Nova bio do Reginaldo"
		patchBody := map[string]interface{}{"biography": newBio}
		body, _ = json.Marshal(patchBody)
		req = authRequest(http.MethodPut, fmt.Sprintf("/users/%d/analyst/profile", analystUserID), body, analystToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp pkg.AnalystProfileResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, newBio, resp.Analyst_profile.Biography)

		// Remover Perfil
		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d/analyst/profile", analystUserID), nil, analystToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Fluxo de Perfil do Cliente", func(t *testing.T) {
		// Criar o Cliente primeiro
		client := pkg.Client{Proposed_budget: 500.0}
		body, _ := json.Marshal(client)
		req := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/client", clientUserID), body, clientToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Criar Perfil
		profile := pkg.ClientProfile{
			UserProfile: pkg.UserProfile{Biography: "Initial client biography"},
		}
		body, _ = json.Marshal(profile)
		req = authRequest(http.MethodPost, fmt.Sprintf("/users/%d/client/profile", clientUserID), body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Alterar Perfil
		newBio := "Nova bio do Marcos"
		patchBody := map[string]interface{}{"biography": newBio}
		body, _ = json.Marshal(patchBody)
		req = authRequest(http.MethodPut, fmt.Sprintf("/users/%d/client/profile", clientUserID), body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp pkg.ClientProfileResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, newBio, resp.Client_profile.Biography)

		// Remover Perfil
		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d/client/profile", clientUserID), nil, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Cleanup de usuários", func(t *testing.T) {
		reqA := authRequest(http.MethodDelete, fmt.Sprintf("/users/%d", analystUserID), nil, analystToken)
		router.ServeHTTP(httptest.NewRecorder(), reqA)

		reqC := authRequest(http.MethodDelete, fmt.Sprintf("/users/%d", clientUserID), nil, clientToken)
		router.ServeHTTP(httptest.NewRecorder(), reqC)
	})
}
