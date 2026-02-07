package display

import (
	"errors"
	"time"
)

// Interval представляет период обновления в виде времени.
// Используется для управления частотой обновления интерфейса или данных.
type Interval time.Duration

// CountProcesses определяет количество процессов, которые следует отображать.
// Используется для ограничения числа выводимых процессов в топ-списке.
type CountProcesses int

// DisplayConfig содержит конфигурацию отображения в программе.
// Управляет интервалом обновления и количеством отображаемых процессов.
type DisplayConfig struct {
	// UpdateInterval задает период, через который происходит обновление данных.
	// Значение должно быть положительным. Использует тип Interval для семантической ясности.
	UpdateInterval Interval

	// TopProcesses определяет максимальное количество процессов, отображаемых в топ-листе.
	// Должно быть неотрицательным.
	TopProcesses CountProcesses
}

// NewDisplayConfig создание экземпляра DisplayConfig
func NewDisplayConfig(seconds, count int) (DisplayConfig, error) {
	if seconds <= 0 {
		return DisplayConfig{},
			errors.New("Количество секунд не может быть нулевым или меньше нуля")
	}
	if count <= 0 {
		return DisplayConfig{},
			errors.New("Число процессов не может быть нулевым или меньше нуля")
	}

	duration := time.Duration(seconds) * time.Second

	return DisplayConfig{
		UpdateInterval: Interval(duration),
		TopProcesses:   CountProcesses(count),
	}, nil
}
