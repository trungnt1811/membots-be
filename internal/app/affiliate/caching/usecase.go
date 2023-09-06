package caching

type Usecase interface {
	GetRepo() Repository
}

type cachingUsecase struct {
	repo Repository
}

func NewCachingUsecase(repo Repository) *cachingUsecase {
	return &cachingUsecase{
		repo: repo,
	}
}

func (u *cachingUsecase) GetRepo() Repository {
	return u.repo
}
