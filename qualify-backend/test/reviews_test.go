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

func TestReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	var analyst = pkg.Analyst{
		User: pkg.User{
			Name:          "Reginaldo Caminhos 8=D",
			Email:         "reginald0@utfpr.edu.br",
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
			Name:          "Alyssa Lynguissa",
			Email:         "alyss4@utfpr.edu.br",
			Phone:         "41965899556",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		Proposed_budget: 112.0,
	}

	var service = pkg.Service{
		Title:       "Arrumar a traseira",
		Content:     "Preciso que arrumem minha traseira",
		Hourly_rate: 69.0,
		Status:      "RUNNING",
	}

	var proposal_letter = pkg.ProposalLetter{
		Title:                "Arrumem a merda da traseira",
		Content:              "Só arruma cara",
		Proposed_hourly_rate: 70.0,
	}

	var review = pkg.Review{
		Rating:  4,
		Comment: "Arrumaram muito bem minha traseira",
	}

	t.Run("Criando review", func(t *testing.T) {
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

		// Agora SIM criamos a carta proposta
		proposal_letter.Analyst_id = analyst.Id
		proposal_letter.Client_id = client.Id

		body, _ = json.Marshal(proposal_letter)
		targetURL = "/proposals"

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o usuário foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var plResponse pkg.ProposalLetterResponse
		json.Unmarshal(w.Body.Bytes(), &plResponse)

		if plResponse.Proposal_letter.Id == 0 {
			t.Error("O ID da carta proposta não deveria ser zero")
		}

		proposal_letter = plResponse.Proposal_letter

		// AGORA SIM CRIAMOS A MERDA DO SERVICO!!!
		service.Proposal_letter_id = proposal_letter.Id
		body, _ = json.Marshal(service)
		targetURL = "/services"

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o servico foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var serviceResponse pkg.ServiceResponse
		json.Unmarshal(w.Body.Bytes(), &serviceResponse)

		if serviceResponse.Service.Id == 0 {
			t.Error("O ID do servico não deveria ser zero")
		}
		service = serviceResponse.Service

		// 	CACETE EU NÃO AGUENTO MAISSS
		review.Service_id = service.Id
		review.Analyst_id = analyst.Id
		review.Client_id = client.Id

		body, _ = json.Marshal(review)
		targetURL = "/reviews"

		req = httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se o servico foi criado com sucesso
		assert.Equal(t, http.StatusCreated, w.Code)
		var reviewResponse pkg.ReviewResponse
		json.Unmarshal(w.Body.Bytes(), &reviewResponse)

		if reviewResponse.Review.Id == 0 {
			t.Error("O ID da review não deveria ser zero")
		}

		review = reviewResponse.Review
	})

	t.Run("Alterando review", func(t *testing.T) {
		newRating := 2
		newComment := "Só arrombaram mais a minha traseira"

		patchBody := map[string]interface{}{
			"rating":  newRating,
			"comment": newComment,
		}

		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/reviews/%d", review.Id)

		req := httptest.NewRequest(http.MethodPatch, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var reviewResponse pkg.ReviewResponse
		json.Unmarshal(w.Body.Bytes(), &reviewResponse)

		if reviewResponse.Review.Id == 0 {
			t.Error("O ID da review não deveria ser zero")
		}

		assert.Equal(t, reviewResponse.Review.Rating, newRating)
		assert.Equal(t, reviewResponse.Review.Comment, newComment)
	})

	t.Run("Removendo review", func(t *testing.T) {
		targetURL := fmt.Sprintf("/reviews/%d", review.Id)

		req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/services/%d", service.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		targetURL = fmt.Sprintf("/proposals/%d", proposal_letter.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)

		w = httptest.NewRecorder()
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
