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

func TestReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Reginaldo Caminhos 8=D",
		Email:         "reginald0@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999999",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	analystToken, analystUserID := registerAndLogin(t, router, analystUser)

	clientUser := pkg.UserRegister{
		Name:          "Elma Maria Aquino Pinto",
		Email:         "elma@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41965899556",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	clientToken, clientUserID := registerAndLogin(t, router, clientUser)

	// IDs para encadeamento
	var analystID int
	var clientID int
	var proposalID int
	var serviceID int
	var reviewID int

	t.Run("Criando toda a cadeia até a review", func(t *testing.T) {
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

		// 3. Criar Proposta (Cliente -> Analista)
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

		// 4. Criar Serviço (Baseado na Proposta)
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

		// 5. Criar Review (Finalmente!)
		rData := pkg.Review{
			Rating: 4, Comment: "Arrumaram muito bem minha traseira",
			Service_id: serviceID, Analyst_id: analystID, Client_id: clientID,
		}
		body, _ = json.Marshal(rData)
		req = authRequest(http.MethodPost, "/reviews", body, clientToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var rResp pkg.ReviewResponse
		json.Unmarshal(w.Body.Bytes(), &rResp)
		reviewID = rResp.Review.Id
	})

	t.Run("Alterando review", func(t *testing.T) {
		newRating := 2
		newComment := "Só arrombaram mais a minha traseira"
		patchBody := map[string]interface{}{"rating": newRating, "comment": newComment}
		body, _ := json.Marshal(patchBody)

		req := authRequest(http.MethodPatch, fmt.Sprintf("/reviews/%d", reviewID), body, clientToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pkg.ReviewResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, newRating, resp.Review.Rating)
		assert.Equal(t, newComment, resp.Review.Comment)
	})

	t.Run("Removendo a cadeia completa", func(t *testing.T) {
		// Deletar na ordem inversa para evitar erros de Foreign Key
		entities := []struct {
			url   string
			token string
		}{
			{fmt.Sprintf("/reviews/%d", reviewID), clientToken},
			{fmt.Sprintf("/services/%d", serviceID), clientToken},
			{fmt.Sprintf("/proposals/%d", proposalID), clientToken},
			{fmt.Sprintf("/users/%d/analyst", analystUserID), analystToken},
			{fmt.Sprintf("/users/%d", analystUserID), analystToken},
			{fmt.Sprintf("/users/%d/client", clientUserID), clientToken},
			{fmt.Sprintf("/users/%d", clientUserID), clientToken},
		}

		for _, e := range entities {
			req := authRequest(http.MethodDelete, e.url, nil, e.token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNoContent, w.Code)
		}
	})
}
