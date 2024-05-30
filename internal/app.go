package internal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/xinggaoya/qwen-sdk/qwen"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// TestDB 测试连接
func (a *App) TestDB(host, port, username, password, database string) bool {
	db := NewDBModel(host, username, password, port, database)
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// GetAllDatabases 获取所有数据库的名称
func (a *App) GetAllDatabases(host, port, username, password string) ([]string, error) {
	db := NewDBModel(host, username, password, port, "")
	rows, err := db.DB.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("查询数据库列表失败: %v", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err = rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("读取数据库名称失败: %v", err)
		}
		databases = append(databases, dbName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历数据库名称失败: %v", err)
	}

	return databases, nil
}

// GetAllTables 获取所有表
func (a *App) GetAllTables(host, port, username, password, database string) ([]string, error) {
	db := NewDBModel(host, username, password, port, database)
	defer db.Close()

	rows, err := db.DB.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("查询表列表失败: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("读取表名称失败: %v", err)
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历表名称失败: %v", err)
	}

	return tables, nil
}

type TableColumn struct {
	Field   string         // 列名
	Type    string         // 类型
	Null    string         // 是否为空
	Key     string         // 主键
	Default sql.NullString // 默认值
	Extra   string         // 其他信息
}

// GetTableStructure 获取表结构
func (a *App) GetTableStructure(host, port, username, password, database, table string) ([]TableColumn, error) {
	db := NewDBModel(host, username, password, port, database)
	defer db.Close()

	sprintf := fmt.Sprintf("DESC %s", table)
	rows, err := db.DB.Query(sprintf)
	if err != nil {
		return nil, fmt.Errorf("查询表结构失败: %v", err)
	}
	defer rows.Close()
	var columns []TableColumn
	for rows.Next() {
		var column TableColumn
		if err = rows.Scan(&column.Field, &column.Type, &column.Null, &column.Key, &column.Default, &column.Extra); err != nil {
			return nil, fmt.Errorf("读取列名称失败: %v", err)
		}
		columns = append(columns, column)

	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历列名称失败: %v", err)
	}

	return columns, nil
}

// GetAIAnswer 获取回答
func (a *App) GetAIAnswer(info, question string) string {
	apiKey := "sk-3437acc67dbf452f9b2f0f2a49640518"
	client := qwen.NewWithDefaultChat(apiKey)
	// 定义一条消息对话的历史记录
	messages := []qwen.Messages{
		{Role: "system", Content: "你现在是SQL大师，下面将会给你发送MySQL数据表结构，你结合进行回答问题。言简意赅，只回答问题，不要过多解释！"},
		{Role: qwen.ChatUser, Content: info},
		{Role: qwen.ChatUser, Content: question},
	}

	// 获取AI对消息的回复
	resp, err := client.GetAIReply(messages)
	if err != nil {
		fmt.Printf("获取AI回复失败：%v\n", err)
		return "Answer: " + question
	}

	return resp.Output.Text
}
