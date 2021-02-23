package user

import (
	"encoding/json"
	"fmt"
	"github.com/audit/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateGetUpdateDeleteUsers(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	handler = NewHandler(model.NewUserStub())

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, httptest.NewRequest("POST", "/", newUser(t, "peter")))
	assert.Equal(t, http.StatusCreated, res.Result().StatusCode)

	var user1 model.User
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&user1))
	assert.Equal(t, "peter", user1.Username)

	res2 := httptest.NewRecorder()
	handler.ServeHTTP(res2, httptest.NewRequest("POST", "/", newUser(t, "john")))
	assert.Equal(t, http.StatusCreated, res2.Result().StatusCode)

	var user2 model.User
	assert.NoError(t, json.NewDecoder(res2.Body).Decode(&user2))
	assert.Equal(t, "john", user2.Username)

	var users []model.User
	t.Run("it returns a list of users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()

		query := req.URL.Query()
		query.Add("page", "0")
		query.Add("limit", "2")
		req.URL.RawQuery = query.Encode()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		assert.NoError(t, json.NewDecoder(res.Body).Decode(&users))
		assert.Equal(t, 2, len(users))
	})

	t.Run(fmt.Sprintf("it returns a single user %d", user1.ID), func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/%d", user1.ID), nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var result model.User
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&result))
		assert.Equal(t, user1.ID, result.ID)
	})

	t.Run("it returns accepted for modifying a user", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/1", newUser(t, "joshua"))
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusAccepted, res.Result().StatusCode)

		var user model.User
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&user))
		assert.Equal(t, "joshua", user.Username)
	})

	for _, user := range users {
		t.Run("it deletes and returns OK", func(t *testing.T) {

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/%d", user.ID), nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		})
	}

	t.Run("it returns a 404 for user not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/%d", user1.ID), nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		assert.Equal(t, 404, res.Result().StatusCode)
	})

	t.Cleanup(func() {
		//put code to clean up db here
		for _, user := range users {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/%d/purge", user.ID), nil)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		}
	})
}
