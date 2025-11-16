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

func TestGetUserReview_EmptyList(t *testing.T) {
	teamName := "test_get_review_team_empty"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_review_empty", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_review_empty", "username": "TestBob", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	req := map[string]interface{}{
		"user_id": "test_u1_review_empty",
	}
	resp := helpers.GetJSON(t, "/users/getReview", req, helpers.UserToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		PR []domain.PullRequest `json:"pull_requests"`
	}
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)
	assert.Equal(t, 0, len(out.PR))
}

func TestGetUserReview_NotEmptyList(t *testing.T) {
	teamName := "test_get_review_team_not_empty"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_review_not_empty", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_review_not_empty", "username": "TestBob", "is_active": true},
			{"user_id": "test_u3_review_not_empty", "username": "TestCharlie", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "test_pr_1_not_empty",
		"pull_request_name": "Test PR",
		"author_id":         "test_u1_review_not_empty",
	}

	respCreate := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respCreate.Body)
	helpers.RequireStatusCode(t, respCreate, http.StatusCreated)

	req := map[string]interface{}{
		"user_id": "test_u2_review_not_empty",
	}
	respGet := helpers.GetJSON(t, "/users/getReview", req, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respGet.Body)

	require.Equal(t, http.StatusOK, respGet.StatusCode)

	var out struct {
		PR []domain.PullRequest `json:"pull_requests"`
	}
	err := json.NewDecoder(respGet.Body).Decode(&out)
	require.NoError(t, err)
	assert.Equal(t, 1, len(out.PR))
}
