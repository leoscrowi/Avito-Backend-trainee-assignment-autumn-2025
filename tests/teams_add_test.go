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

func TestTeamAdd_Created(t *testing.T) {
	team := map[string]interface{}{
		"team_name": "payments",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
		},
	}
	resp := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	helpers.RequireStatusCode(t, resp, http.StatusCreated)

	var out struct {
		Team domain.Team `json:"team"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out), "decode response")

	assert.Equal(t, team["team_name"], out.Team.TeamName, "team name should match")
	assert.Len(t, out.Team.Members, 2, "should have 2 team members")
}

func TestTeamAdd_AlreadyExists(t *testing.T) {
	teamName := "test_existing_1"

	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u3", "username": "TestCharlie", "is_active": true},
		},
	}

	resp1 := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = resp1.Body.Close()
	helpers.RequireStatusCode(t, resp1, http.StatusCreated)

	resp2 := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp2.Body)

	helpers.RequireStatusCode(t, resp2, http.StatusBadRequest)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	require.NoError(t, json.Unmarshal(helpers.ReadBody(t, resp2), &errorResp), "decode error response")
	assert.Equal(t, domain.TEAM_EXISTS, errorResp.Error.Code, "error code should be TEAM_EXISTS")
}
