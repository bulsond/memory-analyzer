package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bulsond/memory-analyzer/display"
	"github.com/bulsond/memory-analyzer/info"
)

type App struct {
	reader info.MemoryReader
	config display.DisplayConfig
}

// NewApp создание экземпляра приложения
func NewApp(interval, tops int) (*App, error) {
	config, err := display.NewDisplayConfig(interval, tops)
	if err != nil {
		return &App{}, err
	}

	// Выбор MemoryReader
	var reader info.MemoryReader
	switch runtime.GOOS {
	case "darwin":
		reader = &info.DarwinMemoryReader{}
	case "linux":
		reader = &info.LinuxMemoryReader{}
	default:
		msg := fmt.Sprintf("Unsupported operating system: %s\n", runtime.GOOS)
		return &App{}, errors.New(msg)
	}

	return &App{
		reader: reader,
		config: config,
	}, nil
}

func (a *App) Run() {
	// Настройка обработки сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Создание ticker
	ticker := time.NewTicker(a.config.UpdateInterval.Duration())
	defer ticker.Stop()

	fmt.Printf("Starting Memory Analyzer on %s\n", runtime.GOOS)

	// Основной цикл
	for {
		select {
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal. Exiting...")
			return
		case <-ticker.C:
			// Получение системной информации
			sysInfo, err := a.reader.ReadSystemMemory()
			if err != nil {
				fmt.Printf("Error reading system memory: %v\n", err)
				continue
			}

			// Получение списка процессов
			pids, err := a.reader.GetProcessList()
			if err != nil {
				fmt.Printf("Error getting process list: %v\n", err)
				continue
			}

			// Сбор информации о процессах
			processes := a.collectProcesses(pids)

			// Отображение информационной панели
			display.DisplayDashboard(sysInfo, processes, a.config)
		}
	}
}

// collectProcesses сбор информации о процессах
func (a *App) collectProcesses(pids []info.PID) []info.ProcessInfo {
	var processes []info.ProcessInfo
	for _, pid := range pids {
		if mem, err := a.reader.ReadProcessMemory(pid); err == nil {
			cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "comm=")
			if output, err := cmd.Output(); err == nil {
				name := strings.TrimSpace(string(output))
				process, _ := info.NewProcessInfo(pid.Int(), name, mem.Uint64())
				processes = append(processes, process)
			}
		}
	}

	return processes
}
