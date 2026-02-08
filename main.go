package main

import (
	"github.com/bulsond/memory-analyzer/display"
	"github.com/bulsond/memory-analyzer/info"
)

func main() {
	// fmt.Println("Привет")

	// Создаём тестовый набор процессов
	// processes := []info.ProcessInfo{
	// 	{
	// 		PID:         1234,
	// 		Name:        "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	// 		MemoryUsage: 1024 * 1024 * 1024, // 1 GB
	// 	},
	// 	{
	// 		PID:         5678,
	// 		Name:        "/Applications/Visual Studio Code.app/Contents/MacOS/Electron",
	// 		MemoryUsage: 512 * 1024 * 1024, // 512 MB
	// 	},
	// 	{
	// 		PID:         9101,
	// 		Name:        "/usr/bin/ssh-agent",
	// 		MemoryUsage: 2 * 1024 * 1024, // 2 MB
	// 	},
	// }

	// // Форматируем и выводим таблицу
	// fmt.Println("Process List:")
	// fmt.Print(display.FormatTable(processes))

	// Создаём тестовую информацию о системной памяти
	sysInfo := info.SystemMemoryInfo{
		TotalMemory:     16 * 1024 * 1024 * 1024, // 16 GB
		FreeMemory:      4 * 1024 * 1024 * 1024,  // 4 GB
		AvailableMemory: 6 * 1024 * 1024 * 1024,  // 6 GB
		SwapTotal:       8 * 1024 * 1024 * 1024,  // 8 GB
		SwapFree:        7 * 1024 * 1024 * 1024,  // 7 GB
	}

	// Форматируем и выводим статистику
	// fmt.Print(sysInfo)

	processes := []info.ProcessInfo{
		{
			PID:         1234,
			Name:        "chrome",
			MemoryUsage: 1024 * 1024 * 1024, // 1 GB
		},
		{
			PID:         5678,
			Name:        "vscode",
			MemoryUsage: 512 * 1024 * 1024, // 512 MB
		},
		{
			PID:         7890,
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

	config, err := display.NewDisplayConfig(3, 10)
	if err != nil {
		panic(err)
	}

	//Отображаем информационную панель
	display.DisplayDashboard(sysInfo, processes, config)
}
