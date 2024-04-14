package goid

import "testing"

func BenchmarkID3_Generate(b *testing.B) {
	id := NewID3()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id.Generate()
	}
}
