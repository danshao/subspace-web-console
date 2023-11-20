package helpers

import (
	"fmt"
	"strings"
	"time"
)

const (
	BYTE     = 1.0             // 8 bits
	KIBIBYTE = 1024 * BYTE     // 2^10 bytes
	MEBIBYTE = 1024 * KIBIBYTE // 2^20 bytes
	GIBIBYTE = 1024 * MEBIBYTE // 2^30 bytes
	TEBIBYTE = 1024 * GIBIBYTE // 2^40 bytes
)

func ByteSizeFmt(bytes uint64) string {
	// (float64) 2^40 =  1.099511627776e+12
	// (uint64)  2^40 =  1099511627776
	// (float32) 2^40 =  1.0995116e+12
	var (
		unit  = ""
		value = float32(bytes)
	)

	switch {
	case bytes >= TEBIBYTE:
		unit = "TiB"
		value = value / TEBIBYTE
	case bytes >= GIBIBYTE:
		unit = "GiB"
		value = value / GIBIBYTE
	case bytes >= MEBIBYTE:
		unit = "MiB"
		value = value / MEBIBYTE
	case bytes >= KIBIBYTE:
		unit = "KiB"
		value = value / KIBIBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}
	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s %s", stringValue, unit)
}

func LocalTimeFmt(t interface{}) string {

	switch t.(type) {
	case string:
		localTime, _ := time.Parse("2006-01-02 (Mon) 15:04:05", t.(string))
		if localTime.IsZero() {
			return "Never"
		}
		return localTime.In(time.UTC).Format("2006-01-02 (Mon) 15:04:05 (MST -0700)")
	case time.Time:
		if t.(time.Time).IsZero() {
			return "Never"
		}
		return t.(time.Time).In(time.UTC).Format("2006-01-02 (Mon) 15:04:05 (MST -0700)")
	default:
		// _ = t
		return ""
	}
}
