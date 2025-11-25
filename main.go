package main

import (
	"errors"

	"github.com/GoMateoGo/inp/pkg/logger"
)

func main() {

	// 打印
	userID := 10086
	userName := "Mateo"
	err := errors.New("connection timeout")

	// 2. 直接使用 (Printf 风格)

	// INFO (绿色)
	logger.Info("程序启动成功，监听端口: %d", 8080)

	logger.Warn("这里测试了一个警告:%d", 9999)

	// DEBUG (紫色)
	logger.Debug("正在处理用户: %s (ID: %d)", userName, userID)

	// ERROR (红色)
	// 完美满足您的需求：直接用占位符
	logger.Error("数据库查询失败，用户: %s, 错误详情: %v", userName, err)
}
