package test

import (
	"encoding/json"
	"fmt"
	"main/internal/routes"
	"main/pkg"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var userRegisters = []pkg.UserRegister{
		{
			Name:          "Marcos Calvaro",
			Email:         "markos@utfpr.edu.br",
			Password:      "aabbccddee",
			Phone:         "41999999989",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		{
			Name:          "Frank do Tank",
			Email:         "frank@utfpr.edu.br",
			Password:      "aabbccddee",
			Phone:         "41969696969",
			Country_code:  "RO",
			Country_name:  "Romania",
			Country_state: "AA",
			City:          "Bucareste",
			Timezone:      "America/Sao_Paulo",
		},
	}
	var clients = []pkg.Client{
		{
			User:            pkg.User{},
			Proposed_budget: 100.0,
		},
		{
			User:            pkg.User{},
			Proposed_budget: 69.0,
		},
	}
	tokens := []string{}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

	// Primeiros testes para criação de clientes, que dependem da criação prévia de usuários
	for i := range len(userRegisters) {
		t.Run("Criar Cliente para "+userRegisters[i].Name, func(t *testing.T) {
			// Agora pegamos o token e o ID de uma vez
			token, userID := registerAndLogin(t, router, userRegisters[i])

			// Atribui o ID ao analista antes de enviar o POST
			clients[i].User.Id = userID

			body, _ := json.Marshal(clients[i])
			targetURL := fmt.Sprintf("/users/%d/client", userID)

			req := authRequest(http.MethodPost, targetURL, body, token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var clientResponse pkg.ClientResponse
			json.Unmarshal(w.Body.Bytes(), &clientResponse)
			assert.NotZero(t, clientResponse.Client.Id)

			clients[i] = clientResponse.Client
			tokens = append(tokens, token)
		})
	}

	// Agora vemos se todos os clientes foram criados corretamente
	t.Run("Listando todos os clientes", func(t *testing.T) {
		targetURL := "/clients"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientsResponse pkg.ClientsResponse
		json.Unmarshal(w.Body.Bytes(), &clientsResponse)

		if len(clientsResponse.Clients) != len(clients) {
			t.Errorf("Número de clientes retornados (%d) diferente do número de clientes criados (%d)", len(clientsResponse.Clients), len(clients))
		}
		assert.ElementsMatch(t, clientsResponse.Clients, clients)
	})

	// Agora pegaremos cada cliente em específico para ver se eles estão sendo retornados corretamente
	for _, a := range clients {
		t.Run("Pegando cliente "+a.User.Name+" por ID", func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d/client", a.User.Id)

			req := httptest.NewRequest(http.MethodGet, targetURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var clientResponse pkg.ClientResponse
			json.Unmarshal(w.Body.Bytes(), &clientResponse)
			assert.Equal(t, clientResponse.Client, a)
		})
	}

	// Agora testaremos os filtros de listagem de clientes, para ver se eles estão funcionando corretamente
	t.Run("Listando clientes com filtro de nome", func(t *testing.T) {
		frank := slices.IndexFunc(clients, func(a pkg.Client) bool {
			return strings.HasPrefix(a.User.Name, "Frank")
		})

		targetURL := "/clients?name=Frank"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientsResponse pkg.ClientsResponse
		json.Unmarshal(w.Body.Bytes(), &clientsResponse)

		if len(clientsResponse.Clients) != 1 {
			t.Errorf("Número de clientes retornados (%d) diferente do esperado (1)", len(clientsResponse.Clients))
		}
		assert.Equal(t, clientsResponse.Clients[0], clients[frank])
	})

	t.Run("Listando clientes com filtro de país", func(t *testing.T) {
		frank := slices.IndexFunc(clients, func(a pkg.Client) bool {
			return strings.HasPrefix(a.User.Country_name, "Romania")
		})

		targetURL := "/clients?country=Romania"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientsResponse pkg.ClientsResponse
		json.Unmarshal(w.Body.Bytes(), &clientsResponse)

		if len(clientsResponse.Clients) != 1 {
			t.Errorf("Número de clientes retornados (%d) diferente do esperado (1)", len(clientsResponse.Clients))
		}
		assert.Equal(t, clientsResponse.Clients[0], clients[frank])
	})

	t.Run("Listando clientes com filtro de cidade", func(t *testing.T) {
		client := slices.IndexFunc(clients, func(a pkg.Client) bool {
			return strings.HasPrefix(a.User.City, "Campo Mourão")
		})

		targetURL := "/clients?city=Campo"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientsResponse pkg.ClientsResponse
		json.Unmarshal(w.Body.Bytes(), &clientsResponse)

		if len(clientsResponse.Clients) != 1 {
			t.Errorf("Número de clientes retornados (%d) diferente do esperado (1)", len(clientsResponse.Clients))
		}
		assert.Equal(t, clientsResponse.Clients[0], clients[client])
	})

	t.Run("Listando clientes com filtro de maior valor", func(t *testing.T) {
		maxim := slices.IndexFunc(clients, func(a pkg.Client) bool {
			return a.Proposed_budget >= 100.0
		})

		targetURL := "/clients?min_proposed_budget=100"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientsResponse pkg.ClientsResponse
		json.Unmarshal(w.Body.Bytes(), &clientsResponse)

		if len(clientsResponse.Clients) != 1 {
			t.Errorf("Número de clientes retornados (%d) diferente do esperado (1)", len(clientsResponse.Clients))
		}
		assert.Equal(t, clientsResponse.Clients[0], clients[maxim])
	})

	t.Run("Listando clientes com filtro de menor valor", func(t *testing.T) {
		maxim := slices.IndexFunc(clients, func(a pkg.Client) bool {
			return a.Proposed_budget <= 69.0
		})

		targetURL := "/clients?max_proposed_budget=69"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientsResponse pkg.ClientsResponse
		json.Unmarshal(w.Body.Bytes(), &clientsResponse)

		if len(clientsResponse.Clients) != 1 {
			t.Errorf("Número de clientes retornados (%d) diferente do esperado (1)", len(clientsResponse.Clients))
		}
		assert.Equal(t, clientsResponse.Clients[0], clients[maxim])
	})

	/**
	 * 'Ah mas e o PUT?' Foda-se o PUT, QUEM RAIOS VAI QUERER ATUALIZAR UM CLIENTE INTEIRO?
	 */

	// Agora modificaremos um cliente para ver se a atualização está funcionando corretamente
	t.Run("Atualizando cliente "+clients[0].User.Name, func(t *testing.T) {
		client := clients[0]
		patchBody := map[string]interface{}{
			"proposed_budget": 150.0,
		}
		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/users/%d/client", client.User.Id)

		req := authRequest(http.MethodPatch, targetURL, body, tokens[0])

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var clientResponse pkg.ClientResponse
		json.Unmarshal(w.Body.Bytes(), &clientResponse)

		if clientResponse.Client.Proposed_budget != 150.0 {
			t.Errorf("Valor do campo 'proposed_budget' diferente do esperado (150.0), obtemos %f", clientResponse.Client.Proposed_budget)
		}
	})

	// Agora deletaremos cada cliente em específico para ver se eles estão sendo deletados corretamente
	for i, c := range clients {
		t.Run("Deletando cliente "+c.User.Name+" por ID", func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d/client", c.User.Id)
			req := authRequest(http.MethodDelete, targetURL, nil, tokens[i])
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})
	}

	// Deletaremos também os usuários associados aos clientes para limpar o banco de dados
	for i, u := range clients {
		t.Run("Deletando usuário associado ao cliente "+u.User.Name, func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d", u.User.Id)
			req := authRequest(http.MethodDelete, targetURL, nil, tokens[i])
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})
	}
}
