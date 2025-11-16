package usecase

import (
	"context"
	"log"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/stats"
)

type Usecase struct {
	StatsRepository stats.Repository
}

func NewUsecase(sRepository stats.Repository) *Usecase {
	return &Usecase{StatsRepository: sRepository}
}

func (u *Usecase) GetPullRequestStats(ctx context.Context) ([]domain.PullRequestStats, error) {
	const op = "users.Usecase.GetPullRequestStats"

	fail := func(code domain.ErrorCode, message string, err error) ([]domain.PullRequestStats, error) {
		log.Printf("%s, %v\n", op, err)
		return nil, domain.NewError(code, message, err)
	}

	stats, err := u.StatsRepository.GetPullRequestStats(ctx)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return stats, nil
}
