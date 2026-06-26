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

func TestConversationRequiresValidServiceAndProposal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Conversation Analyst",
		Email:         "conversationanalyst@utfpr.edu.br",
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
		Name:          "Conversation Client",
		Email:         "conversationclient@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999998",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	clientToken, clientUserID := registerAndLogin(t, router, clientUser)

	var analystID int
	var clientID int
	var proposalID int
	var serviceID int

	t.Run("cria perfil do analista e usa o cliente já criado pelo cadastro", func(t *testing.T) {
		analystBody, _ := json.Marshal(pkg.AnalystCreateRequest{Hourly_rate: 120.0})
		analystReq := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst", analystUserID), analystBody, analystToken)
		analystResp := httptest.NewRecorder()
		router.ServeHTTP(analystResp, analystReq)
		assert.Equal(t, http.StatusCreated, analystResp.Code)
		var analystPayload pkg.AnalystResponse
		json.Unmarshal(analystResp.Body.Bytes(), &analystPayload)
		analystID = analystPayload.Analyst.Id

		clientID = clientUserID
	})

	t.Run("cria proposta e serviço para a conversa", func(t *testing.T) {
		proposalBody, _ := json.Marshal(pkg.ProposalLetterCreateRequest{
			Title:                "Proposta de conversa",
			Content:              "Conteúdo da proposta",
			Analyst_id:           analystID,
			Client_id:            clientID,
			Proposed_hourly_rate: 90.0,
		})
		proposalReq := authRequest(http.MethodPost, "/proposals", proposalBody, clientToken)
		proposalResp := httptest.NewRecorder()
		router.ServeHTTP(proposalResp, proposalReq)
		assert.Equal(t, http.StatusCreated, proposalResp.Code)
		var proposalPayload pkg.ProposalLetterResponse
		json.Unmarshal(proposalResp.Body.Bytes(), &proposalPayload)
		proposalID = proposalPayload.Proposal_letter.Id

		serviceBody, _ := json.Marshal(pkg.ServiceCreateRequest{
			Title:              "Serviço de conversa",
			Content:            "Conteúdo do serviço",
			Hourly_rate:        95.0,
			Status:             "RUNNING",
			Proposal_letter_id: proposalID,
		})
		serviceReq := authRequest(http.MethodPost, "/services", serviceBody, clientToken)
		serviceResp := httptest.NewRecorder()
		router.ServeHTTP(serviceResp, serviceReq)
		assert.Equal(t, http.StatusCreated, serviceResp.Code)
		var servicePayload pkg.ServiceResponse
		json.Unmarshal(serviceResp.Body.Bytes(), &servicePayload)
		serviceID = servicePayload.Service.Id
	})

	t.Run("aceita criação somente quando service_id e proposal_id existem", func(t *testing.T) {
		validBody, _ := json.Marshal(map[string]interface{}{
			"analyst_id":  analystID,
			"client_id":   clientID,
			"service_id":  serviceID,
			"proposal_id": proposalID,
		})
		validReq := authRequest(http.MethodPost, "/conversations", validBody, clientToken)
		validResp := httptest.NewRecorder()
		router.ServeHTTP(validResp, validReq)
		assert.Equal(t, http.StatusCreated, validResp.Code)

		missingProposalBody, _ := json.Marshal(map[string]interface{}{
			"analyst_id": analystID,
			"client_id":  clientID,
			"service_id": serviceID,
		})
		missingProposalReq := authRequest(http.MethodPost, "/conversations", missingProposalBody, clientToken)
		missingProposalResp := httptest.NewRecorder()
		router.ServeHTTP(missingProposalResp, missingProposalReq)
		assert.Equal(t, http.StatusBadRequest, missingProposalResp.Code)

		invalidBody, _ := json.Marshal(map[string]interface{}{
			"analyst_id":  analystID,
			"client_id":   clientID,
			"service_id":  999999,
			"proposal_id": proposalID,
		})
		invalidReq := authRequest(http.MethodPost, "/conversations", invalidBody, clientToken)
		invalidResp := httptest.NewRecorder()
		router.ServeHTTP(invalidResp, invalidReq)
		assert.Equal(t, http.StatusNotFound, invalidResp.Code)
	})

	t.Run("rejeita atualização com service_id ou proposal_id inválidos", func(t *testing.T) {
		createBody, _ := json.Marshal(map[string]interface{}{
			"analyst_id":  analystID,
			"client_id":   clientID,
			"service_id":  serviceID,
			"proposal_id": proposalID,
		})
		createReq := authRequest(http.MethodPost, "/conversations", createBody, clientToken)
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq)
		assert.Equal(t, http.StatusCreated, createResp.Code)

		var createdConversation pkg.Conversation
		json.Unmarshal(createResp.Body.Bytes(), &createdConversation)

		invalidServiceBody, _ := json.Marshal(map[string]interface{}{"service_id": 999999})
		invalidServiceReq := authRequest(http.MethodPut, fmt.Sprintf("/conversations/%d", createdConversation.Id), invalidServiceBody, clientToken)
		invalidServiceResp := httptest.NewRecorder()
		router.ServeHTTP(invalidServiceResp, invalidServiceReq)
		assert.Equal(t, http.StatusNotFound, invalidServiceResp.Code)

		invalidProposalBody, _ := json.Marshal(map[string]interface{}{"proposal_id": 999999})
		invalidProposalReq := authRequest(http.MethodPatch, fmt.Sprintf("/conversations/%d", createdConversation.Id), invalidProposalBody, clientToken)
		invalidProposalResp := httptest.NewRecorder()
		router.ServeHTTP(invalidProposalResp, invalidProposalReq)
		assert.Equal(t, http.StatusNotFound, invalidProposalResp.Code)
	})
}

func TestMessageAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, TestPool)

	analystUser := pkg.UserRegister{
		Name:          "Message Analyst",
		Email:         "messageanalyst@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999997",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	analystToken, analystUserID := registerAndLogin(t, router, analystUser)

	clientUser := pkg.UserRegister{
		Name:          "Message Client",
		Email:         "messageclient@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999996",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	clientToken, clientUserID := registerAndLogin(t, router, clientUser)

	thirdUser := pkg.UserRegister{
		Name:          "Message Third",
		Email:         "messagethird@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "41999999995",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}
	thirdToken, _ := registerAndLogin(t, router, thirdUser)

	analystBody, _ := json.Marshal(pkg.AnalystCreateRequest{Hourly_rate: 100.0})
	analystReq := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/analyst", analystUserID), analystBody, analystToken)
	analystResp := httptest.NewRecorder()
	router.ServeHTTP(analystResp, analystReq)
	assert.Equal(t, http.StatusCreated, analystResp.Code)

	clientBody, _ := json.Marshal(pkg.ClientCreateRequest{Proposed_budget: 120.0})
	clientReq := authRequest(http.MethodPost, fmt.Sprintf("/users/%d/client", clientUserID), clientBody, clientToken)
	clientResp := httptest.NewRecorder()
	router.ServeHTTP(clientResp, clientReq)
	assert.Equal(t, http.StatusConflict, clientResp.Code)

	proposalBody, _ := json.Marshal(pkg.ProposalLetterCreateRequest{
		Title:                "Mensagem de teste",
		Content:              "Conteúdo",
		Analyst_id:           analystUserID,
		Client_id:            clientUserID,
		Proposed_hourly_rate: 80.0,
	})
	proposalReq := authRequest(http.MethodPost, "/proposals", proposalBody, clientToken)
	proposalResp := httptest.NewRecorder()
	router.ServeHTTP(proposalResp, proposalReq)
	assert.Equal(t, http.StatusCreated, proposalResp.Code)
	var proposalPayload pkg.ProposalLetterResponse
	json.Unmarshal(proposalResp.Body.Bytes(), &proposalPayload)

	serviceBody, _ := json.Marshal(pkg.ServiceCreateRequest{
		Title:              "Serviço de mensagem",
		Content:            "Conteúdo",
		Hourly_rate:        85.0,
		Status:             "RUNNING",
		Proposal_letter_id: proposalPayload.Proposal_letter.Id,
	})
	serviceReq := authRequest(http.MethodPost, "/services", serviceBody, clientToken)
	serviceResp := httptest.NewRecorder()
	router.ServeHTTP(serviceResp, serviceReq)
	assert.Equal(t, http.StatusCreated, serviceResp.Code)
	var servicePayload pkg.ServiceResponse
	json.Unmarshal(serviceResp.Body.Bytes(), &servicePayload)

	conversationBody, _ := json.Marshal(map[string]interface{}{
		"analyst_id":  analystUserID,
		"client_id":   clientUserID,
		"service_id":  servicePayload.Service.Id,
		"proposal_id": proposalPayload.Proposal_letter.Id,
	})
	conversationReq := authRequest(http.MethodPost, "/conversations", conversationBody, clientToken)
	conversationResp := httptest.NewRecorder()
	router.ServeHTTP(conversationResp, conversationReq)
	assert.Equal(t, http.StatusCreated, conversationResp.Code)
	var conversation pkg.Conversation
	json.Unmarshal(conversationResp.Body.Bytes(), &conversation)

	t.Run("only participants can send messages", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{"sender_id": 999999, "content": "hello"})
		messageReq := authRequest(http.MethodPost, fmt.Sprintf("/conversations/%d/messages", conversation.Id), messageBody, thirdToken)
		messageResp := httptest.NewRecorder()
		router.ServeHTTP(messageResp, messageReq)
		assert.Equal(t, http.StatusForbidden, messageResp.Code)
	})

	t.Run("only the sender can edit or delete a message", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{"sender_id": clientUserID, "content": "hello world"})
		messageReq := authRequest(http.MethodPost, fmt.Sprintf("/conversations/%d/messages", conversation.Id), messageBody, clientToken)
		messageResp := httptest.NewRecorder()
		router.ServeHTTP(messageResp, messageReq)
		assert.Equal(t, http.StatusCreated, messageResp.Code)
		var message pkg.Message
		json.Unmarshal(messageResp.Body.Bytes(), &message)

		updateBody, _ := json.Marshal(map[string]interface{}{"content": "tampered"})
		updateReq := authRequest(http.MethodPut, fmt.Sprintf("/conversations/%d/messages/%d", conversation.Id, message.Id), updateBody, analystToken)
		updateResp := httptest.NewRecorder()
		router.ServeHTTP(updateResp, updateReq)
		assert.Equal(t, http.StatusForbidden, updateResp.Code)

		deleteReq := authRequest(http.MethodDelete, fmt.Sprintf("/conversations/%d/messages/%d", conversation.Id, message.Id), nil, analystToken)
		deleteResp := httptest.NewRecorder()
		router.ServeHTTP(deleteResp, deleteReq)
		assert.Equal(t, http.StatusForbidden, deleteResp.Code)
	})
}
