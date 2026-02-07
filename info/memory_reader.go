package info

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// MemoryReader методы для работы с процессами
type MemoryReader interface {
	// ReadSystemMemory получить информацию о системной памяти
	ReadSystemMemory() (SystemMemoryInfo, error)

	// GetProcessList получить PIDы всех активных процессов
	GetProcessList() ([]PID, error)

	// ReadProcessMemory по PID возвращает
	// количество используемой резидентной памяти
	ReadProcessMemory(pid PID) (Bytes, error)
}

// DarwinMemoryReader для работы с процессами в Darwin (macOS)
type DarwinMemoryReader struct{}

func (d *DarwinMemoryReader) ReadSystemMemory() (SystemMemoryInfo, error) {
	var info SystemMemoryInfo

	// 1. Получаем общий объем памяти: sysctl -n hw.memsize
	totalMem, err := totalMemory()
	if err != nil {
		return info, fmt.Errorf("не получилось получить общий объем памяти: %w", err)
	}
	info.TotalMemory = Bytes(totalMem)

	// 2. Получаем статистику vm_stat
	vmStat, err := vmStat()
	if err != nil {
		return info, fmt.Errorf("не получилось выполнить vm_stat: %w", err)
	}
	pageSize := Bytes(4096)
	info.FreeMemory = pageSize * Bytes(vmStat["free"])
	info.AvailableMemory = pageSize * Bytes(vmStat["free"]+vmStat["inactive"])

	// 3. Получаем информацию о swap: sysctl -n vm.swapusage
	swapTotal, swapFree, err := swapUsage()
	if err != nil {
		return info, fmt.Errorf("не получилось получить swap usage: %w", err)
	}
	info.SwapTotal = Bytes(swapTotal)
	info.SwapFree = Bytes(swapFree)

	return info, nil
}

// swapUsage парсит вывод sysctl vm.swapusage
func swapUsage() (total, free uint64, err error) {
	cmd := exec.Command("sysctl", "-n", "vm.swapusage")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	line := strings.TrimSpace(string(output))

	// Парсим total
	totalRe := regexp.MustCompile(`total\s*=\s*([0-9.]+)([KMGT]?)`)
	totalMatch := totalRe.FindStringSubmatch(line)
	if len(totalMatch) > 0 {
		total = parseSize(totalMatch[1], totalMatch[2])
	}

	// Парсим free
	freeRe := regexp.MustCompile(`free\s*=\s*([0-9.]+)([KMGT]?)`)
	freeMatch := freeRe.FindStringSubmatch(line)
	if len(freeMatch) > 0 {
		free = parseSize(freeMatch[1], freeMatch[2])
	}

	return total, free, nil
}

// parseSize конвертирует строку с суффиксом (K, M, G, T) в байты
func parseSize(valueStr, unit string) uint64 {
	value, _ := strconv.ParseFloat(valueStr, 64)

	multiplier := 1.0
	switch unit {
	case "K":
		multiplier = 1024
	case "M":
		multiplier = 1024 * 1024
	case "G":
		multiplier = 1024 * 1024 * 1024
	case "T":
		multiplier = 1024 * 1024 * 1024 * 1024
	}

	return uint64(value * multiplier)
}

// totalMemory выполняет sysctl и парсит uint64 значение
func totalMemory() (uint64, error) {
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	value := strings.TrimSpace(string(output))
	return strconv.ParseUint(value, 10, 64)
}

// vmStat парсит вывод vm_stat
func vmStat() (map[string]uint64, error) {
	cmd := exec.Command("vm_stat")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make(map[string]uint64)
	lines := strings.Split(string(output), "\n")

	// Регулярное выражение для парсинга: "Pages free: 12345."
	re := regexp.MustCompile(`Pages\s+(\w+):\s+(\d+)\.`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := matches[1]
			value, err := strconv.ParseUint(matches[2], 10, 64)
			if err == nil {
				result[key] = value
			}
		}
	}

	return result, nil
}

func (d *DarwinMemoryReader) GetProcessList() ([]PID, error) {
	//  выполняем команду ps -e -o pid=
	cmd := exec.Command("ps", "-e", "-o", "pid=")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Ошибка запуска команды ps: %w", err)
	}

	// разбираем вывод команды
	outTxt := string(output)
	reader := strings.NewReader(outTxt)
	scanner := bufio.NewScanner(reader)
	pids := []PID{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil,
				fmt.Errorf("не удалось определить номер PID %q: %w", line, err)
		}
		pids = append(pids, PID(num))
	}
	// проверяем ошибки сканера
	if err := scanner.Err(); err != nil {
		return nil,
			fmt.Errorf("ошибка чтения данных из команды ps: %w", err)
	}

	return pids, nil
}

func (d *DarwinMemoryReader) ReadProcessMemory(pid PID) (Bytes, error) {
	numPID := strconv.Itoa(int(pid))
	// выполняем команду
	cmd := exec.Command("ps", "-p", numPID, "-o", "rss=")
	output, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return 0,
				fmt.Errorf("процесс %s не найден или доступ запрещен: %w", numPID, err)
		}
		return 0,
			fmt.Errorf("не удалось выполнить команду ps: %w", err)
	}

	// читаем вывод
	str := strings.TrimSpace(string(output))
	if len(str) == 0 {
		return 0,
			fmt.Errorf("команда ps вернула пустое значение %q: %w", str, err)
	}
	kbs, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0,
			fmt.Errorf("не удалось распарсить значение RSS %q: %w", str, err)
	}

	// конвертируем, т.к. ps отдает результат в килобайтах
	bytes := Bytes(kbs * 1024)

	return bytes, nil
}

// LinuxMemoryReader для работы с процессами в Linux
type LinuxMemoryReader struct{}

func (l *LinuxMemoryReader) ReadSystemMemory() (SystemMemoryInfo, error) {
	var info SystemMemoryInfo

	// Открываем файл /proc/meminfo
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return info,
			fmt.Errorf("не удалось открыть /proc/meminfo: %w", err)
	}
	defer file.Close()

	// Создаем аккумулятор и сканер для чтения файла построчно
	acc := map[string]uint64{}
	scanner := bufio.NewScanner(file)

	// Читаем файл построчно и парсим значения
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		//убираем двоеточие в конце ключа
		key := strings.TrimSuffix(fields[0], ":")
		// парсим значение
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		acc[key] = val
	}

	// Проверяем ошибки сканера
	if err := scanner.Err(); err != nil {
		return info,
			fmt.Errorf("не удалось прочитать /proc/meminfo: %w", err)
	}

	// Проверяем наличие обязательных полей
	fields := []string{"MemTotal", "MemFree", "MemAvailable", "SwapTotal", "SwapFree"}

	for _, field := range fields {
		if _, ok := acc[field]; !ok {
			return info,
				fmt.Errorf("не найдено требуемое поле %s в /proc/meminfo", field)
		}
	}

	// Заполняем структуру, конвертируя килобайты в байты
	info.TotalMemory = Bytes(acc["MemTotal"] * 1024)
	info.FreeMemory = Bytes(acc["MemFree"] * 1024)
	info.AvailableMemory = Bytes(acc["MemAvailable"] * 1024)
	info.SwapTotal = Bytes(acc["SwapTotal"] * 1024)
	info.SwapFree = Bytes(acc["SwapFree"] * 1024)

	return info, nil
}

func (l *LinuxMemoryReader) GetProcessList() ([]PID, error) {
	items, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("не удалость прочитать /proc: %w", err)
	}

	pids := []PID{}
	for _, item := range items {
		// пропускаем файлы
		if !item.IsDir() {
			continue
		}
		name := item.Name()
		num, err := strconv.Atoi(name)
		if err != nil {
			// пропускаем не номерные директории
			continue
		}
		pids = append(pids, PID(num))
	}

	return pids, nil
}

func (l *LinuxMemoryReader) ReadProcessMemory(pid PID) (Bytes, error) {
	numPID := strconv.Itoa(int(pid))
	// читаем файл с информацией о процессе
	file := fmt.Sprintf("/proc/%s/status", numPID)
	data, err := os.ReadFile(file)
	if err != nil {
		return 0,
			fmt.Errorf("не удалось прочитать %s: %w", file, err)
	}

	// ищем строку "VmRSS:" в выводе
	fields := []string{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "VmRSS:") {
			continue
		}
		fields = strings.Fields(line)
		break
	}
	if len(fields) != 3 {
		return 0,
			errors.New("Не найден или не определен формат VmRSS у процесса")
	}

	// получаем байты
	kbs, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось спарсить значение VmRSS: %w", err)
	}

	//
	return Bytes(kbs * 1024), nil
}
