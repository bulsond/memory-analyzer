package info

import "testing"

func TestSystemMemoryInfo(t *testing.T) {
	t.Run("Успешное создание экземпляра SystemMemoryInfo",
		func(t *testing.T) {
			want := SystemMemoryInfo{
				TotalMemory:     16 * 1024 * 1024 * 1024, // 16 GB
				FreeMemory:      4 * 1024 * 1024 * 1024,  // 4 GB
				AvailableMemory: 6 * 1024 * 1024 * 1024,  // 6 GB
				SwapTotal:       8 * 1024 * 1024 * 1024,  // 8 GB
				SwapFree:        7 * 1024 * 1024 * 1024,  // 7 GB
			}
			got := SystemMemoryInfo{
				TotalMemory:     GBytes(16),
				FreeMemory:      GBytes(4),
				AvailableMemory: GBytes(6),
				SwapTotal:       GBytes(8),
				SwapFree:        GBytes(7),
			}
			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
}
