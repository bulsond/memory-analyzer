package display

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bulsond/memory-analyzer/info"
)

// GetShortProcessName форматирование названия процесса
func GetShortProcessName(fullName string) string {
	// 1. Извлекаем базовое имя из полного пути
	base := filepath.Base(fullName)

	// 2. Удаляем стандартные суффиксы
	suffixes := [4]string{"-helper (Renderer)", " Helper (Renderer)", " Helper", ".app"}
	for _, suffix := range suffixes {
		if idx := strings.Index(base, suffix); idx != -1 {
			base = base[:idx]
		}
	}

	// 3. Если многословное название берем последнее слово
	fields := strings.Fields(base)
	length := len(fields)
	if length > 1 {
		base = fields[length-1]
	}

	// 4. Ограничиваем длину до 15 символов с добавлением "..."
	const maxLen = 15
	const suffixLen = 3
	if len(base) > maxLen {
		base = base[:(maxLen-suffixLen)] + "..."
	}

	return base
}

// FormatTable формирование таблицы процессов
func FormatTable(processes []info.ProcessInfo) string {
	var sb strings.Builder
	// Заголовок
	sb.WriteString("PID      NAME            MEMORY\n")
	// Разделитель
	sb.WriteString("--------------------------------\n")

	// Вывод процессов
	for _, proc := range processes {
		shortName := GetShortProcessName(string(proc.Name))
		memStr := proc.MemoryUsage.String()
		// Форматирование: PID (8, лево), Name (15), Memory (10)
		line := fmt.Sprintf("%-8d %-15s %10s\n", proc.PID, shortName, memStr)
		sb.WriteString(line)
	}

	return sb.String()
}

func DisplayDashboard(
	stats info.SystemMemoryInfo,
	processes []info.ProcessInfo,
	config DisplayConfig) {
	// Очистка экрана (работает на Unix/Linux/macOS и Windows)
	fmt.Print("\033[2J\033[H")

	// Заголовок
	fmt.Println()
	fmt.Println("=== Memory Analyzer ===")
	fmt.Println()

	// О памяти
	fmt.Println(stats)

	// Подзаголовок
	fmt.Println("Top Memory Processes:")

	// О процессах
	fmt.Println(FormatTable(processes))

	// Временная метка
	fmt.Printf("Updated: %s\n",
		time.Now().Format("2006-01-02 15:04:05"))

	// Подсказка о выходе
	fmt.Println("Press Ctrl+C to exit")
}
