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

func TestService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Reginaldo Caminhos 8=D",
		Email:         "reginaldo00@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41239999999",
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
		Phone:         "41965999556",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	clientToken, clientUserID := registerAndLogin(t, router, clientUser)

	// Variáveis de persistência
	var analystID int
	var clientID int
	var proposalID int
	var serviceID int

	t.Run("Criando toda a cadeia até o serviço", func(t *testing.T) {
		// 1. Criar Analista
		aData := pkg.Analyst{Hourly_rate: 100.0}
		body, _ := json.Marshal(aData)
		req := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst", analystUserID), body, analystToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var aResp pkg.AnalystResponse
		json.Unmarshal(w.Body.Bytes(), &aResp)
		analystID = aResp.Analyst.Id

		// 2. Criar Cliente
		cData := pkg.Client{Proposed_budget: 112.0}
		body, _ = json.Marshal(cData)
		req = authRequest(http.MethodPost, fmt.Sprintf("/users/%d/client", clientUserID), body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var cResp pkg.ClientResponse
		json.Unmarshal(w.Body.Bytes(), &cResp)
		clientID = cResp.Client.Id

		// 3. Criar Proposta
		pData := pkg.ProposalLetter{
			Title: "Arrumem a merda da traseira", Content: "Só arruma cara",
			Proposed_hourly_rate: 70.0, Analyst_id: analystID, Client_id: clientID,
		}
		body, _ = json.Marshal(pData)
		req = authRequest(http.MethodPost, "/proposals", body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var pResp pkg.ProposalLetterResponse
		json.Unmarshal(w.Body.Bytes(), &pResp)
		proposalID = pResp.Proposal_letter.Id

		// 4. Criar Serviço (O clímax do teste)
		sData := pkg.Service{
			Title: "Arrumar a traseira", Content: "Preciso que arrumem minha traseira",
			Hourly_rate: 69.0, Status: "RUNNING", Proposal_letter_id: proposalID,
		}
		body, _ = json.Marshal(sData)
		req = authRequest(http.MethodPost, "/services", body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var sResp pkg.ServiceResponse
		json.Unmarshal(w.Body.Bytes(), &sResp)
		serviceID = sResp.Service.Id
	})

	t.Run("Alterando serviço", func(t *testing.T) {
		newTitle := "Título atualizado por Reginaldo Caminhas"
		newContent := "Só arruma a minha traseira vei"
		patchBody := map[string]interface{}{"title": newTitle, "content": newContent}
		body, _ := json.Marshal(patchBody)

		req := authRequest(http.MethodPatch, fmt.Sprintf("/services/%d", serviceID), body, clientToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pkg.ServiceResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, newTitle, resp.Service.Title)
		assert.Equal(t, newContent, resp.Service.Content)
	})

	t.Run("Removendo a cadeia completa", func(t *testing.T) {
		// Ordem de remoção para respeitar FKs
		steps := []struct {
			url   string
			token string
		}{
			{fmt.Sprintf("/services/%d", serviceID), clientToken},
			{fmt.Sprintf("/proposals/%d", proposalID), clientToken},
			{fmt.Sprintf("/users/%d/analyst", analystUserID), analystToken},
			{fmt.Sprintf("/users/%d", analystUserID), analystToken},
			{fmt.Sprintf("/users/%d/client", clientUserID), clientToken},
			{fmt.Sprintf("/users/%d", clientUserID), clientToken},
		}

		for _, s := range steps {
			req := authRequest(http.MethodDelete, s.url, nil, s.token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNoContent, w.Code)
		}
	})
}
