package info

import (
	"strings"
	"testing"
)

func TestSystemMemoryInfo(t *testing.T) {
	smi := SystemMemoryInfo{
		TotalMemory:     16 * 1024 * 1024 * 1024, // 16 GB
		FreeMemory:      4 * 1024 * 1024 * 1024,  // 4 GB
		AvailableMemory: 6 * 1024 * 1024 * 1024,  // 6 GB
		SwapTotal:       8 * 1024 * 1024 * 1024,  // 8 GB
		SwapFree:        7 * 1024 * 1024 * 1024,  // 7 GB
	}
	t.Run("Успешное создание экземпляра SystemMemoryInfo",
		func(t *testing.T) {
			got := SystemMemoryInfo{
				TotalMemory:     GBytes(16),
				FreeMemory:      GBytes(4),
				AvailableMemory: GBytes(6),
				SwapTotal:       GBytes(8),
				SwapFree:        GBytes(7),
			}
			if got != smi {
				t.Errorf("got %q want %q", got, smi)
			}
		})

	t.Run("Правильное отображение экземпляра SystemMemoryInfo",
		func(t *testing.T) {
			got := smi.String()
			if !strings.Contains(got, "System Memory:") {
				t.Error("Не найдена строка System Memory")
			}
			if !strings.Contains(got, "Total:") {
				t.Error("Не найдена строка Total")
			}
			if !strings.Contains(got, "Used:") {
				t.Error("Не найдена строка Used")
			}
			if !strings.Contains(got, "Available:") {
				t.Error("Не найдена строка Available")
			}
			if !strings.Contains(got, "Swap Used:") {
				t.Error("Не найдена строка Swap Used")
			}
		})
}
