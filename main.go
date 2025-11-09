package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var todos Todos
var storage *Storage
var storagePath string

func init() {
	flag.StringVar(&storagePath, "storage-path", "todos.json", "数据文件的路径")
	flag.Parse()
}

func main() {
	// 1. 初始化日志器
	if err := initLogger(); err != nil {
		// 若失败：向标准错误输出一条清晰的启动失败信息，然后以非零退出码退出程序。
		fmt.Fprintf(os.Stderr, "Logger initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer zap.L().Sync() // 确保程序退出前，把日志刷新到文件。
	zap.L().Debug("Logger initialized")

	// 2. 加载数据
	loadTodos()

	for {
		todos.print()
		showMenu()
		choice, err := getUserChoice()
		if err != nil {
			zap.L().Warn("Invalid user input.", zap.Int("choice", *choice), zap.Error(err))
			fmt.Println("输入有误，请重试。")
			continue
		}
		shouldExit := Execute(*choice)
		if shouldExit {
			return
		}
	}

}

// initLogger 用于初始化日志器。
// 它是一个 “只负责创建并返回结果” 的函数。
// 如果初始化失败，返回一个包装后的错误，不在这个函数里 panic 或退出；成功时替换全局 logger 后返回 nil。
func initLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("无法初始化日志器: %w", err)
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func Execute(choice int) bool {
	switch choice {
	case 1:
		addTodo()
	case 2:
		deleteTodo()
	case 3:
		toggleTodo()
	case 4:
		editTodo()
	case 5:
		saveTodos()
	case 6:
		fmt.Println("Bye.")
		zap.L().Debug("Exiting application.")
		return true
	default:
		zap.L().Warn("Invalid menu choice.", zap.Int("choice", choice))
		fmt.Println("无效的选择，请重新输入")
	}
	return false
}

func editTodo() {
	var index int
	fmt.Print("请输入待办事项索引: ")
	if _, err := fmt.Scan(&index); err != nil {
		zap.L().Warn("Invalid user input.", zap.Int("index", index), zap.Error(err))
		return
	}
	var title string
	fmt.Print("请输入待办事项标题: ")
	if _, err := fmt.Scan(&title); err != nil {
		zap.L().Warn("Invalid user input.", zap.String("title", title), zap.Error(err))
		return
	}
	todos.edit(index, title)
	zap.L().Debug("Todo edited.", zap.Int("index", index), zap.String("title", title))
}

func saveTodos() {
	if err := storage.Save(todos); err != nil {
		zap.L().Error("Failed to save todos", zap.String("file", storage.FileName), zap.Error(err))
		fmt.Println("保存失败，请重试。")
		return
	}
	zap.L().Debug("Todos saved", zap.String("file", storage.FileName), zap.Int("count", len(todos)))
}

func toggleTodo() {
	var index int
	fmt.Print("请输入待办事项索引: ")
	if _, err := fmt.Scan(&index); err != nil {
		zap.L().Warn("Invalid user input.", zap.Int("index", index), zap.Error(err))
		return
	}
	todos.toggle(index)
	zap.L().Debug("Todo toggled.", zap.Int("index", index))
}

func deleteTodo() {
	var index int
	fmt.Print("请输入待办事项索引: ")
	if _, err := fmt.Scan(&index); err != nil {
		zap.L().Warn("Invalid user input.", zap.Int("index", index), zap.Error(err))
		return
	}
	todos.delete(index)
	zap.L().Debug("Todo deleted.", zap.Int("index", index))
}

func addTodo() {
	var title string
	fmt.Print("请输入待办事项标题: ")
	if _, err := fmt.Scan(&title); err != nil {
		zap.L().Warn("Invalid user input.", zap.String("title", title), zap.Error(err))
		return
	}
	todos.add(title)
	zap.L().Debug("Todo added.", zap.String("title", title))
}

func getUserChoice() (*int, error) {
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		return nil, fmt.Errorf("输入错误: %w", err)
	}
	return &choice, nil
}

func showMenu() {
	fmt.Println("1. 添加待办事项")
	fmt.Println("2. 删除待办事项")
	fmt.Println("3. 完成待办事项")
	fmt.Println("4. 编辑待办事项")
	fmt.Println("5. 保存待办事项")
	fmt.Println("6. 退出")
	fmt.Print("请输入您的选择: ")
}

// loadTodos 用于从文件加载待办事项。
func loadTodos() {
	storage = NewStorage(storagePath)
	if err := storage.Load(&todos); err != nil {
		todos = Todos{}
		zap.L().Warn("Failed to load todos, using empty todos.", zap.String("file", storage.FileName), zap.Error(err))
		return
	}
	zap.L().Debug("Todos loaded.", zap.Int("count", len(todos)), zap.String("file", storage.FileName))
}
