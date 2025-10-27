package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GoFurry/gofurry-nav-backend/apps/nav/navPage/dao"
	"github.com/GoFurry/gofurry-nav-backend/apps/nav/navPage/models"
	"github.com/GoFurry/gofurry-nav-backend/common"
	cs "github.com/GoFurry/gofurry-nav-backend/common/service"
	"github.com/GoFurry/gofurry-nav-backend/common/util"
	"github.com/GoFurry/gofurry-nav-backend/roof/env"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type navPageService struct{}

var navPageSingleton = new(navPageService)

func GetNavPageService() *navPageService { return navPageSingleton }

// Logo 支持的图片格式
var validExts = env.GetServerConfig().Resource.ImageExts

// 获取导航站点信息
func (svc *navPageService) GetSiteList(lang string) (res []models.SiteVo, err common.GFError) {
	records, gfError := dao.GetNavPageDao().GetSiteList() // 所有站点记录
	if gfError != nil {
		return nil, gfError
	}

	switch lang {
	case "zh":
		for _, v := range records {
			newRecord := models.SiteVo{
				ID:      util.Int642String(v.ID),
				Name:    v.Name,
				Info:    v.Info,
				Domain:  v.Domain,
				Country: v.Country,
				Nsfw:    v.Nsfw,
				Welfare: v.Welfare,
				Icon:    v.Icon,
			}
			res = append(res, newRecord)
		}
	case "en":
		for _, v := range records {
			newRecord := models.SiteVo{
				ID:      util.Int642String(v.ID),
				Name:    v.NameEn,
				Info:    v.InfoEn,
				Domain:  v.Domain,
				Country: v.Country,
				Nsfw:    v.Nsfw,
				Welfare: v.Welfare,
				Icon:    v.Icon,
			}
			res = append(res, newRecord)
		}
	default:
		for _, v := range records {
			newRecord := models.SiteVo{
				ID:      util.Int642String(v.ID),
				Name:    v.Name,
				Info:    v.Info,
				Domain:  v.Domain,
				Country: v.Country,
				Nsfw:    v.Nsfw,
				Welfare: v.Welfare,
				Icon:    v.Icon,
			}
			res = append(res, newRecord)
		}
	}

	// 成功查询站点增加浏览量
	addCount()

	return res, nil
}

// 获取导航站点分组信息
func (svc *navPageService) GetGroupList(lang string) (res []models.GroupVo, err common.GFError) {
	groupRecords, gfError := dao.GetNavPageDao().GetGroupList() // 所有分组记录
	if gfError != nil {
		return nil, gfError
	}
	mappingRecords, gfError := dao.GetNavPageDao().GetGroupMapList() // 所有分组-站点映射
	if gfError != nil {
		return nil, gfError
	}

	idList := []int64{}
	voMap := map[int64]*models.GroupVo{} // 准备返回的VO
	for _, v := range groupRecords {
		idList = append(idList, v.ID)
		voMap[v.ID] = &models.GroupVo{
			ID:       util.Int642String(v.ID),
			Priority: v.Priority,
			Sites:    []string{},
		}

		// 默认返回中文
		switch lang {
		case "zh":
			voMap[v.ID].Name = v.Name
			voMap[v.ID].Info = v.Info
		case "en":
			voMap[v.ID].Name = v.NameEn
			voMap[v.ID].Info = v.InfoEn
		default:
			voMap[v.ID].Name = v.Name
			voMap[v.ID].Info = v.Info
		}

	}

	// 根据映射把站点放入分组
	for _, v := range mappingRecords {
		voMap[v.GroupID].Sites = append(voMap[v.GroupID].Sites, util.Int642String(v.SiteID))
	}
	for _, v := range idList {
		res = append(res, *voMap[v])
	}

	return
}

// 获取导航站点延迟信息
func (svc *navPageService) GetPingList() (res map[string]string, err common.GFError) {
	return cs.HGetAll("ping:result")
}

func (svc *navPageService) GetBaiduSuggestion(q string) ([]string, common.GFError) {
	url := fmt.Sprintf("http://suggestion.baidu.com/su?wd=%s&p=3&cb=window.bdsug.sug", q)
	resp, err := http.Get(url)
	if err != nil {
		return nil, common.NewServiceError("请求百度搜索建议接口出错: " + err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// 转 GBK -> UTF-8
	reader := transform.NewReader(bytes.NewReader(body), simplifiedchinese.GBK.NewDecoder())
	utf8Body, _ := ioutil.ReadAll(reader)
	strBody := string(utf8Body)

	// 提取 JSON 字符串
	prefix := "window.bdsug.sug("
	suffix := ");"
	start := strings.Index(strBody, prefix)
	end := strings.LastIndex(strBody, suffix)
	if start == -1 || end == -1 || end <= start+len(prefix) {
		return []string{}, nil
	}
	jsonStr := strBody[start+len(prefix) : end]

	// 把非标准 JSON 的键加上双引号
	replacer := strings.NewReplacer(
		"q:", `"q":`,
		"p:", `"p":`,
		"s:", `"s":`,
	)
	jsonStr = replacer.Replace(jsonStr)

	// 定义结构体
	type BaiduResp struct {
		S []string `json:"s"`
	}

	var result BaiduResp
	if jsonErr := json.Unmarshal([]byte(jsonStr), &result); jsonErr != nil {
		return []string{}, nil
	}

	return result.S, nil
}

type BingResponse struct {
	AS struct {
		Results []struct {
			Suggests []struct {
				Txt string `json:"Txt"`
			} `json:"Suggests"`
		} `json:"Results"`
	} `json:"AS"`
}

func (svc *navPageService) GetBingSuggestion(q string) ([]string, common.GFError) {
	url := fmt.Sprintf("https://api.bing.com/qsonhs.aspx?type=cb&q=%s", q)
	resp, err := http.Get(url)
	if err != nil {
		return nil, common.NewServiceError("请求必应搜索建议接口出错: " + err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	strBody := string(body)

	// 固定前缀和固定尾部
	prefix := "if(typeof  == 'function') ("
	suffix := "/* pageview_candidate */);"

	start := strings.Index(strBody, prefix)
	end := strings.LastIndex(strBody, suffix)
	if start == -1 || end == -1 || end <= start+len(prefix) {
		return []string{}, nil
	}

	jsonStr := strBody[start+len(prefix) : end]

	// 定义结构体解析
	type Suggest struct {
		Txt string `json:"Txt"`
	}
	type Result struct {
		Suggests []Suggest `json:"Suggests"`
	}
	type AS struct {
		Results []Result `json:"Results"`
	}
	type BingResp struct {
		AS AS `json:"AS"`
	}

	var bingResp BingResp
	if err := json.Unmarshal([]byte(jsonStr), &bingResp); err != nil {
		return []string{}, nil
	}

	items := []string{}
	for _, r := range bingResp.AS.Results {
		for _, s := range r.Suggests {
			items = append(items, s.Txt)
		}
	}

	return items, nil
}

type GoogleXML struct {
	XMLName             xml.Name `xml:"toplevel"`
	CompleteSuggestions []struct {
		Suggestion struct {
			Data string `xml:"data,attr"`
		} `xml:"suggestion"`
	} `xml:"CompleteSuggestion"`
}

func (svc *navPageService) GetGoogleSuggestion(q string) ([]string, common.GFError) {
	proxyURL, _ := url.Parse(env.GetServerConfig().Proxy.Url)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	url := fmt.Sprintf("http://suggestqueries.google.com/complete/search?output=toolbar&hl=zh&q=%s", q)
	resp, err := client.Get(url)
	if err != nil {
		return nil, common.NewServiceError("请求谷歌搜索建议接口出错: " + err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	type GoogleXML struct {
		XMLName             xml.Name `xml:"toplevel"`
		CompleteSuggestions []struct {
			Suggestion struct {
				Data string `xml:"data,attr"`
			} `xml:"suggestion"`
		} `xml:"CompleteSuggestion"`
	}

	var xmlResp GoogleXML
	if err := xml.Unmarshal(body, &xmlResp); err != nil {
		return []string{}, nil
	}

	items := []string{}
	for _, s := range xmlResp.CompleteSuggestions {
		items = append(items, s.Suggestion.Data)
	}
	return items, nil
}

func (svc *navPageService) GetBiliBiliSuggestion(q string) ([]string, common.GFError) {
	if q == "" {
		return []string{}, nil
	}

	url := fmt.Sprintf("https://s.search.bilibili.com/main/suggest?func=suggest&suggest_type=accurate&sub_type=tag&main_ver=v1&term=%s", q)
	resp, err := http.Get(url)
	if err != nil {
		return nil, common.NewServiceError("请求B站搜索建议接口出错: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, common.NewServiceError("读取B站响应出错: " + err.Error())
	}

	// 定义响应结构体
	type TagItem struct {
		Value string `json:"value"`
		Term  string `json:"term"`
		Name  string `json:"name"`
	}
	type BiliResp struct {
		Result struct {
			Tag []TagItem `json:"tag"`
		} `json:"result"`
	}

	var result BiliResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, common.NewServiceError("解析B站响应JSON出错: " + err.Error())
	}

	// 只返回 value
	suggestions := make([]string, len(result.Result.Tag))
	for i, item := range result.Result.Tag {
		suggestions[i] = item.Value
	}

	return suggestions, nil
}

func (svc *navPageService) GetSayingService() (string, common.GFError) {
	// 查总条数并生成一个随机值
	count, gfError := dao.GetNavPageDao().Count(models.GfnSaying{})
	if gfError != nil {
		return "", gfError
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(int(count))
	record, err := dao.GetNavPageDao().GetSayingByOrder(randomIndex)
	if err != nil || record == nil {
		return "", common.NewServiceError(fmt.Sprintf("查询金句记录失败: %v", err))
	}
	return record.Saying, nil
}

func (svc *navPageService) GetImageService(c *fiber.Ctx, t string) (string, common.GFError) {
	// 禁止浏览器/代理缓存
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	var imagePath = ""
	switch t {
	case "normal":
		imagePath = env.GetServerConfig().Resource.ImagePath
	case "resized":
		imagePath = env.GetServerConfig().Resource.ResizeImagePath
	default:
		imagePath = env.GetServerConfig().Resource.ImagePath
	}

	var imageFiles []string
	err := filepath.Walk(imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isValidImage(info.Name()) {
			imageFiles = append(imageFiles, path)
		}
		return nil
	})

	if err != nil {
		return "", common.NewServiceError(fmt.Sprintf("读取图片目录失败: %v", err))
	}

	if len(imageFiles) == 0 {
		return "", common.NewServiceError("未找到任何图片")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(imageFiles))

	return imageFiles[randomIndex], nil
}

// 检查文件是否为有效的图片格式
func isValidImage(filename string) bool {
	ext := filepath.Ext(filename)
	for _, validExt := range strings.Split(validExts, ",") {
		if strings.EqualFold(ext, validExt) {
			return true
		}
	}
	return false
}

// 增加浏览量
func addCount() {
	cs.Incr("stat-count:total") // 总访问量
	year := util.Int642String(int64(time.Now().Year()))
	month := util.Int642String(int64(time.Now().Month()))
	day := util.Int642String(int64(time.Now().Day()))
	cs.Incr("stat-count:" + year)                           // 年访问量
	cs.Incr("stat-count:" + year + "-" + month)             // 月访问量
	cs.Incr("stat-count:" + year + "-" + month + "-" + day) // 日访问量
}
