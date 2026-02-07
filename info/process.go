package info

import "errors"

// PID идентификатор процесса в операционной системе.
type PID int

// ProcessName содержит имя исполняемого файла или отображаемое имя процесса.
type ProcessName string

// ProcessInfo представляет информацию о запущенном процессе.
type ProcessInfo struct {
	// PID уникальный идентификатор процесса в системе.
	PID PID

	// Name отображаемое имя процесса.
	Name ProcessName

	// MemoryUsage объём памяти, используемой процессом, в байтах.
	MemoryUsage Bytes
}

// NewProcessInfo создаёт ProcessInfo.
// pid — ID процесса, name — имя процесса,
// memory — потребление памяти в байтах.
// Возвращает ошибку при невалидном pid или пустом имени.
func NewProcessInfo(pid int, name string, memory uint64) (ProcessInfo, error) {
	vPID, err := newPID(pid)
	if err != nil {
		return ProcessInfo{}, err
	}

	vName, err := newProcessName(name)
	if err != nil {
		return ProcessInfo{}, err
	}

	return ProcessInfo{
		PID:         vPID,
		Name:        vName,
		MemoryUsage: Bytes(memory),
	}, nil
}

func newPID(pid int) (PID, error) {
	if pid < 0 {
		return 0,
			errors.New("Идентификатор процесса не может быть меньше 0")
	}
	return PID(pid), nil
}

func newProcessName(name string) (ProcessName, error) {
	if len(name) == 0 {
		return ProcessName("?"),
			errors.New("Имя процесса не может быть пустым")
	}
	return ProcessName(name), nil
}
