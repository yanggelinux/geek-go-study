package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const DB_USER = "root"
const DB_PASSWORD = ""
const DB_HOST = ""
const DB_PORT = "3306"
const DB_NAME = "myapp"

var ErrNoRows = sql.ErrNoRows

type MySQLDB struct {
	db          *sql.DB
	DSN         string
	Active      int           // pool
	Idle        int           // pool
	IdleTimeout time.Duration // connect max life time.
}

func (m *MySQLDB) MysqlOpen() error {
	db, err := sql.Open("mysql", m.DSN)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(m.Active)
	db.SetMaxIdleConns(m.Idle)
	db.SetConnMaxLifetime(time.Minute * m.IdleTimeout)
	m.db = db
	return nil
}

func (m *MySQLDB) MysqlClose() {
	err := m.db.Close()
	if err != nil {
		log.Println("close mysql error:", err)
		return
	}
}

//select 查询数据

func (m *MySQLDB) Query(sql string) (*sql.Rows, error) {
	rows, err := m.db.Query(sql)
	return rows, err

}

func (m *MySQLDB) QueryRow(sql string) *sql.Row {
	row := m.db.QueryRow(sql)
	return row

}

func NewMySQLDB(DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME string) *MySQLDB {
	var MySQLDB *MySQLDB = new(MySQLDB)
	MySQLDB.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	MySQLDB.Active = 1000
	MySQLDB.Idle = 500
	MySQLDB.IdleTimeout = time.Duration(10)
	return MySQLDB
}

type UserDao struct {
	db *MySQLDB
}

func NewUserDao(db *MySQLDB) *UserDao {
	return &UserDao{db: db}
}

//查询用户数据库
func (us *UserDao) getUserInfo(s string) (error, map[string]string) {

	result := make(map[string]string)

	row := us.db.QueryRow(s)

	var user_name, full_name string
	err := row.Scan(&user_name, &full_name)
	if err != nil {
		//包装err信息抛给上层
		err := fmt.Errorf("查询用户信息失败,%w", err)
		return err, result
	}
	result["user_name"] = user_name
	result["full_name"] = full_name
	return nil, result

}

type UserService struct {
	db *MySQLDB
}

func NewUserService(db *MySQLDB) *UserService {
	return &UserService{db: db}
}

//查询用户服务
func (us *UserService) query() map[string]string {

	result := make(map[string]string)

	s := "select user_name,full_name from user where id = 1"
	userDao := NewUserDao(us.db)
	err, result := userDao.getUserInfo(s)
	// 判断err类型是否是ErrNoRows
	if errors.Is(err, ErrNoRows) {
		log.Println("数据不存在", err)
		return result
	}
	return result

}

func main() {

	db := NewMySQLDB(DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	err := db.MysqlOpen()
	if err != nil {
		log.Fatal("获取数据库链接失败")
	}
	defer db.MysqlClose()
	userSvc := NewUserService(db)
	result := userSvc.query()
	log.Println("result:", result)

}

//我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么？
// 我认为错误应该在dao层包装后抛给上层，在dao的上层service层处理错误，service层处理业务逻辑后，将处理信息返回给调用者
