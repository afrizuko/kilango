package auth

import (
	"bytes"
	"encoding/json"
	"github.com/afrizuko/kilango/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_AuthenticateUser(t *testing.T) {
	handler := NewHandler(NewStub())
	t.Run("it returns a valid JWT token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", newAuthData(t, "1", "001"))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var token model.AuthResponse
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&token))
		assert.Equal(t, 3, len(strings.Split(token.Token, ".")))
	})

	t.Run("it returns unauthorized for invalid credentials", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", newAuthData(t, "1", "002"))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
	})
}

func Test_GetAuthenticatedUser(t *testing.T) {
	handler := NewHandler(NewStub())
	req := httptest.NewRequest("POST", "/", newAuthData(t, "1", "001"))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)

	var token model.AuthResponse
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&token))
	assert.Equal(t, 3, len(strings.Split(token.Token, ".")))

	t.Run("it returns an authenticated user's profile", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}

func newAuthData(t *testing.T, username, password string) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	assert.NoError(t, json.NewEncoder(buffer).Encode(model.AuthRequest{
		Username: username,
		Password: password,
	}))
	return buffer
}
