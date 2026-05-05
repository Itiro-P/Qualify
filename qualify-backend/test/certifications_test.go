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

var certifications = []pkg.Certification{
	{
		Name:        "Mestrado em Teste Baseado em Aspectos",
		Description: "Que aspectos?",
		Year:        1900,
	},
	{
		Name:        "Certificado de Sobrevivência em Reuniões Infinitas",
		Description: "Participou e saiu vivo de uma maratona de reuniões.",
		Year:        2024,
	},
	{
		Name:        "Especialista em Copiar/Colar Estratégico",
		Description: "Domina o Ctrl+C/Ctrl+V com elegância.",
		Year:        2023,
	},
	{
		Name:        "Mestre em Café e Debug",
		Description: "Resolve bugs consistentemente após a terceira xícara.",
		Year:        2025,
	},
	{
		Name:        "Ninja do Merge sem Conflito",
		Description: "Resolve conflitos de git com silêncio e honra.",
		Year:        2022,
	},
}

func TestCertifications(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

	postCertResponse := []pkg.Certification{}

	// Primeiros testes para criação de certificações
	for _, c := range certifications {
		t.Run("Criar Certificação para "+c.Name, func(t *testing.T) {
			certification := c
			body, _ := json.Marshal(certification)
			targetURL := "/certifications"

			req := httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificamos se a certificação foi criada com sucesso
			assert.Equal(t, http.StatusCreated, w.Code)
			var certificationResponse pkg.CertificationResponse
			json.Unmarshal(w.Body.Bytes(), &certificationResponse)

			if certificationResponse.Certification.Id == 0 {
				t.Error("O ID da certificação não deveria ser zero")
			}

			postCertResponse = append(postCertResponse, certificationResponse.Certification)
		})
	}

	// Agora veremos se as listagens estão funcionando
	t.Run("Listando todos os analistas", func(t *testing.T) {
		targetURL := "/certifications"
		req := httptest.NewRequest(http.MethodGet, targetURL, nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificamos se a certificação foi criada com sucesso
		assert.Equal(t, http.StatusOK, w.Code)

		var certificationResponse pkg.CertificationsResponse
		json.Unmarshal(w.Body.Bytes(), &certificationResponse)
		assert.ElementsMatch(t, postCertResponse, certificationResponse.Certifications)
	})

	// Agora veremos o GET por ID
	for _, c := range postCertResponse {
		t.Run("Listando Certificação para "+c.Name, func(t *testing.T) {
			targetURL := fmt.Sprintf("/certifications/%d", c.Id)

			req := httptest.NewRequest(http.MethodGet, targetURL, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var certificationReponse pkg.CertificationResponse
			json.Unmarshal(w.Body.Bytes(), &certificationReponse)
			assert.Equal(t, c, certificationReponse.Certification)
		})
	}

	// Agora atualizaremos uma certificação
	t.Run("Atualizando uma certificação ", func(t *testing.T) {
		var cert = postCertResponse[0]
		patchBody := map[string]interface{}{
			"year": 1920,
		}

		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/certifications/%d", cert.Id)

		req := httptest.NewRequest(http.MethodPatch, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var certResponse pkg.CertificationResponse
		json.Unmarshal(w.Body.Bytes(), &certResponse)

		assert.Equal(t, certResponse.Certification.Year, 1920)
	})

	// Agora testaremos as certificações no analista
	t.Run("Apontando certificação para um analista", func(t *testing.T) {
		analyst := pkg.Analyst{
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
		}

		t.Run("Criar Analista para "+analyst.User.Name, func(t *testing.T) {
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

			analyst = analystResponse.Analyst
		})

		certification := postCertResponse[0]

		targetURL := fmt.Sprintf("/users/%d/analyst/certifications", analyst.Id)

		reqBody := pkg.AnalystCertification{
			Certification_id: certification.Id,
			Analyst_id:       analyst.Id,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Agora faremos o cleanup do analista criado
		targetURL = fmt.Sprintf("/users/%d/analyst", analyst.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		targetURL = fmt.Sprintf("/users/%d", analyst.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})

	// Agora deletaremos as certificações criadas
	for _, c := range postCertResponse {
		t.Run("Removendo Certificação para "+c.Name, func(t *testing.T) {
			certification := c
			targetURL := fmt.Sprintf("/certifications/%d", certification.Id)

			req := httptest.NewRequest(http.MethodDelete, targetURL, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificamos se a certificação foi criada com sucesso
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
