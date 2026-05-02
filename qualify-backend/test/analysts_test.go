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

var analysts = []pkg.Analyst{
	{
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
	},
	{
		User: pkg.User{
			Name:          "John Xina",
			Email:         "xina@utfpr.edu.br",
			Phone:         "41969696969",
			Country_code:  "CN",
			Country_name:  "China",
			Country_state: "Beijing",
			City:          "Beijing",
			Timezone:      "Asia/Shanghai",
		},
		Hourly_rate:   69.0,
		Total_reviews: 5,
		Mean_rating:   ToPtrFloat64(4.2),
	},
	{
		User: pkg.User{
			Name:          "Ivanilton Pelado",
			Email:         "nudismo@utfpr.edu.br",
			Phone:         "44999999999",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		Hourly_rate:   67.0,
		Total_reviews: 60,
		Mean_rating:   ToPtrFloat64(3.9),
	},
	{
		User: pkg.User{
			Name:          "João Paumolence",
			Email:         "joaum@utfpr.edu.br",
			Phone:         "44999990000",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Roncador",
			Timezone:      "America/Sao_Paulo",
		},
		Hourly_rate:   57.0,
		Total_reviews: 30,
		Mean_rating:   ToPtrFloat64(4.0),
	},
	{
		User: pkg.User{
			Name:          "Alex do Durex",
			Email:         "alex@utfpr.edu.br",
			Phone:         "44999690000",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
		},
		Hourly_rate:   50.0,
		Total_reviews: 25,
		Mean_rating:   ToPtrFloat64(4.0),
	},
	{
		User: pkg.User{
			Name:          "Rodrigo do Piolho",
			Email:         "rodrigo@utfpr.edu.br",
			Phone:         "44999690000",
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "GO",
			City:          "Goiânia",
			Timezone:      "America/Sao_Paulo",
		},
		Hourly_rate:   128.0,
		Total_reviews: 28,
		Mean_rating:   ToPtrFloat64(1.6),
	},
}

func TestAnalyst(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

	postAnalystResponse := []pkg.Analyst{}

	// Primeiros testes para criação de analistas, que dependem da criação prévia de usuários
	for _, a := range analysts {
		t.Run("Criar Analista para "+a.User.Name, func(t *testing.T) {
			analyst := a
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

			if analystResponse.Analyst.Id == 0 {
				t.Error("O ID do analista não deveria ser zero")
			}

			// Então adicionamos o analista criado à resposta para verificarmos depois se todos foram criados corretamente

			postAnalystResponse = append(postAnalystResponse, analystResponse.Analyst)
		})
	}

	// Agora vemos se todos os analistas foram criados corretamente
	t.Run("Listando todos os analistas", func(t *testing.T) {
		targetURL := "/analysts"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != len(postAnalystResponse) {
			t.Errorf("Número de analistas retornados (%d) diferente do número de analistas criados (%d)", len(analystsResponse.Analysts), len(postAnalystResponse))
		}
		assert.ElementsMatch(t, analystsResponse.Analysts, postAnalystResponse)
	})

	// Agora pegaremos cada analista em específico para ver se eles estão sendo retornados corretamente
	for _, a := range postAnalystResponse {
		t.Run("Pegando analista "+a.User.Name+" por ID", func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d/analyst", a.User.Id)

			req := httptest.NewRequest(http.MethodGet, targetURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var analystResponse pkg.AnalystResponse
			json.Unmarshal(w.Body.Bytes(), &analystResponse)
			assert.Equal(t, analystResponse.Analyst, a)
		})
	}

	// Agora testaremos os filtros de listagem de analistas, para ver se eles estão funcionando corretamente
	t.Run("Listando analistas com filtro de nome", func(t *testing.T) {
		johnXina := slices.IndexFunc(postAnalystResponse, func(a pkg.Analyst) bool {
			return strings.HasPrefix(a.User.Name, "John Xina")
		})

		targetURL := "/analysts?name=John"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != 1 {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (1)", len(analystsResponse.Analysts))
		}
		assert.Equal(t, analystsResponse.Analysts[0], postAnalystResponse[johnXina])
	})

	t.Run("Listando analistas com filtro de país", func(t *testing.T) {
		johnXina := slices.IndexFunc(postAnalystResponse, func(a pkg.Analyst) bool {
			return strings.HasPrefix(a.User.Country_name, "China")
		})

		targetURL := "/analysts?country=China"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != 1 {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (1)", len(analystsResponse.Analysts))
		}
		assert.Equal(t, analystsResponse.Analysts[0], postAnalystResponse[johnXina])
	})

	t.Run("Listando analistas com filtro de cidade", func(t *testing.T) {
		joao := slices.IndexFunc(postAnalystResponse, func(a pkg.Analyst) bool {
			return strings.HasPrefix(a.User.City, "Roncador")
		})

		targetURL := "/analysts?city=Roncador"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != 1 {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (1)", len(analystsResponse.Analysts))
		}
		assert.Equal(t, analystsResponse.Analysts[0], postAnalystResponse[joao])
	})

	t.Run("Listando analistas com filtro de maior valor", func(t *testing.T) {
		maxim := slices.IndexFunc(postAnalystResponse, func(a pkg.Analyst) bool {
			return a.Hourly_rate >= 128.0
		})

		targetURL := "/analysts?min_hourly_rate=128"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != 1 {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (1)", len(analystsResponse.Analysts))
		}
		assert.Equal(t, analystsResponse.Analysts[0], postAnalystResponse[maxim])
	})

	t.Run("Listando analistas com filtro de menor valor", func(t *testing.T) {
		maxim := slices.IndexFunc(postAnalystResponse, func(a pkg.Analyst) bool {
			return a.Hourly_rate <= 50.0
		})

		targetURL := "/analysts?max_hourly_rate=50"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != 1 {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (1)", len(analystsResponse.Analysts))
		}
		assert.Equal(t, analystsResponse.Analysts[0], postAnalystResponse[maxim])
	})

	t.Run("Listando analistas com filtro de mínimo de avaliações", func(t *testing.T) {
		expected := []pkg.Analyst{}
		for _, a := range postAnalystResponse {
			if a.Total_reviews >= 60 {
				expected = append(expected, a)
			}
		}

		targetURL := "/analysts?min_total_reviews=60"
		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != len(expected) {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (%d)", len(analystsResponse.Analysts), len(expected))
		}
		if len(expected) > 0 {
			assert.Equal(t, analystsResponse.Analysts[0], expected[0])
		}
	})

	t.Run("Listando analistas com filtro de mínimo de média de avaliações", func(t *testing.T) {
		expected := []pkg.Analyst{}
		for _, a := range postAnalystResponse {
			if a.Mean_rating != nil && *a.Mean_rating >= 4.5 {
				expected = append(expected, a)
			}
		}

		targetURL := "/analysts?min_mean_rating=4.5"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != len(expected) {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (%d)", len(analystsResponse.Analysts), len(expected))
		}
		if len(expected) > 0 {
			assert.ElementsMatch(t, analystsResponse.Analysts, expected)
		}
	})

	t.Run("Listando analistas ordenando do menor para o maior", func(t *testing.T) {
		arr := postAnalystResponse
		slices.SortFunc(arr, func(a, b pkg.Analyst) int {
			// strings.Compare retorna -1 se a < b, 0 se a == b, 1 se a > b
			return strings.Compare(a.User.Name, b.User.Name)
		})

		targetURL := "/analysts?sort_by=name&order=ASC"

		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystsResponse pkg.AnalystsResponse
		json.Unmarshal(w.Body.Bytes(), &analystsResponse)

		if len(analystsResponse.Analysts) != len(arr) {
			t.Errorf("Número de analistas retornados (%d) diferente do esperado (%d)", len(analystsResponse.Analysts), len(arr))
		}
		assert.Equal(t, analystsResponse.Analysts, arr)
	})

	/**
	 * 'Ah mas e o PUT?' Foda-se o PUT, QUEM RAIOS VAI QUERER ATUALIZAR UM ANALISTA INTEIRO?
	 */

	// Agora modificaremos um analista para ver se a atualização está funcionando corretamente
	t.Run("Atualizando analista "+postAnalystResponse[0].User.Name, func(t *testing.T) {
		analyst := postAnalystResponse[0]
		patchBody := map[string]interface{}{
			"hourly_rate": 150.0,
		}
		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/users/%d/analyst", analyst.User.Id)

		req := httptest.NewRequest(http.MethodPatch, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var analystResponse pkg.AnalystResponse
		json.Unmarshal(w.Body.Bytes(), &analystResponse)

		if analystResponse.Analyst.Hourly_rate != 150.0 {
			t.Errorf("Valor do campo 'hourly_rate' diferente do esperado (150.0), obtemos %f", analystResponse.Analyst.Hourly_rate)
		}
	})

	// Agora deletaremos cada analista em específico para ver se eles estão sendo deletados corretamente
	for _, a := range postAnalystResponse {
		t.Run("Deletando analista "+a.User.Name+" por ID", func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d/analyst", a.User.Id)

			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}

	// Deletaremos também os usuários associados aos analistas para limpar o banco de dados
	for _, a := range postAnalystResponse {
		t.Run("Deletando usuário associado ao analista "+a.User.Name, func(t *testing.T) {
			targetURL := fmt.Sprintf("/users/%d", a.User.Id)

			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
