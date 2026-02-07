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

func (b Bytes) String() string {
	if b == 0 {
		return "0 B"
	}

	units := []string{"TB", "GB", "MB", "KB"}
	sizes := []Bytes{TB, GB, MB, KB}
	for i, size := range sizes {
		if b < size {
			continue
		}
		val := float64(b) / float64(size)
		if val >= 100 {
			return fmt.Sprintf("%.0f %s", val, units[i])
		}
		if val >= 10 {
			return fmt.Sprintf("%.1f %s", val, units[i])
		}
		return fmt.Sprintf("%.2f %s", val, units[i])
	}

	return fmt.Sprintf("%d B", b)
}
