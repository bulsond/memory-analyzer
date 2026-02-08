package info

import "fmt"

// Bytes тип хранения значения величины памяти
type Bytes uint64

const (
	KB Bytes = 1024
	MB       = 1024 * KB
	GB       = 1024 * MB
	TB       = 1024 * GB
)

// KBytes размер в килобайтах
func KBytes(n uint64) Bytes {
	return Bytes(n) * KB
}

// MBytes размер в мегабайтах
func MBytes(n uint64) Bytes {
	return Bytes(n) * MB
}

// GBytes размер в гигабайтах
func GBytes(n uint64) Bytes {
	return Bytes(n) * GB
}

// TBytes размер в терабайтах
func TBytes(n uint64) Bytes {
	return Bytes(n) * TB
}

// Uint64 обратное преобразование
func (b Bytes) Uint64() uint64 {
	return uint64(b)
}

// String отображение объема с указанием единиц размерности
func (b Bytes) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
