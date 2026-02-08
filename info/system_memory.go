package info

import (
	"fmt"
	"strings"
)

// SystemMemoryInfo содержит информацию
// о памяти системы, включая физическую память и своп.
type SystemMemoryInfo struct {
	TotalMemory     Bytes
	FreeMemory      Bytes
	AvailableMemory Bytes
	SwapTotal       Bytes
	SwapFree        Bytes
}

// String реализация fmt.Stringer
func (s SystemMemoryInfo) String() string {
	const str1 = "System Memory:"
	const str2 = "Total:"
	const str3 = "Used:"
	const str4 = "Available:"
	const str5 = "Swap Used:"
	// позиция, с которой начинается число
	const numStartPos = 11

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\n", str1)

	formatLine := func(label string, value string) {
		labelLen := len(label)
		spaces := numStartPos - labelLen
		if spaces < 0 {
			spaces = 0
		}

		idx := strings.Index(value, ".")
		if idx == 1 {
			value = " " + value
		}

		fmt.Fprintf(&sb, "%s%*s%s\n", label, spaces, "", value)
	}
	formatLine(str2, s.TotalMemory.String())
	formatLine(str3, s.UsedMemory())
	formatLine(str4, s.AvailableMemory.String())
	formatLine(str5, s.UsedSwap())
	return sb.String()
}

func (s SystemMemoryInfo) UsedMemory() string {
	used := s.TotalMemory - s.AvailableMemory
	percent := float64(used) / float64(s.TotalMemory) * 100
	return fmt.Sprintf("%s (%.1f%%)", used, percent)
}

func (s SystemMemoryInfo) UsedSwap() string {
	used := s.SwapTotal - s.SwapFree
	percent := float64(used) / float64(s.SwapTotal) * 100
	return fmt.Sprintf("%s (%.1f%%)", used, percent)
}
