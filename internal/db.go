package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

/**
  @author: XingGao
  @date: 2024/5/24
**/

type DBModel struct {
	DB *sql.DB
}

type DBConfig struct {
	username string
	password string
	host     string
	port     string
	database string
}

func NewDBModel(host, username, password, port, database string) *DBModel {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("无法连接到数据库: %v", err)
		return nil
	}
	return &DBModel{DB: db}
}

// Ping ping测试
func (d *DBModel) Ping() error {
	return d.DB.Ping()
}

// Close 关闭连接
func (d *DBModel) Close() {
	d.DB.Close()
}

// executeSQL 执行sql语句
func (d *DBModel) executeSQL(sql string) error {
	// 按分号分割 SQL 语句
	sqlStmts := strings.Split(sql, ";")

	for _, stmt := range sqlStmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// 执行 SQL 语句
		if _, err := d.DB.Exec(stmt); err != nil {
			return fmt.Errorf("执行 SQL 语句失败: %v", err)
		}
	}
	return nil
}
