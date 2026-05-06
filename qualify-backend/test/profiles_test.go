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

func TestProfile(t *testing.T) {
	var analyst = pkg.Analyst{
		User: pkg.User{
			Name:          "Reginaldo Ré",
			Email:         "reginaldo@utfpr.edu.br",
			Phone:         "41999999999",
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

	var client = pkg.Client{
		User: pkg.User{
			Name:          "Marcos Calvaro",
			Email:         "markos@utfpr.edu.br",
			Phone:         "41999999989",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		Proposed_budget: 100.0,
	}

	var analystProfile = pkg.AnalystProfile{
		UserProfile: pkg.UserProfile{
			User_id:   0,
			Biography: "",
		},
	}

	var clientProfile = pkg.ClientProfile{
		UserProfile: pkg.UserProfile{
			User_id:   0,
			Biography: "Initial client biography",
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

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

		analyst = analystResponse.Analyst

		// Agora sim criamos seu perfil
		body, _ = json.Marshal(analystProfile)
		targetURL = fmt.Sprintf("/users/%d/analyst/profile", analyst.User.Id)

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var profileResponse = pkg.AnalystProfileResponse{}
		json.Unmarshal(w.Body.Bytes(), &profileResponse)
		analystProfile = profileResponse.Analyst_profile

		assert.Equal(t, analyst.User.Id, analystProfile.User_id)
	})

	t.Run("Alterando perfil de analista", func(t *testing.T) {
		newBio := "aaaa"
		patchBody := map[string]interface{}{
			"biography": newBio,
		}
		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/users/%d/analyst/profile", analyst.User.Id)

		req := httptest.NewRequest(http.MethodPut, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystProfileResponse pkg.AnalystProfileResponse
		json.Unmarshal(w.Body.Bytes(), &analystProfileResponse)

		if analystProfileResponse.Analyst_profile.Biography != newBio {
			t.Errorf("Valor do campo 'biography' diferente do esperado")
		}
	})

	t.Run("Removendo perfil de analista", func(t *testing.T) {
		targetURL := fmt.Sprintf("/users/%d/analyst/profile", analyst.User.Id)

		req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Criando perfil de cliente", func(t *testing.T) {
		// Primeiro, criamos o usuário associado ao cliente
		body, _ := json.Marshal(client.User)
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
			t.Error("O ID do cliente não deveria ser zero")
		}

		// Agora criamos o cliente usando o ID do usuário criado
		client.User = userResponse.User
		body, _ = json.Marshal(client)
		targetURL = fmt.Sprintf("/users/%d/client", client.User.Id)

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o cliente foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var clientResponse pkg.ClientResponse
		json.Unmarshal(w.Body.Bytes(), &clientResponse)

		if clientResponse.Client.Id == 0 {
			t.Error("O ID do analista não deveria ser zero")
		}

		client = clientResponse.Client

		// Agora sim criamos seu perfil
		body, _ = json.Marshal(clientProfile)
		targetURL = fmt.Sprintf("/users/%d/client/profile", client.User.Id)

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var profileResponse = pkg.ClientProfileResponse{}
		json.Unmarshal(w.Body.Bytes(), &profileResponse)
		clientProfile = profileResponse.Client_profile

		assert.Equal(t, client.User.Id, clientProfile.User_id)
	})

	t.Run("Alterando perfil de cliente", func(t *testing.T) {
		newBio := "aaaa"
		patchBody := map[string]interface{}{
			"biography": newBio,
		}
		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/users/%d/client/profile", clientProfile.User_id)

		req := httptest.NewRequest(http.MethodPut, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientProfileResponse pkg.ClientProfileResponse
		json.Unmarshal(w.Body.Bytes(), &clientProfileResponse)

		if clientProfileResponse.Client_profile.Biography != newBio {
			t.Errorf("Valor do campo 'biography' diferente do esperado")
		}
	})

	t.Run("Removendo perfil de cliente", func(t *testing.T) {
		targetURL := fmt.Sprintf("/users/%d/client/profile", clientProfile.User_id)

		req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
