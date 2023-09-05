package caching

type Usecase interface {
	GetRepo() Repository
}

type CachingUsecase struct {
	repo Repository
}

func NewCachingUsecase(repo Repository) *CachingUsecase {
	return &CachingUsecase{
		repo: repo,
	}
}

func (u *CachingUsecase) GetRepo() Repository {
	return u.repo
}
