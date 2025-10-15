package dao

import (
	navModels "github.com/GoFurry/gofurry-nav-backend/apps/nav/navPage/models"
	"github.com/GoFurry/gofurry-nav-backend/apps/system/stat/models"
	"github.com/GoFurry/gofurry-nav-backend/common"
	"github.com/GoFurry/gofurry-nav-backend/common/abstract"
)

var newStatDao = new(statDao)

func init() {
	newStatDao.Init()
}

type statDao struct{ abstract.Dao }

func GetStatDao() *statDao { return newStatDao }

// 获取内容最多的分组
func (dao *statDao) GetGroupCount() (res []models.GroupCountVo, err common.GFError) {
	db := dao.Gm.Table(navModels.TableNameGfnSiteGroupMap).Select("group_id, COUNT(*) AS count")
	db.Group("group_id").Order("count DESC")
	db.Limit(4)
	db.Find(&res)
	if err := db.Error; err != nil {
		return res, common.NewDaoError(err.Error())
	}
	return
}
