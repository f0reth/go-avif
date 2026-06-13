//go:build !unix && !darwin && !windows

package avif

import (
	"fmt"
	"image"
	"io"
)

var (
	dynamic    = false
	dynamicErr = fmt.Errorf("avif: platform not supported")
)

func decode(r io.Reader, configOnly, decodeAll bool) (*AVIF, image.Config, error) {
	return nil, image.Config{}, dynamicErr
}

func encode(w io.Writer, m image.Image, quality, qualityAlpha, speed int, subsampleRatio image.YCbCrSubsampleRatio) error {
	return dynamicErr
}

func loadLibrary() (uintptr, error) {
	return 0, dynamicErr
}
