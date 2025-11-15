package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
)

type Usecase struct {
	PullRequestRepository pull_requests.Repository
	UsersRepository       users.Repository
}

func NewUsecase(prRepository pull_requests.Repository, usRepository users.Repository) *Usecase {
	return &Usecase{PullRequestRepository: prRepository, UsersRepository: usRepository}
}

func (u *Usecase) ReassignPullRequest(ctx context.Context, pullRequestID string, oldUserID string) (domain.PullRequest, error) {
	const op = "pull_request.Usecase.ReassignPullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	pr, err := u.PullRequestRepository.FetchByID(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	if pr.Status == domain.MERGED {
		return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
	}

	revs, err := u.PullRequestRepository.GetReviewersID(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
	}
	for _, rev := range revs {
		if rev == oldUserID {
			err = u.PullRequestRepository.DeleteReviewer(ctx, pullRequestID, oldUserID)
			if err != nil {
				return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
			}

			user, err := u.UsersRepository.FetchByID(ctx, oldUserID)
			if err != nil {
				return fail(domain.NOT_FOUND, "resource not found", err)
			}

			activeUsers, err := u.UsersRepository.GetActiveUsersIDByTeam(ctx, user.TeamName)
			if err != nil {
				return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
			}

			if len(activeUsers) == 0 {
				return fail(domain.NO_CANDIDATE, "no active replacement candidate in team", err)
			}

			var addUserID string
			for _, id := range activeUsers {
				if id != oldUserID {
					addUserID = id
				}
			}

			if addUserID == "" {
				return fail(domain.NO_CANDIDATE, "no active replacement candidate in team", err)
			}

			err = u.PullRequestRepository.AddReviewer(ctx, pullRequestID, addUserID)
			if err != nil {
				return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
			}

			pr, err := u.PullRequestRepository.FetchByID(ctx, pullRequestID)
			if err != nil {
				return fail(domain.NOT_ASSIGNED, "reviewer is not assigned to this PR", err)
			}

			return pr, nil
		}
	}
	return domain.PullRequest{}, fmt.Errorf("can't find user or pr")
}

func (u *Usecase) MergePullRequest(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	const op = "pull_request.Usecase.MergePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	pr, err := u.PullRequestRepository.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	return pr, nil
}

func (u *Usecase) CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error) {
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

	pullRequest.AssignedReviewers = reviewers
	pullRequest.Status = domain.OPEN

	err = u.PullRequestRepository.CreatePullRequest(ctx, pullRequest)
	if err != nil {
		return fail(domain.PR_EXISTS, "PR is already exists", err)
	}
	return *pullRequest, nil
}
