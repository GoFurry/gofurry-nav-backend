package controller

import (
	"github.com/GoFurry/gofurry-nav-backend/apps/system/stat/service"
	"github.com/GoFurry/gofurry-nav-backend/common"
	"github.com/gofiber/fiber/v2"
)

type statApi struct{}

var StatApi *statApi

func init() {
	StatApi = &statApi{}
}

// @Summary 增加浏览数
// @Schemes
// @Description 增加浏览数
// @Tags Stat
// @Accept json
// @Produce json
// @Success 200 {object} common.ResultData
// @Router /api/stat/add/count [Get]
func (api *statApi) AddCount(c *fiber.Ctx) error {
	// TODO:IP限流和令牌桶限流
	err := service.GetStatService().AddViewCount()
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).Success()
}

// @Summary 查询浏览数
// @Schemes
// @Description 查询浏览数
// @Tags Stat
// @Accept json
// @Produce json
// @Success 200 {object} common.ResultData
// @Router /api/stat/chart/views/count [Get]
func (api *statApi) GetViewsCount(c *fiber.Ctx) error {
	data, err := service.GetStatService().ViewsCount()
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}

// @Summary 查询内容最多的分组
// @Schemes
// @Description 查询内容最多的分组
// @Tags Stat
// @Accept json
// @Produce json
// @Success 200 {object} common.ResultData
// @Router /api/stat/chart/group/count [Get]
func (api *statApi) GetGroupCount(c *fiber.Ctx) error {
	data, err := service.GetStatService().GroupCount()
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}

// @Summary 查询访问最多的国家
// @Schemes
// @Description 查询访问最多的国家
// @Tags Stat
// @Accept json
// @Produce json
// @Success 200 {object} common.ResultData
// @Router /api/stat/chart/views/region/country [Get]
func (api *statApi) GetCountryCount(c *fiber.Ctx) error {
	data, err := service.GetStatService().CountryCount()
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}

// @Summary 查询访问最多的省
// @Schemes
// @Description 查询访问最多的省
// @Tags Stat
// @Accept json
// @Produce json
// @Success 200 {object} common.ResultData
// @Router /api/stat/chart/views/region/province [Get]
func (api *statApi) GetProvinceCount(c *fiber.Ctx) error {
	data, err := service.GetStatService().ProvinceCount()
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}

// @Summary 查询访问最多的城市
// @Schemes
// @Description 查询访问最多的城市
// @Tags Stat
// @Accept json
// @Produce json
// @Success 200 {object} common.ResultData
// @Router /api/stat/chart/views/region/city [Get]
func (api *statApi) GetCityCount(c *fiber.Ctx) error {
	data, err := service.GetStatService().CityCount()
	if err != nil {
		return common.NewResponse(c).Error(err.GetMsg())
	}

	return common.NewResponse(c).SuccessWithData(data)
}
