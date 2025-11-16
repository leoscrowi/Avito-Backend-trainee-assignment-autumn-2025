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

func TestTeamGet_Success(t *testing.T) {
	teamName := "test_get_team"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2", "username": "TestBob", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	type response struct {
		TeamName string `json:"team_name"`
	}
	respGet := helpers.GetJSON(t, "/team/get", response{TeamName: teamName}, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respGet.Body)

	helpers.RequireStatusCode(t, respGet, http.StatusOK)

	var result struct {
		TeamName string `json:"team_name"`
		Members  []struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	}

	require.NoError(t, json.NewDecoder(respGet.Body).Decode(&result), "decode response")

	assert.Equal(t, teamName, result.TeamName, "team name should match")
	assert.Len(t, result.Members, 2, "should have 2 team members")

	assert.Equal(t, "test_u1", result.Members[0].UserID, "first member user ID should match")
	assert.Equal(t, "TestAlice", result.Members[0].Username, "first member username should match")
	assert.Equal(t, "test_u2", result.Members[1].UserID, "second member user ID should match")
	assert.Equal(t, "TestBob", result.Members[1].Username, "second member username should match")
}

func TestTeamGet_EmptyTeamName(t *testing.T) {
	resp := helpers.GetJSON(t, "/team/get", map[string]string{"team_name": ""}, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTeamGet_Unauthorized(t *testing.T) {
	resp := helpers.GetJSON(t, "/team/get", map[string]string{"team_name": "some_team"}, "")
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	helpers.RequireStatusCode(t, resp, http.StatusUnauthorized)
}

func TestTeamGet_NotFound(t *testing.T) {
	nonExistentTeam := "non_existent_team_123"

	resp := helpers.GetJSON(t, "/team/get", map[string]string{"team_name": nonExistentTeam}, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	helpers.RequireStatusCode(t, resp, http.StatusNotFound)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	require.NoError(t, json.Unmarshal(helpers.ReadBody(t, resp), &errorResp), "decode error response")
	assert.NotEmpty(t, errorResp.Error.Message, "expected error message in response")
}
