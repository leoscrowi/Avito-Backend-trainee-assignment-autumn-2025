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

func TestPullRequestReassign_Success(t *testing.T) {
	teamName := "test_pr_team_success"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_success", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_success", "username": "TestBob", "is_active": true},
			{"user_id": "test_u3_success", "username": "TestCharlie", "is_active": true},
			{"user_id": "test_u4_success", "username": "TestSarah", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "pr-1001_success",
		"pull_request_name": "Add search",
		"author_id":         "test_u1_success",
	}

	createResp := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	_ = createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	reassignRequest := map[string]interface{}{
		"pull_request_id": "pr-1001_success",
		"old_user_id":     "test_u2_success",
	}

	resp := helpers.PatchJSON(t, "/pullRequest/reassign", reassignRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		PR         domain.PullRequest `json:"pr"`
		ReplacedBy string             `json:"replaced_by"`
	}
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)

	assert.Equal(t, "pr-1001_success", out.PR.PullRequestID)
	assert.NotEmpty(t, out.ReplacedBy)
	assert.NotContains(t, out.PR.AssignedReviewers, "test_u2_success", "old reviewer should not be in assigned reviewers")
	assert.Len(t, out.PR.AssignedReviewers, 2, "should have same number of reviewers")
}

func TestPullRequestReassign_NotFound(t *testing.T) {
	reassignRequest := map[string]interface{}{
		"pull_request_id": "pr-nonexistent_notfound",
		"old_user_id":     "u2_notfound",
	}

	resp := helpers.PatchJSON(t, "/pullRequest/reassign", reassignRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	err := json.Unmarshal(helpers.ReadBody(t, resp), &errorResp)
	require.NoError(t, err)

	assert.True(t, errorResp.Error.Code == domain.NOT_FOUND,
		"expected NOT_FOUND, got: ", errorResp.Error.Code)
}

func TestPullRequestReassign_NotAssigned(t *testing.T) {
	teamName := "test_pr_team_notassigned"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_notassigned", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_notassigned", "username": "TestBob", "is_active": true},
			{"user_id": "test_u3_notassigned", "username": "TestCharlie", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "pr-1002_notassigned",
		"pull_request_name": "Fix bug",
		"author_id":         "test_u1_notassigned",
	}

	createResp := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	_ = createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	mergeResp := helpers.PatchJSON(t, "/pullRequest/merge", map[string]interface{}{"pull_request_id": "pr-1002_notassigned"}, helpers.AdminToken)
	_ = mergeResp.Body.Close()

	reassignRequest := map[string]interface{}{
		"pull_request_id": "pr-1002_notassigned",
		"old_user_id":     "test_u2_notassigned",
	}

	resp := helpers.PatchJSON(t, "/pullRequest/reassign", reassignRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	err := json.Unmarshal(helpers.ReadBody(t, resp), &errorResp)
	require.NoError(t, err)

	assert.Equal(t, domain.PR_MERGED, errorResp.Error.Code)
}

func TestPullRequestReassign_NoCandidate(t *testing.T) {
	teamName := "test_pr_team_nocandidate"
	team := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "test_u1_nocandidate", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2_nocandidate", "username": "TestBob", "is_active": true},
			{"user_id": "test_u3_nocandidate", "username": "TestCharlie", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "pr-1004_nocandidate",
		"pull_request_name": "Special feature",
		"author_id":         "test_u1_nocandidate",
	}

	createResp := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	_ = createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	reassignRequest := map[string]interface{}{
		"pull_request_id": "pr-1004_nocandidate",
		"old_user_id":     "test_u2_nocandidate",
	}

	resp := helpers.PatchJSON(t, "/pullRequest/reassign", reassignRequest, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	err := json.Unmarshal(helpers.ReadBody(t, resp), &errorResp)
	require.NoError(t, err)

	assert.Equal(t, domain.NO_CANDIDATE, errorResp.Error.Code)
}

func TestPullRequestReassign_InvalidToken(t *testing.T) {
	reassignRequest := map[string]interface{}{
		"pull_request_id": "pr-1001_invalidtoken",
		"old_user_id":     "u2_invalidtoken",
	}

	resp := helpers.PostJSON(t, "/pullRequest/reassign", reassignRequest, "")
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
