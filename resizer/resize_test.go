package resizer

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/tokubai/kinu/engine"
)

func TestMain(m *testing.M) {
	engine.Initialize()
	code := m.Run()
	engine.Finalize()
	os.Exit(code)
}

func newTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode returned error: %v", err)
	}
	return buf.Bytes()
}

func TestResize(t *testing.T) {
	tests := []struct {
		name       string
		option     *ResizeOption
		wantWidth  int
		wantHeight int
		wantFormat string
	}{
		{"fixed width", &ResizeOption{Width: 100, Format: "png"}, 100, 50, "png"},
		{"auto crop", &ResizeOption{Width: 100, Height: 100, NeedsAutoCrop: true, Format: "png"}, 100, 100, "png"},
		{"to jpeg with alpha removed", &ResizeOption{Width: 200, Quality: 70, SourceContentType: "image/png", Format: "jpeg"}, 200, 100, "jpeg"},
	}

	src := newTestPNG(t, 400, 200)
	for _, tt := range tests {
		result := Resize(src, tt.option)
		if result.err != nil {
			t.Errorf("%s: Resize returned error: %v", tt.name, result.err)
			continue
		}

		config, format, err := image.DecodeConfig(bytes.NewReader(result.image))
		if err != nil {
			t.Errorf("%s: failed to decode result: %v", tt.name, err)
			continue
		}
		if config.Width != tt.wantWidth || config.Height != tt.wantHeight || format != tt.wantFormat {
			t.Errorf("%s: got %dx%d %s, want %dx%d %s", tt.name,
				config.Width, config.Height, format, tt.wantWidth, tt.wantHeight, tt.wantFormat)
		}
	}
}

func TestResizeInvalidImage(t *testing.T) {
	if result := Resize([]byte("not an image"), &ResizeOption{Width: 100}); result.err == nil {
		t.Error("Resize with invalid image should return error")
	}
}
