package imgresizer

import (
	"path/filepath"
	"testing"
)

// bunch of images for testing
var imgs, _ = filepath.Glob("original/windsurf/boards/*.jpg")

func BenchmarkImgCompress(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ImgCompress(imgs, 50)
	}
}
