package display

import (
	"testing"
	"time"
)

func TestDisplayConfig(t *testing.T) {
	t.Run("Успешное создание экземпляра DisplayConfig", func(t *testing.T) {
		want := DisplayConfig{
			UpdateInterval: Interval(3 * time.Second),
			TopProcesses:   CountProcesses(10),
		}
		got, err := NewDisplayConfig(3, 10)
		if err != nil {
			t.Errorf("Неверные аргументы для создания конфигурации: %q", err)
		}
		if want != got {
			t.Errorf("Хотели создать %v, а создали % v", want, got)
		}
	})
}
