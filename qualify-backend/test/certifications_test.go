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

func TestCertifications(t *testing.T) {
	var certifications = []pkg.Certification{
		{
			Name:        "Mestrado em Teste Baseado em Aspectos",
			Institution: "UTFPR",
			Description: "Que aspectos?",
			Year:        1900,
		},
		{
			Name:        "Certificado de Sobrevivência em Reuniões Infinitas",
			Institution: "UTFPR",
			Description: "Participou e saiu vivo de uma maratona de reuniões.",
			Year:        2024,
		},
		{
			Name:        "Especialista em Copiar/Colar Estratégico",
			Institution: "UTFPR",
			Description: "Domina o Ctrl+C/Ctrl+V com elegância.",
			Year:        2023,
		},
		{
			Name:        "Mestre em Café e Debug",
			Institution: "UTFPR",
			Description: "Resolve bugs consistentemente após a terceira xícara.",
			Year:        2025,
		},
		{
			Name:        "Técnico em Computaria",
			Institution: "UTFPR",
			Description: "Possui habilidades extraordinárias em computaria.",
			Year:        1960,
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	routes.SetupRoutes(router, TestPool)

	userRegister := pkg.UserRegister{
		Name:          "Paulo Sabo",
		Email:         "sabo@utfpr.edu.br",
		Password:      "aabbccddee",
		Phone:         "45999697000",
		Country_code:  "BR",
		Country_name:  "Brazil",
		Country_state: "PR",
		City:          "Campo Mourão",
		Timezone:      "America/Sao_Paulo",
	}

	analyst := pkg.Analyst{
		Hourly_rate:   128.0,
		Total_reviews: 28,
		Mean_rating:   ToPtrFloat64(1.6),
	}

	token, userID := registerAndLogin(t, router, userRegister)

	// Primeiros testes para criação de certificações
	for i, c := range certifications {
		t.Run("Criar Certificação para "+c.Name, func(t *testing.T) {
			certification := c
			body, _ := json.Marshal(certification)
			targetURL := "/certifications"

			req := authRequest(http.MethodPost, targetURL, body, token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificamos se a certificação foi criada com sucesso
			assert.Equal(t, http.StatusCreated, w.Code)
			var certificationResponse pkg.CertificationResponse
			json.Unmarshal(w.Body.Bytes(), &certificationResponse)

			if certificationResponse.Certification.Id == 0 {
				t.Error("O ID da certificação não deveria ser zero")
			}
			certifications[i] = certificationResponse.Certification
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
		assert.ElementsMatch(t, certifications, certificationResponse.Certifications)
	})

	// Agora veremos o GET por ID
	for _, c := range certifications {
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
		var cert = certifications[0]
		patchBody := map[string]interface{}{
			"year": 1920,
		}

		body, _ := json.Marshal(patchBody)
		targetURL := fmt.Sprintf("/certifications/%d", cert.Id)

		req := authRequest(http.MethodPatch, targetURL, body, token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var certResponse pkg.CertificationResponse
		json.Unmarshal(w.Body.Bytes(), &certResponse)

		assert.Equal(t, certResponse.Certification.Year, 1920)
	})

	// Agora testaremos as certificações no analista
	t.Run("Apontando certificação para um analista", func(t *testing.T) {
		body, _ := json.Marshal(analyst)
		targetURL := fmt.Sprintf("/users/%d/analyst", userID)

		req := authRequest(http.MethodPost, targetURL, body, token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var analystResponse pkg.AnalystResponse
		json.Unmarshal(w.Body.Bytes(), &analystResponse)
		assert.NotZero(t, analystResponse.Analyst.Id)

		analyst = analystResponse.Analyst

		certification := certifications[0]

		targetURL = fmt.Sprintf("/users/%d/analyst/certifications", analyst.Id)

		reqBody := pkg.AnalystCertification{
			Certification_id: certification.Id,
			Analyst_id:       analyst.Id,
		}
		body, _ = json.Marshal(reqBody)
		req = authRequest(http.MethodPost, targetURL, body, token)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Agora faremos o cleanup do analista criado
		targetURL = fmt.Sprintf("/users/%d/analyst", analyst.User.Id)

		req = httptest.NewRequest(http.MethodDelete, targetURL, nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})

	// Agora deletaremos as certificações criadas
	for _, c := range certifications {
		t.Run("Removendo Certificação para "+c.Name, func(t *testing.T) {
			targetURL := fmt.Sprintf("/certifications/%d", c.Id)
			req := authRequest(http.MethodDelete, targetURL, nil, token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificamos se a certificação foi removida com sucesso
			assert.Equal(t, http.StatusNoContent, w.Code)
		})
	}
}
