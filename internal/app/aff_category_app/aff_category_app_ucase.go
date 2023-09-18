package category

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type categoryUCase struct {
	Repo interfaces.AffCategoryRepository
}

func (c categoryUCase) GetAllCategory(ctx context.Context, page, size int) (dto.AffCategoryResponseDto, error) {
	listCategory, err := c.CategoryRepository.GetAllCategory(ctx, page, size)
	if err != nil {
		return dtos.CategoryResponseDto{}, err
	}
	var categoryDtos []dtos.CategoryDto
	for i := range listCategory {
		if i >= size {
			continue
		}
		categoryDtos = append(categoryDtos, listCategory[i].ToCategoryDto())
	}
	nextPage := page
	if len(listCategory) > size {
		nextPage = page + 1
	}
	return dtos.CategoryResponseDto{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     categoryDtos,
	}, nil
}

func (c categoryUCase) GetAllAffCampaignInCategory(ctx context.Context, categoryId uint32, queryBy, order string, page, size int) (dto.AffCategoryResponseDto, error) {
	orderBy := util.BuildOrderByC(queryBy, order)
	listCouponInCategory, err := c.CategoryRepository.GetAllCouponInCategory(ctx, categoryId, orderBy, page, size)
	if err != nil {
		return dtos.CouponDtoResponse{}, err
	}
	nextPage := page
	if len(listCouponInCategory) > size {
		nextPage = page + 1
	}
	var listCoupon []dtos.CouponDto

	for i := range listCouponInCategory {
		if i >= size {
			break
		}
		listCoupon = append(listCoupon, listCouponInCategory[i].ToCouponDto())
	}

	return dto.AffCategoryResponseDto{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data:     listCoupon,
	}, nil
}

func NewAffCategoryUCase(repo interfaces.AffCategoryRepository) interfaces.AffCategoryUCase {
	return &categoryUCase{
		Repo: repo,
	}
}
