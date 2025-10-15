package routers

import (
	nav "github.com/GoFurry/gofurry-nav-backend/apps/nav/navPage/controller"
	site "github.com/GoFurry/gofurry-nav-backend/apps/nav/sitePage/controller"
	stat "github.com/GoFurry/gofurry-nav-backend/apps/system/stat/controller"
	user "github.com/GoFurry/gofurry-nav-backend/apps/system/user/controller"
	"github.com/gofiber/fiber/v2"
)

/*
 * @Desc: 接口层
 * @author: 福狼
 * @version: v1.0.0
 */

func userApi(g fiber.Router) {
	// 非鉴权
	g.Post("/login", user.UserApi.Login)
	// 鉴权

}

// 导航相关
func navApi(g fiber.Router) {
	// 导航页
	// 导航部分
	g.Get("/page/site/list", nav.NavPageApi.GetSiteList)   // 获取所有导航站点信息
	g.Get("/page/group/list", nav.NavPageApi.GetGroupList) // 获取所有导航站点分组信息
	g.Get("/page/ping/list", nav.NavPageApi.GetPingList)   // 获取所有导航站点延迟信息
	// 导航页头部组件
	g.Get("/page/search/baidu", nav.NavPageApi.GetBaiduSearchSuggestion)       // 解析百度搜索建议框
	g.Get("/page/search/bing", nav.NavPageApi.GetBingSearchSuggestion)         // 解析必应搜索建议框
	g.Get("/page/search/google", nav.NavPageApi.GetGoogleSearchSuggestion)     // 解析谷歌搜索建议框
	g.Get("/page/search/bilibili", nav.NavPageApi.GetBiliBiliSearchSuggestion) // 解析B站搜索建议框
	g.Get("/page/header/getSaying", nav.NavPageApi.GetSaying)                  // 提供随机金句
	g.Get("/page/header/getImage", nav.NavPageApi.GetImage)                    // 提供背景随机图片

	// 详情页
	g.Get("/site/getSiteDetail", site.SitePageApi.GetSiteDetail)         // 获取单个站点的信息
	g.Get("/site/getSitePingRecord", site.SitePageApi.GetSitePingRecord) // 获取单个站点的 Ping 记录
	g.Get("/site/getSiteHttpRecord", site.SitePageApi.GetSiteHttpRecord) // 获取单个站点的 HTTP 记录
	g.Get("/site/getSiteDnsRecord", site.SitePageApi.GetSiteDnsRecord)   // 获取单个站点的 DNS 记录

	// 鉴权

}

func statApi(g fiber.Router) {
	// 数据
	g.Get("/add/count", stat.StatApi.AddCount)                             // 增加访问量
	g.Get(("/chart/views/count"), stat.StatApi.GetViewsCount)              // 获取访问量数据
	g.Get(("/chart/views/region/country"), stat.StatApi.GetCountryCount)   // 获取访问国家统计
	g.Get(("/chart/views/region/province"), stat.StatApi.GetProvinceCount) // 获取访问省份统计
	g.Get(("/chart/views/region/city"), stat.StatApi.GetCityCount)         // 获取访问城市统计
	g.Get(("/chart/group/count"), stat.StatApi.GetGroupCount)              // 获取站点分组统计

	// 鉴权

}
