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

func TestPullRequestCreate_Success(t *testing.T) {
	team := map[string]interface{}{
		"team_name": "test_pr_team",
		"members": []map[string]interface{}{
			{"user_id": "test_u1", "username": "TestAlice", "is_active": true},
			{"user_id": "test_u2", "username": "TestBob", "is_active": true},
			{"user_id": "test_u3", "username": "TestCharlie", "is_active": true},
			{"user_id": "test_u4", "username": "TestSarah", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "test_pr_1",
		"pull_request_name": "Test PR",
		"author_id":         "test_u1",
	}

	respCreate := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respCreate.Body)

	helpers.RequireStatusCode(t, respCreate, http.StatusCreated)

	var result struct {
		PR struct {
			PullRequestID     string   `json:"pull_request_id"`
			PullRequestName   string   `json:"pull_request_name"`
			AuthorID          string   `json:"author_id"`
			Status            string   `json:"status"`
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
	}
	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&result), "decode response")

	assert.Equal(t, "test_pr_1", result.PR.PullRequestID, "pull_request_id should match")
	assert.Equal(t, "Test PR", result.PR.PullRequestName, "pull_request_name should match")
	assert.Equal(t, "test_u1", result.PR.AuthorID, "author_id should match")
	assert.Equal(t, "OPEN", result.PR.Status, "status should be OPEN")

	assert.Len(t, result.PR.AssignedReviewers, 2, "should have 2 assigned reviewers")

	for _, reviewer := range result.PR.AssignedReviewers {
		assert.NotEqual(t, "test_u1", reviewer, "author should not be assigned as reviewer")
	}

	expectedReviewers := []string{"test_u2", "test_u3"}
	assert.ElementsMatch(t, expectedReviewers, result.PR.AssignedReviewers, "reviewers should match expected")
}

func TestPullRequestCreate_AuthorNotFound(t *testing.T) {
	pr := map[string]interface{}{
		"pull_request_id":   "test_pr_2",
		"pull_request_name": "Test PR",
		"author_id":         "non_existent_user",
	}

	respCreate := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respCreate.Body)

	helpers.RequireStatusCode(t, respCreate, http.StatusNotFound)

	var errorResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&errorResp), "decode error response")
	assert.NotEmpty(t, errorResp.Error.Message, "expected error message in response")
}

func TestPullRequestCreate_TeamNotFound(t *testing.T) {
	pr := map[string]interface{}{
		"pull_request_id":   "test_pr_3",
		"pull_request_name": "Test PR",
		"author_id":         "lonely_user",
	}

	respCreate := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respCreate.Body)

	helpers.RequireStatusCode(t, respCreate, http.StatusNotFound)
}

func TestPullRequestCreate_DuplicatePR(t *testing.T) {
	team := map[string]interface{}{
		"team_name": "test_duplicate_team",
		"members": []map[string]interface{}{
			{"user_id": "test_dup_u1", "username": "TestDupAlice", "is_active": true},
			{"user_id": "test_dup_u2", "username": "TestDupBob", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "duplicate_pr",
		"pull_request_name": "First PR",
		"author_id":         "test_dup_u1",
	}

	respCreate1 := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	_ = respCreate1.Body.Close()
	helpers.RequireStatusCode(t, respCreate1, http.StatusCreated)

	respCreate2 := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respCreate2.Body)

	helpers.RequireStatusCode(t, respCreate2, http.StatusConflict)

	var errorResp struct {
		Error struct {
			Code    domain.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}

	require.NoError(t, json.NewDecoder(respCreate2.Body).Decode(&errorResp), "decode error response")
	assert.Equal(t, domain.PR_EXISTS, errorResp.Error.Code, "error code should be PR_EXISTS")
}

func TestPullRequestCreate_NotEnoughReviewers(t *testing.T) {
	team := map[string]interface{}{
		"team_name": "test_solo_team",
		"members": []map[string]interface{}{
			{"user_id": "solo_user", "username": "SoloUser", "is_active": true},
		},
	}

	respAdd := helpers.PostJSON(t, "/team/add", team, helpers.AdminToken)
	_ = respAdd.Body.Close()
	helpers.RequireStatusCode(t, respAdd, http.StatusCreated)

	pr := map[string]interface{}{
		"pull_request_id":   "solo_pr",
		"pull_request_name": "Solo PR",
		"author_id":         "solo_user",
	}

	respCreate := helpers.PostJSON(t, "/pullRequest/create", pr, helpers.AdminToken)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(respCreate.Body)

	helpers.RequireStatusCode(t, respCreate, http.StatusCreated)

	var result struct {
		PR struct {
			PullRequestID     string   `json:"pull_request_id"`
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
	}

	require.NoError(t, json.NewDecoder(respCreate.Body).Decode(&result), "decode response")
	assert.Empty(t, result.PR.AssignedReviewers, "should have 0 assigned reviewers in solo team")
}
