package display

import (
	"path/filepath"
	"strings"
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
// func FormatTable(processes []ProcessInfo) string {
// 	return ""
// }
