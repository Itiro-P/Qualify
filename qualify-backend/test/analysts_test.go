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
			Country_code:  "BR",
			Country_name:  "Brazil",
			Country_state: "PR",
			City:          "Campo Mourão",
			Timezone:      "America/Sao_Paulo",
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
			City:          "Campo Mourão",
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
			Country_state: "PR",
			City:          "Campo Mourão",
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
		})
	}
}
