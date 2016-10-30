package main

import (
	"path/filepath"
	"testing"
)

// bunch of images for testing

func BenchmarkImgCompress(b *testing.B) {
	imgs, _ = filepath.Glob("original/windsurf/boards/*.jpg")
	for n := 0; n < b.N; n++ {
		ImgCompress(imgs, 50)
	}
}
