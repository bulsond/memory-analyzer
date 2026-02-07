package info

// SystemMemoryInfo содержит информацию
// о памяти системы, включая физическую память и своп.
type SystemMemoryInfo struct {
	TotalMemory     Bytes
	FreeMemory      Bytes
	AvailableMemory Bytes
	SwapTotal       Bytes
	SwapFree        Bytes
}
