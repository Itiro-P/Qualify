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
