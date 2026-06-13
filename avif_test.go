package avif

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"sync"
	"testing"
)

//go:embed testdata/test8.avif
var testAvif8 []byte

//go:embed testdata/test10.avif
var testAvif10 []byte

//go:embed testdata/test.avifs
var testAvifAnim []byte

func skipIfNoLibrary(tb testing.TB) {
	if err := Dynamic(); err != nil {
		fmt.Println(err)
		tb.Skip()
	}
}

func TestDecode(t *testing.T) {
	skipIfNoLibrary(t)

	img, _, err := decode(bytes.NewReader(testAvif8), false, false)
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img.Image[0], nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecode10(t *testing.T) {
	skipIfNoLibrary(t)

	img, _, err := decode(bytes.NewReader(testAvif10), false, false)
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img.Image[0], nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeAnim(t *testing.T) {
	skipIfNoLibrary(t)

	ret, _, err := decode(bytes.NewReader(testAvifAnim), false, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret.Image) != len(ret.Delay) {
		t.Errorf("got %d, want %d", len(ret.Delay), len(ret.Image))
	}

	if len(ret.Image) != 17 {
		t.Errorf("got %d, want %d", len(ret.Image), 17)
	}

	for _, img := range ret.Image {
		err = jpeg.Encode(io.Discard, img, nil)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestImageDecode(t *testing.T) {
	skipIfNoLibrary(t)

	img, _, err := image.Decode(bytes.NewReader(testAvif8))
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestImageDecodeAnim(t *testing.T) {
	skipIfNoLibrary(t)

	img, _, err := image.Decode(bytes.NewReader(testAvifAnim))
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeConfig(t *testing.T) {
	skipIfNoLibrary(t)

	_, cfg, err := decode(bytes.NewReader(testAvif8), true, false)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Width != 512 {
		t.Errorf("width: got %d, want %d", cfg.Width, 512)
	}

	if cfg.Height != 512 {
		t.Errorf("height: got %d, want %d", cfg.Height, 512)
	}
}

func TestEncode(t *testing.T) {
	skipIfNoLibrary(t)

	img, err := Decode(bytes.NewReader(testAvif8))
	if err != nil {
		t.Fatal(err)
	}

	w, err := writeCloser()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	err = encode(w, img, DefaultQuality, DefaultQuality, DefaultSpeed, image.YCbCrSubsampleRatio420)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodeSync(t *testing.T) {
	skipIfNoLibrary(t)

	wg := sync.WaitGroup{}
	ch := make(chan bool, 2)

	img, err := Decode(bytes.NewReader(testAvif8))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			ch <- true
			defer func() { <-ch; wg.Done() }()

			err = encode(io.Discard, img, DefaultQuality, DefaultQuality, DefaultSpeed, image.YCbCrSubsampleRatio420)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkDecode(b *testing.B) {
	skipIfNoLibrary(b)

	for b.Loop() {
		_, _, err := decode(bytes.NewReader(testAvif8), false, false)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDecodeConfig(b *testing.B) {
	skipIfNoLibrary(b)

	for b.Loop() {
		_, _, err := decode(bytes.NewReader(testAvif8), true, false)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	skipIfNoLibrary(b)

	img, err := Decode(bytes.NewReader(testAvif8))
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		err := encode(io.Discard, img, DefaultQuality, DefaultQuality, DefaultSpeed, image.YCbCrSubsampleRatio420)
		if err != nil {
			b.Error(err)
		}
	}
}

type discard struct{}

func (d discard) Close() error {
	return nil
}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

var discardCloser io.WriteCloser = discard{}

func writeCloser(s ...string) (io.WriteCloser, error) {
	if len(s) > 0 {
		f, err := os.Create(s[0])
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	return discardCloser, nil
}
