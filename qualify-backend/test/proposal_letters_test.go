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

func TestProposalLetter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Rafael Liberado",
		Email:         "rafael@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41996599909",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	analystToken, analystUserID := registerAndLogin(t, router, analystUser)

	clientUser := pkg.UserRegister{
		Name:          "Alyssa Min Ha Lynguissa",
		Email:         "alyssaaa@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41965099556",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	clientToken, clientUserID := registerAndLogin(t, router, clientUser)

	// Variáveis para persistência entre subtestes
	var analystID int
	var clientID int
	var proposalID int

	t.Run("Criando entidades e carta proposta", func(t *testing.T) {
		// 1. Criar o Analista (usando o token do analista)
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
		analystID = aResp.Analyst.Id

		// 2. Criar o Cliente (usando o token do cliente)
		clientData := pkg.Client{
			Proposed_budget: 112.0,
		}
		body, _ = json.Marshal(clientData)
		req = authRequest(http.MethodPost, fmt.Sprintf("/users/%d/client", clientUserID), body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var cResp pkg.ClientResponse
		json.Unmarshal(w.Body.Bytes(), &cResp)
		clientID = cResp.Client.Id

		// 3. Criar a Carta Proposta (O cliente envia para o analista)
		proposal_letter := pkg.ProposalLetterCreateRequest{
			Title:                "Arrumem a merda da traseira pelo AMOR DE AAAAA",
			Content:              "Só arruma cara. Por favor :)",
			Proposed_hourly_rate: 79.0,
			Analyst_id:           analystID,
			Client_id:            clientID,
		}
		body, _ = json.Marshal(proposal_letter)
		req = authRequest(http.MethodPost, "/proposals", body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var plResp pkg.ProposalLetterResponse
		json.Unmarshal(w.Body.Bytes(), &plResp)
		assert.NotZero(t, plResp.Proposal_letter.Id)
		proposalID = plResp.Proposal_letter.Id
	})

	t.Run("Alterando carta proposta", func(t *testing.T) {
		newTitle := "Título atualizado pelo LIBERADOOOOO"
		newContent := "Isso aqui certamente também foi atualizado"

		patchBody := map[string]interface{}{
			"title":   newTitle,
			"content": newContent,
		}
		body, _ := json.Marshal(patchBody)

		// Patch exige token (usamos o do cliente que criou a proposta)
		req := authRequest(http.MethodPatch, fmt.Sprintf("/proposals/%d", proposalID), body, clientToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var plResp pkg.ProposalLetterResponse
		json.Unmarshal(w.Body.Bytes(), &plResp)
		assert.Equal(t, newTitle, plResp.Proposal_letter.Title)
	})

	t.Run("Removendo entidades", func(t *testing.T) {
		// 1. Deleta a Proposta
		req := authRequest(http.MethodDelete, fmt.Sprintf("/proposals/%d", proposalID), nil, clientToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// 2. Cleanup Analista (Token do Rafael)
		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d/analyst", analystUserID), nil, analystToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d", analystUserID), nil, analystToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		// 3. Cleanup Cliente (Token da Alyssa)
		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d/client", clientUserID), nil, clientToken)
		router.ServeHTTP(httptest.NewRecorder(), req)

		req = authRequest(http.MethodDelete, fmt.Sprintf("/users/%d", clientUserID), nil, clientToken)
		router.ServeHTTP(httptest.NewRecorder(), req)
	})
}
