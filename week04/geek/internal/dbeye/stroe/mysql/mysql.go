package mysql

import (
	"database/sql"
	"fmt"
	"geek/global"
	"geek/pkg/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	DB    *gorm.DB
	sqlDB *sql.DB
	err   error
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type Model struct {
	IsDeleted   uint32 `gorm:"default:0" json:"is_deleted"`
	DeletedTime uint64 `json:"deleted_time"` //自动填充时间到毫秒
	UpdatedTime uint64 `json:"updated_time"`
	CreatedTime uint64 `json:"created_time"`
}

func updateAndCreateTimeForCreateCallback(DB *gorm.DB) {
	if DB.Statement.Schema != nil {
		nowTime := uint64(time.Now().UnixNano() / 1e6)
		createTimeField := DB.Statement.Schema.LookUpField("CreatedTime")
		if createTimeField != nil {
			DB.Statement.SetColumn("CreatedTime", nowTime, true)
		}

		updateTimeField := DB.Statement.Schema.LookUpField("UpdatedTime")
		if updateTimeField != nil {
			DB.Statement.SetColumn("UpdatedTime", nowTime, true)

		}
	}
}

func updateTimeForUpdateCallback(DB *gorm.DB) {
	if DB.Statement.Schema != nil {
		nowTime := uint64(time.Now().UnixNano() / 1e6)
		updateTimeField := DB.Statement.Schema.LookUpField("UpdatedTime")
		if updateTimeField != nil {
			DB.Statement.SetColumn("UpdatedTime", nowTime, true)

		}

	}
}

func deleteTimeForUpdateCallback(DB *gorm.DB) {
	if DB.Statement.Schema != nil {
		nowTime := uint64(time.Now().UnixNano() / 1e6)
		deleteTimeField := DB.Statement.Schema.LookUpField("DeletedTime")
		if deleteTimeField != nil {
			DB.Statement.SetColumn("DeletedTime", nowTime, true)

		}

	}
}

func SetupModel() {
	DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		global.MySQLSetting.Username,
		global.MySQLSetting.Password,
		global.MySQLSetting.Host,
		global.MySQLSetting.Port,
		global.MySQLSetting.DBName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Logger.Error("connect database error:", zap.Error(err))
		return
	}

	sqlDB, err = DB.DB()
	if err != nil {
		log.Logger.Error("get db.DB() error:", zap.Error(err))
		return
	}
	// callback
	err = DB.Callback().Create().Before("gorm:create").Register("createUpdateTime", updateAndCreateTimeForCreateCallback)
	if err != nil {
		log.Logger.Error("updateAndCreateTimeForCreateCallback Register error:", zap.Error(err))
	}
	err = DB.Callback().Update().Before("gorm:update").Register("updateTime", updateTimeForUpdateCallback)
	if err != nil {
		log.Logger.Error("updateTimeForUpdateCallback Register error:", zap.Error(err))
	}
	err = DB.Callback().Delete().Before("gorm:delete").Register("deleteTime", deleteTimeForUpdateCallback)
	if err != nil {
		log.Logger.Error("deleteTimeForDeleteCallback Register error:", zap.Error(err))
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	//创建表 建议手动创建表
	//err = DB.AutoMigrate(&User{}, &Role{})
	//if err != nil {
	//	log.Logger.Error("create db table error:", zap.Error(err))
	//	return
	//}
}
