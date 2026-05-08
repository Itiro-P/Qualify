package test

import (
	"bytes"
	"encoding/json"
	"main/pkg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func registerAndLogin(t *testing.T, router *gin.Engine, user pkg.UserRegister) (string, int) {
	t.Helper()

	// Registra o usuário
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code, "falhou ao registrar usuário")

	var regResp pkg.UserResponse // Ajuste para o seu struct de resposta de registro
	json.Unmarshal(w.Body.Bytes(), &regResp)
	userID := regResp.User.Id

	// Faz login e pega o token
	loginBody, _ := json.Marshal(pkg.UserLogin{
		Email:    user.Email,
		Password: user.Password,
	})
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, "falhou ao fazer login")

	var loginResp pkg.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NotEmpty(t, loginResp.AccessToken, "token não pode ser vazio")

	return loginResp.AccessToken, userID
}

// authRequest cria uma request já com o header de autorização
func authRequest(method, url string, body []byte, token string) *http.Request {
	var req *http.Request

	if body != nil {
		req = httptest.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return req
}
