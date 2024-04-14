package goid

import "testing"

func BenchmarkID_Generate(b *testing.B) {
	id := NewID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id.Generate()
	}
}
