package test

import (
	"bytes"
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

var clients = []pkg.Client{
	{
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
	},
	{
		User: pkg.User{
			Name:          "Frank",
			Email:         "frank@utfpr.edu.br",
			Phone:         "41969696969",
			Country_code:  "RO",
			Country_name:  "Romania",
			Country_state: "AA",
			City:          "Bucareste",
			Timezone:      "America/Sao_Paulo",
		},
		Proposed_budget: 69.0,
	},
}

func TestClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

	postClientResponse := []pkg.Client{}

	// Primeiros testes para criação de clientes, que dependem da criação prévia de usuários
	for _, c := range clients {
		t.Run("Criar Cliente para "+c.User.Name, func(t *testing.T) {
			client := c
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
				t.Error("O ID do cliente não deveria ser zero")
			}

			// Então adicionamos o cliente criado à resposta para verificarmos depois se todos foram criados corretamente

			postClientResponse = append(postClientResponse, clientResponse.Client)
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

		if len(clientsResponse.Clients) != len(postClientResponse) {
			t.Errorf("Número de clientes retornados (%d) diferente do número de clientes criados (%d)", len(clientsResponse.Clients), len(postClientResponse))
		}
		assert.ElementsMatch(t, clientsResponse.Clients, postClientResponse)
	})

	// Agora pegaremos cada cliente em específico para ver se eles estão sendo retornados corretamente
	for _, a := range postClientResponse {
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
		frank := slices.IndexFunc(postClientResponse, func(a pkg.Client) bool {
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
		assert.Equal(t, clientsResponse.Clients[0], postClientResponse[frank])
	})

	t.Run("Listando clientes com filtro de país", func(t *testing.T) {
		frank := slices.IndexFunc(postClientResponse, func(a pkg.Client) bool {
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
		assert.Equal(t, clientsResponse.Clients[0], postClientResponse[frank])
	})

	t.Run("Listando clientes com filtro de cidade", func(t *testing.T) {
		client := slices.IndexFunc(postClientResponse, func(a pkg.Client) bool {
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
		assert.Equal(t, clientsResponse.Clients[0], postClientResponse[client])
	})

	t.Run("Listando clientes com filtro de maior valor", func(t *testing.T) {
		maxim := slices.IndexFunc(postClientResponse, func(a pkg.Client) bool {
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
		assert.Equal(t, clientsResponse.Clients[0], postClientResponse[maxim])
	})

	t.Run("Listando clientes com filtro de menor valor", func(t *testing.T) {
		maxim := slices.IndexFunc(postClientResponse, func(a pkg.Client) bool {
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
		assert.Equal(t, clientsResponse.Clients[0], postClientResponse[maxim])
	})

	/**
	 * 'Ah mas e o PUT?' Foda-se o PUT, QUEM RAIOS VAI QUERER ATUALIZAR UM CLIENTE INTEIRO?
	 */

	// Agora modificaremos um cliente para ver se a atualização está funcionando corretamente
	t.Run("Atualizando cliente "+postClientResponse[0].User.Name, func(t *testing.T) {
		client := postClientResponse[0]
		patchBody := map[string]interface{}{
			"proposed_budget": 150.0,
		}
		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/users/%d/client", client.User.Id)

		req := httptest.NewRequest(http.MethodPatch, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

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
	for _, c := range postClientResponse {
		t.Run("Deletando cliente "+c.User.Name+" por ID", func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d/client", c.User.Id)

			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}

	// Deletaremos também os usuários associados aos clientes para limpar o banco de dados
	for _, a := range postClientResponse {
		t.Run("Deletando usuário associado ao cliente "+a.User.Name, func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d", a.User.Id)

			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
