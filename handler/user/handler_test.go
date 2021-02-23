package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/audit/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var handler *Handler

func init() {
	handler = NewHandler(NewStub())
}

func Test_GetUsers(t *testing.T) {

	t.Run("it returns a list of users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()

		query := req.URL.Query()
		query.Add("page", "0")
		query.Add("limit", "2")
		req.URL.RawQuery = query.Encode()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var users []model.User
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&users))
		assert.Equal(t, 2, len(users))
	})

	for i := 1; i <= 2; i++ {
		t.Run(fmt.Sprintf("it returns a single user %d", i), func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/%d", i), nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Result().StatusCode)

			var user model.User
			assert.NoError(t, json.NewDecoder(res.Body).Decode(&user))
			assert.Equal(t, uint(i), user.ID)
		})
	}
}

func Test_GetUser(t *testing.T) {

	for i := 1; i <= 2; i++ {
		t.Run(fmt.Sprintf("it returns a single user %d", i), func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/%d", i), nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Result().StatusCode)

			var user model.User
			assert.NoError(t, json.NewDecoder(res.Body).Decode(&user))
			assert.Equal(t, uint(i), user.ID)
		})
	}

	t.Run("it returns a 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/0", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)

		errorVal := make(map[string]string)
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&errorVal))
		assert.Contains(t, errorVal, "error")
	})
}

func Test_CreateUser(t *testing.T) {
	req := httptest.NewRequest("POST", "/", newUser(t, "peter"))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	assert.Equal(t, http.StatusCreated, res.Result().StatusCode)

	var user model.User
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&user))
	assert.Equal(t, "peter", user.Username)
}

func Test_ModifyUser(t *testing.T) {
	req := httptest.NewRequest("PUT", "/1", newUser(t, "joshua"))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	assert.Equal(t, http.StatusAccepted, res.Result().StatusCode)

	var user model.User
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&user))
	assert.Equal(t, "joshua", user.Username)
}

func Test_DeleteUser(t *testing.T) {
	t.Run("it returns an OK after deleting a user", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/1", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	})

	t.Run("it returns a 404 for user not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/1", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, 404, res.Result().StatusCode)
	})
}

func newUser(t *testing.T, s string) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	assert.NoError(t, json.NewEncoder(buffer).Encode(model.User{
		Username: s,
		Password: s,
	}))
	return buffer
}
