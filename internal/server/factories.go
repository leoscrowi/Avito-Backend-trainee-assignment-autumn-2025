package server

import (
	"github.com/jmoiron/sqlx"
	pr_ "github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests/delivery/http/v1"
	prr_ "github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests/repository/postgresql"
	prc_ "github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests/usecase"
	t_ "github.com/leoscrowi/pr-assignment-service/internal/app/teams/delivery/http/v1"
	tr_ "github.com/leoscrowi/pr-assignment-service/internal/app/teams/repository/postgresql"
	tc_ "github.com/leoscrowi/pr-assignment-service/internal/app/teams/usecase"
	u_ "github.com/leoscrowi/pr-assignment-service/internal/app/users/delivery/http/v1"
	ur_ "github.com/leoscrowi/pr-assignment-service/internal/app/users/repository/postgresql"
	uc_ "github.com/leoscrowi/pr-assignment-service/internal/app/users/usecase"
)

func GetControllers(db *sqlx.DB) []RouteSetup {
	ur := ur_.NewUsersRepository(db)
	prR := prr_.NewPullRequestsRepository(db)
	tr := tr_.NewTeamsRepository(db)

	uc := u_.NewUsersController(uc_.NewUsecase(ur, prR))
	prc := pr_.NewPullRequestController(prc_.NewUsecase(prR, ur))
	t := t_.NewTeamsController(tc_.NewUsecase(ur, tr))

	var res = make([]RouteSetup, 0, 3)
	res = append(res, uc)
	res = append(res, prc)
	res = append(res, t)

	return res
}
