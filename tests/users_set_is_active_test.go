package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUserIsActive_Success(t *testing.T) {
	teamName := "test_set_active_team"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_set_active", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_set_active", "username": "TestBob", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	setActiveRequest := map[string]interface{}{
		"user_id":   "test_u2_set_active",
		"is_active": false,
	}

	resp := helpers.PatchJSON(t, "/users/setIsActive", setActiveRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		User domain.User `json:"user"`
	}
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)

	assert.Equal(t, "test_u2_set_active", out.User.UserID)
	assert.Equal(t, "TestBob", out.User.Username)
	assert.Equal(t, teamName, out.User.TeamName)
	assert.False(t, out.User.IsActive, "user should be deactivated")
}

func TestSetUserIsActive_Reactivate(t *testing.T) {
	teamName := "test_reactivate_team"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_reactivate", "username": "TestCharlie", "is_active": false},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	reactivateRequest := map[string]interface{}{
		"user_id":   "test_u1_reactivate",
		"is_active": true,
	}

	resp := helpers.PatchJSON(t, "/users/setIsActive", reactivateRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		User domain.User `json:"user"`
	}
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)

	assert.Equal(t, "test_u1_reactivate", out.User.UserID)
	assert.Equal(t, "TestCharlie", out.User.Username)
	assert.True(t, out.User.IsActive, "user should be reactivated")
}

func TestSetUserIsActive_UserNotFound(t *testing.T) {
	setActiveRequest := map[string]interface{}{
		"user_id":   "non_existent_user",
		"is_active": false,
	}

	resp := helpers.PatchJSON(t, "/users/setIsActive", setActiveRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errorResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)

	assert.Contains(t, errorResp, "error")
	assert.NotEmpty(t, errorResp["error"])
}

func TestSetUserIsActive_Unauthorized(t *testing.T) {
	setActiveRequest := map[string]interface{}{
		"user_id":   "user_test",
		"is_active": false,
	}

	resp := helpers.PatchJSON(t, "/users/setIsActive", setActiveRequest, "")
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errorResp map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)

	assert.Contains(t, errorResp, "error")
	assert.NotEmpty(t, errorResp["error"])
}
