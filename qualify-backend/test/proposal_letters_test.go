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

func TestProposalLetter(t *testing.T) {
	var analyst = pkg.Analyst{
		User: pkg.User{
			Name:          "Rafael Liberado",
			Email:         "rafael@utfpr.edu.br",
			Phone:         "41996599909",
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
			Name:          "Alyssa Min Ha Lynguissa",
			Email:         "alyssaaa@utfpr.edu.br",
			Phone:         "41965099556",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		Proposed_budget: 112.0,
	}

	var proposal_letter = pkg.ProposalLetter{
		Title:                "Arrumem a merda da traseira pelo AMOR DE AAAAA",
		Content:              "Só arruma cara. Por favor :)",
		Proposed_hourly_rate: 79.0,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	t.Run("Criando carta proposta", func(t *testing.T) {
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

		// Criamos o usuário associado ao cliente
		body, _ = json.Marshal(client.User)
		targetURL = "/users"

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o usuário foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
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

		if clientResponse.Client.User.Id == 0 {
			t.Error("O ID do cliente não deveria ser zero")
		}

		client = clientResponse.Client

		proposal_letter.Analyst_id = analyst.Id
		proposal_letter.Client_id = client.Id

		// Agora SIM criamos a carta proposta
		body, _ = json.Marshal(proposal_letter)
		targetURL = "/proposals"

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se a carta proposta foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var plResponse pkg.ProposalLetterResponse
		json.Unmarshal(w.Body.Bytes(), &plResponse)

		if plResponse.Proposal_letter.Id == 0 {
			t.Error("O ID da carta proposta não deveria ser zero")
		}

		proposal_letter = plResponse.Proposal_letter
	})

	t.Run("Alterando carta proposta", func(t *testing.T) {
		newTitle := "Título atualizado pelo LIBERADOOOOO"
		newContent := "Isso aqui certamente também foi atualizado"

		patchBody := map[string]interface{}{
			"title":   newTitle,
			"content": newContent,
		}

		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/proposals/%d", proposal_letter.Id)

		req := httptest.NewRequest(http.MethodPatch, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var plResponse pkg.ProposalLetterResponse
		err := json.Unmarshal(w.Body.Bytes(), &plResponse)
		assert.NoError(t, err)

		// Verificações
		assert.Equal(t, newTitle, plResponse.Proposal_letter.Title)
		assert.Equal(t, newContent, plResponse.Proposal_letter.Content)

		// Atualiza a variável compartilhada para garantir que o estado está correto
		proposal_letter = plResponse.Proposal_letter
	})

	t.Run("Removendo carta proposta", func(t *testing.T) {
		targetURL := fmt.Sprintf("/proposals/%d", proposal_letter.Id)

		req := httptest.NewRequest(http.MethodDelete, targetURL, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/users/%d/analyst", analyst.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/users/%d", analyst.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/users/%d/client", client.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/users/%d", client.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
