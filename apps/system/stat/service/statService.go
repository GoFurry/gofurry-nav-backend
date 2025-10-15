package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoFurry/gofurry-nav-backend/apps/system/stat/dao"
	"github.com/GoFurry/gofurry-nav-backend/apps/system/stat/models"
	"github.com/GoFurry/gofurry-nav-backend/common"
	"github.com/GoFurry/gofurry-nav-backend/common/log"
	cs "github.com/GoFurry/gofurry-nav-backend/common/service"
	"github.com/GoFurry/gofurry-nav-backend/common/util"
)

type statService struct{}

var statSingleton = new(statService)

func GetStatService() *statService { return statSingleton }

// 增加访问量
func (s statService) AddViewCount() common.GFError {
	cs.Incr("stat-count:total") // 总访问量
	year := util.Int642String(int64(time.Now().Year()))
	month := util.Int642String(int64(time.Now().Month()))
	day := util.Int642String(int64(time.Now().Day()))
	cs.Incr("stat-count:" + year)                           // 年访问量
	cs.Incr("stat-count:" + year + "-" + month)             // 月访问量
	cs.Incr("stat-count:" + year + "-" + month + "-" + day) // 日访问量

	return nil
}

// 返回访问量统计
func (s statService) ViewsCount() (res models.ViewsCountVo, err common.GFError) {
	// 获取时间
	now := time.Now()
	year := now.Year()
	month := now.Month()

	// 获取 redis 缓存
	getViewsCount(&res, year, month)

	// 获取最近7日浏览量
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		key := fmt.Sprintf("stat-count:%d-%d-%d", day.Year(), int(day.Month()), day.Day())

		val, gfsError := cs.GetString(key)
		if gfsError != nil {
			log.Warn("redis未找到key:", key)
			res.Date = append(res.Date, day.Format("2006-01-02"))
			res.Count = append(res.Count, 0)
			continue
		}

		intVal, utilErr := util.String2Int64(val)
		if utilErr != nil {
			intVal = 0
		}

		res.Date = append(res.Date, day.Format("2006-01-02"))
		res.Count = append(res.Count, intVal)
	}

	return
}

// 获取数量最多的分组
func (s statService) GroupCount() (res []models.GroupCountVo, err common.GFError) {
	return dao.GetStatDao().GetGroupCount()
}

// 获取访问国家统计
func (s statService) CountryCount() (res models.RegionCountVo, err common.GFError) {
	res.RegionMap = readTopCache("stat-geoip-country:top")
	return
}

// 获取访问省份统计
func (s statService) ProvinceCount() (res models.RegionCountVo, err common.GFError) {
	res.RegionMap = readTopCache("stat-geoip-province:top")
	return
}

// 获取访问城市统计
func (s statService) CityCount() (res models.RegionCountVo, err common.GFError) {
	res.RegionMap = readTopCache("stat-geoip-city:top")
	return
}

// 读取缓存
func readTopCache(cacheKey string) map[string]int64 {
	var res map[string]int64
	if val, err := cs.GetString(cacheKey); err == nil && val != "" {
		_ = json.Unmarshal([]byte(val), &res)
	}
	if res == nil {
		res = make(map[string]int64)
	}
	return res
}

func getViewsCount(res *models.ViewsCountVo, year int, month time.Month) {
	total, gfsError := cs.GetString("stat-count:total")
	if gfsError != nil {
		log.Error("stat-count:total获取失败: ", gfsError)
	}
	intTotal, utilErr := util.String2Int64(total)
	if utilErr != nil {
		log.Error("String转换Int64获取失败: ", gfsError)
	}
	res.Total = intTotal

	strYear := util.Int642String(int64(year))
	yearCount, gfsError := cs.GetString("stat-count:" + strYear)
	if gfsError != nil {
		log.Error("stat-count:"+strYear+"获取失败: ", strYear, gfsError)
	}
	intYearCount, utilErr := util.String2Int64(yearCount)
	if utilErr != nil {
		log.Error("String转换Int64获取失败: ", gfsError)
	}
	res.YearCount = intYearCount

	strMonth := util.Int642String(int64(month))
	monthCount, gfsError := cs.GetString("stat-count:" + strYear + "-" + strMonth)
	if gfsError != nil {
		log.Error("stat-count:"+strYear+"-"+strMonth+"获取失败: ", gfsError)
	}
	intMonthCount, utilErr := util.String2Int64(monthCount)
	if utilErr != nil {
		log.Error("String转换Int64获取失败: ", gfsError)
	}
	res.MonthCount = intMonthCount
}
