package info

import "testing"

func TestProcessInfo(t *testing.T) {
	t.Run("Успешное создание экземпляра ProcessInfo",
		func(t *testing.T) {
			want := ProcessInfo{
				PID:         PID(1234),
				Name:        ProcessName("chrome"),
				MemoryUsage: MBytes(512),
			}
			got, err := NewProcessInfo(1234, "chrome", 512*1024*1024)
			if err != nil {
				t.Errorf("Ошибка создания ProcessInfo: %v", err)
			}
			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
}
