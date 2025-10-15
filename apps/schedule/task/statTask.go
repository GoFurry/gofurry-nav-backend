package task

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/GoFurry/gofurry-nav-backend/common/log"
	cs "github.com/GoFurry/gofurry-nav-backend/common/service"
	"github.com/GoFurry/gofurry-nav-backend/common/util"
)

// 更新访问量最多的几个键
func UpdateTopCountCache() {
	log.Info("StatTask 开始...")

	type regionType struct {
		Prefix   string
		CacheKey string
	}

	// 最多的 国家、省份、城市
	regions := []regionType{
		{"stat-geoip-country:", "top"},
		{"stat-geoip-province:", "top"},
		{"stat-geoip-city:", "top"},
	}

	for _, r := range regions {
		// 最多的 8 个
		topMap := getTopRegion(r.Prefix, 8)
		// 缓存一天
		if b, err := json.Marshal(topMap); err == nil {
			cs.SetExpire(r.Prefix+r.CacheKey, string(b), 24*time.Hour)
		}
	}

	log.Info("StatTask 结束...")
}

func getTopRegion(prefix string, top int) map[string]int64 {
	res := make(map[string]int64)

	// Redis 扫描所有相关 key
	ctx := context.Background()
	iter := cs.GetRedisService().Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		valStr, err := cs.GetString(key)
		if err != nil {
			continue
		}
		val, convErr := util.String2Int64(valStr)
		if convErr != nil {
			continue
		}
		// 区域名
		region := strings.TrimPrefix(key, prefix)
		res[region] = val
	}
	if err := iter.Err(); err != nil {
		log.Error("扫描redis失败:", err)
	}

	// 排序取前 top 个
	type kv struct {
		Key string
		Val int64
	}
	var kvs []kv
	for k, v := range res {
		kvs = append(kvs, kv{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Val > kvs[j].Val
	})

	topMap := make(map[string]int64)
	for i := 0; i < len(kvs) && i < top; i++ {
		topMap[kvs[i].Key] = kvs[i].Val
	}
	return topMap
}
