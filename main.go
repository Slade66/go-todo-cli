package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var todos Todos
var storage *Storage

func main() {
	// 1. 初始化日志器
	if err := initLogger(); err != nil {
		// 若失败：向标准错误输出一条清晰的启动失败信息，然后以非零退出码退出程序。
		fmt.Fprintf(os.Stderr, "Logger initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer zap.L().Sync() // 确保程序退出前，把日志刷新到文件。
	zap.L().Info("Logger initialized")
	
	// 2. 加载数据
	loadTodos()

	for {
		todos.print()
		showMenu()
		choice, err := getUserChoice()
		if err != nil {
			fmt.Println(err)
			continue
		}
		shouldExit := Execute(choice)
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
		return true
	default:
		fmt.Println("无效的选择，请重新输入")
	}
	return false
}

func editTodo() {
	var index int
	fmt.Print("请输入待办事项索引: ")
	fmt.Scan(&index)
	var title string
	fmt.Print("请输入待办事项标题: ")
	fmt.Scan(&title)
	todos.edit(index, title)
}

func saveTodos() {
	storage.Save(todos)
}

func toggleTodo() {
	var index int
	fmt.Print("请输入待办事项索引: ")
	fmt.Scan(&index)
	todos.toggle(index)
}

func deleteTodo() {
	var index int
	fmt.Print("请输入待办事项索引: ")
	fmt.Scan(&index)
	todos.delete(index)
}

func addTodo() {
	var title string
	fmt.Print("请输入待办事项标题: ")
	fmt.Scan(&title)
	todos.add(title)
}

func getUserChoice() (int, error) {
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		return choice, fmt.Errorf("输入错误: %w", err)
	}
	return choice, nil
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

func loadTodos() {
	storagePath := flag.String("storage-path", "todos.json", "数据文件的路径")
	flag.Parse()
	storage = NewStorage(*storagePath)
	storage.Load(&todos)
}
