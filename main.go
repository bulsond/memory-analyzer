package main

import (
	"fmt"

	"github.com/bulsond/memory-analyzer/display"
	"github.com/bulsond/memory-analyzer/info"
)

func main() {
	fmt.Println("Привет")

	// Создаём тестовый набор процессов
	processes := []info.ProcessInfo{
		{
			PID:         1234,
			Name:        "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			MemoryUsage: 1024 * 1024 * 1024, // 1 GB
		},
		{
			PID:         5678,
			Name:        "/Applications/Visual Studio Code.app/Contents/MacOS/Electron",
			MemoryUsage: 512 * 1024 * 1024, // 512 MB
		},
		{
			PID:         9101,
			Name:        "/usr/bin/ssh-agent",
			MemoryUsage: 2 * 1024 * 1024, // 2 MB
		},
	}

	// Форматируем и выводим таблицу
	fmt.Println("Process List:")
	fmt.Print(display.FormatTable(processes))
}
