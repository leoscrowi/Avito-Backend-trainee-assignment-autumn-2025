package usecase

import (
	"context"
	"fmt"

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

	pr, err := u.PullRequestRepository.FetchByID(ctx, pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if pr.Status == domain.MERGED {
		return domain.PullRequest{}, fmt.Errorf("Merged")
	}

	revs, err := u.PullRequestRepository.GetReviewersID(ctx, pullRequestID)
	for _, rev := range revs {
		if rev == oldUserID {
			err = u.PullRequestRepository.DeleteReviewer(ctx, pullRequestID, oldUserID)
			if err != nil {
				return domain.PullRequest{}, err
			}

			user, err := u.UsersRepository.FetchByID(ctx, oldUserID)
			if err != nil {
				return domain.PullRequest{}, err
			}

			activeUsers, err := u.UsersRepository.GetActiveUsersIDByTeam(ctx, user.TeamName)
			if err != nil {
				return domain.PullRequest{}, err
			}

			if len(activeUsers) == 0 {
				return domain.PullRequest{}, fmt.Errorf("no active users")
			}

			var addUserID string
			for _, id := range activeUsers {
				if id != oldUserID {
					addUserID = id
				}
			}

			err = u.PullRequestRepository.AddReviewer(ctx, pullRequestID, addUserID)
			if err != nil {
				return domain.PullRequest{}, err
			}

			pr, err := u.PullRequestRepository.FetchByID(ctx, pullRequestID)
			if err != nil {
				return domain.PullRequest{}, err
			}

			return pr, nil
		}
	}
	return domain.PullRequest{}, fmt.Errorf("can't find user or pr")
}

func (u *Usecase) MergePullRequest(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	pr, err := u.PullRequestRepository.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return pr, err
}

func (u *Usecase) CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error) {
	user, err := u.UsersRepository.FetchByID(ctx, pullRequest.AuthorID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	teamMembersID, err := u.UsersRepository.GetActiveUsersIDByTeam(ctx, user.TeamName)
	if err != nil {
		return domain.PullRequest{}, err
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
	pullRequest.Status = domain.MERGED

	err = u.PullRequestRepository.CreatePullRequest(ctx, pullRequest)
	if err != nil {
		return domain.PullRequest{}, err
	}
	return *pullRequest, nil
}
