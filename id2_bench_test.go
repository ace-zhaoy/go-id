package goid

import "testing"

func BenchmarkID2_Generate(b *testing.B) {
	id := NewID2()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id.Generate()
	}
}
