package usecase

import (
	"context"
	"log"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
)

type usecase struct {
	PullRequestRepository pull_requests.Repository
	UsersRepository       users.Repository
}

func NewUsecase(prRepository pull_requests.Repository, usRepository users.Repository) *usecase {
	return &usecase{PullRequestRepository: prRepository, UsersRepository: usRepository}
}

func (u *usecase) ReassignPullRequest(ctx context.Context, pullRequestID string, oldUserID string) (domain.PullRequest, string, error) {
	const op = "pull_request.Usecase.ReassignPullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, string, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, "", domain.NewError(code, message, err)
	}

	pr, err := u.PullRequestRepository.FetchByID(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	if pr.Status == domain.MERGED {
		return fail(domain.PR_MERGED, "PR was merged", err)
	}

	revs, err := u.PullRequestRepository.GetReviewersID(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
	}

	isReviewerAssigned := false
	for _, reviewer := range revs {
		if reviewer == oldUserID {
			isReviewerAssigned = true
			break
		}
	}

	if !isReviewerAssigned {
		return fail(domain.NOT_FOUND, "reviewer is not found", nil)
	}

	oldUser, err := u.UsersRepository.FetchByID(ctx, oldUserID)
	if err != nil {
		return fail(domain.NOT_FOUND, "user to replace not found", err)
	}

	activeUsers, err := u.UsersRepository.GetActiveUsersIDByTeam(ctx, oldUser.TeamName)
	if err != nil {
		return fail(domain.INTERNAL, "failed to get active team members", err)
	}

	currentReviewersSet := make(map[string]bool)
	for _, revID := range revs {
		currentReviewersSet[revID] = true
	}

	var newUserID string
	for _, userID := range activeUsers {
		if userID != oldUserID && !currentReviewersSet[userID] && userID != pr.AuthorID {
			newUserID = userID
			break
		}
	}

	if newUserID == "" {
		return fail(domain.NO_CANDIDATE, "no active replacement candidate in team", nil)
	}

	err = u.PullRequestRepository.DeleteReviewer(ctx, pullRequestID, oldUserID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	err = u.PullRequestRepository.AddReviewer(ctx, pullRequestID, newUserID)
	if err != nil {
		err := u.PullRequestRepository.AddReviewer(ctx, pullRequestID, oldUserID)
		if err != nil {
			return fail(domain.NOT_ASSIGNED, "failer to add new reviewer", err)
		}
		return fail(domain.NOT_ASSIGNED, "failed to add new reviewer", err)
	}

	updatedPR, err := u.PullRequestRepository.FetchByID(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource is not found", err)
	}

	return updatedPR, newUserID, nil
}

func (u *usecase) MergePullRequest(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	const op = "pull_request.Usecase.MergePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	pr, err := u.PullRequestRepository.FetchByIDWithMergeAt(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	if pr.Status == domain.MERGED {
		return pr, nil
	}

	newPr, err := u.PullRequestRepository.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	return newPr, nil
}

func (u *usecase) CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error) {
	const op = "pull_request.Usecase.CreatePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	user, err := u.UsersRepository.FetchByID(ctx, pullRequest.AuthorID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	teamMembersID, err := u.UsersRepository.GetActiveUsersIDByTeam(ctx, user.TeamName)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	var reviewers []string
	for _, teamMemberID := range teamMembersID {
		if teamMemberID != pullRequest.AuthorID {
			reviewers = append(reviewers, teamMemberID)
		}

		if len(reviewers) >= 2 {
			break
		}
	}

	pullRequest.Status = domain.OPEN

	err = u.PullRequestRepository.CreatePullRequest(ctx, pullRequest)
	if err != nil {
		return fail(domain.PR_EXISTS, "PR is already exists", err)
	}

	for _, rev := range reviewers {
		err = u.PullRequestRepository.AddReviewer(ctx, pullRequest.PullRequestID, rev)
		if err != nil {
			return fail(domain.NOT_ASSIGNED, "Not assigned after creating", err)
		}
	}

	pullRequest.AssignedReviewers = reviewers

	return *pullRequest, nil
}
