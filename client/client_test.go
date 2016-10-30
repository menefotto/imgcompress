package main

import "testing"

func BenchmarkImgCompress(b *testing.B) {

	for n := 0; n < b.N; n++ {
		Send()
	}
}
