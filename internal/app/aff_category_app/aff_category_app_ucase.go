package category

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
)

type categoryUCase struct {
	Repo interfaces.AffCategoryRepository
}

func (c *categoryUCase) GetAllCategory(ctx context.Context, page, size int) (dto.AffCategoryResponseDto, error) {
	listCategory, err := c.Repo.GetAllCategory(ctx, page, size)
	if err != nil {
		return dto.AffCategoryResponseDto{}, err
	}
	var categoryDtos []dto.AffCategoryDto
	for i := range listCategory {
		if i >= size {
			continue
		}
		categoryDtos = append(categoryDtos, listCategory[i].ToDto())
	}
	nextPage := page
	if len(listCategory) > size {
		nextPage = page + 1
	}
	return dto.AffCategoryResponseDto{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     categoryDtos,
	}, nil
}

func (c *categoryUCase) GetAllAffCampaignInCategory(ctx context.Context, categoryId uint32, queryBy, order string, page, size int) (dto.AffCampaignAppDtoResponse, error) {
	orderBy := util.BuildOrderBy(queryBy, order)
	listCouponInCategory, err := c.Repo.GetAllAffCampaignInCategory(ctx, categoryId, orderBy, page, size)
	if err != nil {
		return dto.AffCampaignAppDtoResponse{}, err
	}
	nextPage := page
	if len(listCouponInCategory) > size {
		nextPage = page + 1
	}
	var listCampaign []dto.AffCampaignLessDto

	for i := range listCouponInCategory {
		if i >= size {
			break
		}
		listCampaign = append(listCampaign, listCouponInCategory[i].ToDto())
	}

	return dto.AffCampaignAppDtoResponse{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listCampaign,
	}, nil
}

func NewAffCategoryUCase(repo interfaces.AffCategoryRepository) interfaces.AffCategoryUCase {
	return &categoryUCase{
		Repo: repo,
	}
}
