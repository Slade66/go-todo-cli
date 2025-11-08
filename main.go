package main

func main() {
	var todos Todos
	// todos.add("吃一个榴莲大福")
	// todos.add("喝一口豆浆")
	
	storage := NewStorage("todos.json")
	storage.Load(&todos)
	
	todos.print()
	
	todos.delete(0)
	todos.edit(0, "喝完豆浆")
	todos.toggle(0)

	todos.print()
	storage.Save(todos)
}
