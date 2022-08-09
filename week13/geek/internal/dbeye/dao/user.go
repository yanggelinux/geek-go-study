package dao

import (
	"fmt"
	"geek/internal/dbeye/stroe/mysql"
	"geek/internal/pkg/ce"
	"geek/pkg/util"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

type userRes struct {
	ID          uint32
	UserName    string
	FullName    string
	UpdatedTime uint64
	CreatedTime uint64
}

func (us *UserDao) GetUserList(userName *string, page, pageSize *int) (error, int, map[string]interface{}) {

	type res map[string]interface{}
	var (
		condition  string
		args       []interface{}
		resList    []res
		totalCount int64
		results    []userRes
	)
	code := ce.SUCCESS
	data := make(map[string]interface{})

	condition += "user.is_deleted = ?"
	args = append(args, 0)

	if userName != nil {
		condition += " AND user.user_name Like ?"
		args = append(args, fmt.Sprintf("%%%s%%", *userName))
	}

	//数据库操作
	err := mysql.DB.Table("user").Select("user.id").
		Where(condition, args...).Count(&totalCount).Error
	if err != nil {
		code = ce.ERROR_DB_TABLE_QUERY_FAILED
		err = fmt.Errorf("数据库查询表记录失败:%w", err)
		return err, code, data
	}
	if page != nil && pageSize != nil {
		err = mysql.DB.Table("user").Select(
			"user.id,user.user_name,user.full_name,user.updated_time,user.created_time").
			Where(condition, args...).Order("created_time desc").
			Offset((*page - 1) * *pageSize).Limit(*pageSize).Scan(&results).Error
	} else {
		err = mysql.DB.Table("user").Select(
			"user.id,user.user_name,user.full_name,user.updated_time,user.created_time").
			Where(condition, args...).Order("created_time desc").Scan(&results).Error
	}

	if err != nil {
		code = ce.ERROR_DB_TABLE_QUERY_FAILED
		err = fmt.Errorf("数据库查询表记录失败:%w", err)
		return err, code, data
	}
	//
	if pageSize != nil {
		resList = make([]res, 0, *pageSize)
	} else {
		resList = make([]res, 0, totalCount)
	}
	for _, result := range results {
		res := make(map[string]interface{})
		res["id"] = result.ID
		res["userName"] = result.UserName
		res["fullName"] = result.FullName
		res["updatedTime"] = util.FormatTimeToString(result.UpdatedTime)
		res["createdTime"] = util.FormatTimeToString(result.CreatedTime)
		resList = append(resList, res)
	}
	data["totalCount"] = totalCount
	data["resList"] = resList
	return err, code, data
}
