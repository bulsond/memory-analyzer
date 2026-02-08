package info

import "testing"

// тестовая реализация интерфейса MemoryReader
type FakeMemoryReader struct{}

func (f *FakeMemoryReader) ReadSystemMemory() (SystemMemoryInfo, error) {
	return SystemMemoryInfo{
		TotalMemory:     GBytes(8),
		FreeMemory:      GBytes(4),
		AvailableMemory: GBytes(4),
		SwapTotal:       GBytes(2),
		SwapFree:        GBytes(2),
	}, nil
}

func (f *FakeMemoryReader) GetProcessList() ([]PID, error) {
	return []PID{1, 100, 200, 300}, nil
}

func (f *FakeMemoryReader) ReadProcessMemory(pid PID) (Bytes, error) {
	return MBytes(256), nil
}

func TestFakeMemoryReader(t *testing.T) {
	var reader MemoryReader = &FakeMemoryReader{}
	t.Run("Успешное получение экземпляра SystemMemoryInfo", func(t *testing.T) {
		gb8 := Bytes(8 * 1024 * 1024 * 1024)
		gb4 := Bytes(4 * 1024 * 1024 * 1024)
		gb2 := Bytes(2 * 1024 * 1024 * 1024)

		sysInfo, err := reader.ReadSystemMemory()
		if err != nil {
			t.Errorf("Ошибка получения экземпляра SystemMemoryInfo: %q", err)
		}

		if sysInfo.TotalMemory != gb8 {
			t.Errorf("total got: %d want: %d", sysInfo.TotalMemory, gb8)
		}
		if sysInfo.FreeMemory != gb4 {
			t.Errorf("free got: %d want: %d", sysInfo.FreeMemory, gb4)
		}
		if sysInfo.AvailableMemory != gb4 {
			t.Errorf("available got: %d want: %d", sysInfo.AvailableMemory, gb4)
		}
		if sysInfo.SwapTotal != gb2 {
			t.Errorf("swap total got: %d want: %d", sysInfo.SwapTotal, gb2)
		}
		if sysInfo.SwapFree != gb2 {
			t.Errorf("swap free got: %d want: %d", sysInfo.SwapFree, gb2)
		}
	})

	t.Run("Успешное получение списка процессов", func(t *testing.T) {
		pids, err := reader.GetProcessList()
		if err != nil {
			t.Errorf("Ошибка получения списка процессов: %q", err)
		}
		if len(pids) == 0 {
			t.Error("Получен пустой список процессов")
		}
	})

	t.Run("Успешное чтение памяти процесса", func(t *testing.T) {
		memory, err := reader.ReadProcessMemory(PID(0))
		if err != nil {
			t.Errorf("Ошибка получения памяти процесса: %q", err)
		}
		if memory != 256*MB {
			t.Error("Получен неверный размер памяти процесса")
		}
	})
}

func TestDarwinMemoryReader(t *testing.T) {
	var reader MemoryReader = &DarwinMemoryReader{}
	t.Run("Успешное получение списка процессов", func(t *testing.T) {
		pids, err := reader.GetProcessList()
		if err != nil {
			t.Errorf("Ошибка получения списка процессов: %q", err)
		}
		if len(pids) == 0 {
			t.Error("Получен пустой список процессов")
		}
	})

	t.Run("Успешное чтение памяти процесса", func(t *testing.T) {
		pid := PID(1) // PID 1 (init/systemd)
		memory, err := reader.ReadProcessMemory(pid)
		if err != nil {
			t.Errorf("Ошибка получения памяти процесса: %q", err)
		}
		if memory == 0 {
			t.Error("Получен неверный размер памяти процесса")
		}
	})
}

func TestLinuxMemoryReader(t *testing.T) {
	var reader MemoryReader = &LinuxMemoryReader{}
	t.Run("Успешное получение списка процессов", func(t *testing.T) {
		pids, err := reader.GetProcessList()
		if err != nil {
			t.Errorf("Ошибка получения списка процессов: %q", err)
		}
		if len(pids) == 0 {
			t.Error("Получен пустой список процессов")
		}
	})

	t.Run("Успешное чтение памяти процесса", func(t *testing.T) {
		pid := PID(1) // PID 1 (init/systemd)
		memory, err := reader.ReadProcessMemory(pid)
		if err != nil {
			t.Errorf("Ошибка получения памяти процесса: %q", err)
		}
		if memory == 0 {
			t.Error("Получен неверный размер памяти процесса")
		}
	})

	t.Run("Успешное получение экземпляра SystemMemoryInfo", func(t *testing.T) {
		sysInfo, err := reader.ReadSystemMemory()
		if err != nil {
			t.Errorf("Ошибка получения экземпляра SystemMemoryInfo: %q", err)
		}

		if sysInfo.TotalMemory == 0 {
			t.Errorf("не получено значение TotalMemory")
		}
		if sysInfo.FreeMemory == 0 {
			t.Errorf("не получено значение FreeMemory")
		}
		if sysInfo.AvailableMemory == 0 {
			t.Errorf("не получено значение AvailableMemory")
		}
		if sysInfo.SwapTotal == 0 {
			t.Errorf("не получено значение SwapTotal")
		}
		if sysInfo.SwapFree == 0 {
			t.Errorf("не получено значение SwapFree")
		}
	})
}
