package isoduration

import (
	"testing"
)

func BenchmarkParseDuration(b *testing.B) {
	str := "P30Y11.9M29.5D29.1WT10H30.5M5S"
	for i := 0; i < b.N; i++ {
		v, err := ParseDuration(str)
		if err != nil {
			panic(err)
		}
		v.ToTimeDuration()
	}
}
