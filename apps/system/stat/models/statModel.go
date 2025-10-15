package models

type ViewsCountVo struct {
	Total      int64    `json:"total"`
	YearCount  int64    `json:"year_count"`
	MonthCount int64    `json:"month_count"`
	Date       []string `json:"date"`
	Count      []int64  `json:"count"`
}

type GroupCountVo struct {
	GroupID string `gorm:"column:group_id" json:"group_id"`
	Count   int64  `gorm:"column:count" json:"count"`
}

type RegionCountVo struct {
	RegionMap map[string]int64 `gorm:"column:region_map" json:"region_map"`
}
