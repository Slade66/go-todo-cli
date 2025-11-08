package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type Todos []Todo

func (todos *Todos) add(title string) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   time.Now(),
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) delete(index int) {
	t := *todos
	*todos = append(t[:index], t[index+1:]...)
}

func (todos *Todos) toggle(index int) {
	t := *todos
	if t[index].Completed {
		t[index].Completed = false
		t[index].CompletedAt = nil
	} else {
		t[index].Completed = true
		now := time.Now()
		t[index].CompletedAt = &now
	}
}

func (todos *Todos) edit(index int, title string) {
	t := *todos
	t[index].Title = title
}

func (todos *Todos) print() {
	fmt.Print("\033[H\033[2J") // ANSI 转义序列
	t := table.New(os.Stdout)
	t.SetHeaders("#", "Title", "Completed", "Created At", "Completed At")
	for index, todo := range *todos {
		completedAt := ""
		if todo.Completed {
			completedAt = todo.CompletedAt.Format(time.DateTime)
		}
		t.AddRow(strconv.Itoa(index),
			todo.Title,
			strconv.FormatBool(todo.Completed),
			todo.CreatedAt.Format(time.DateTime),
			completedAt,
		)
	}
	t.Render()
}
