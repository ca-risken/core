package ai

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"
)

func encodeSensitiveContentTestImage(t *testing.T, format string, width, height int) []byte {
	t.Helper()

	var encoded bytes.Buffer
	frame := image.NewRGBA(image.Rect(0, 0, width, height))
	var err error
	switch format {
	case "jpeg":
		err = jpeg.Encode(&encoded, frame, nil)
	case "png":
		err = png.Encode(&encoded, frame)
	case "gif":
		err = gif.Encode(&encoded, frame, nil)
	default:
		t.Fatalf("unexpected image format %q", format)
	}
	if err != nil {
		t.Fatal(err)
	}
	return encoded.Bytes()
}

func TestDecisionImagePartValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format string
		width  int
		height int
		want   error
	}{
		{name: "PNG at width limit", format: "png", width: 2048, height: 1},
		{name: "JPEG at height limit", format: "jpeg", width: 1, height: 2048},
		{name: "over width limit", format: "jpeg", width: 2049, height: 1, want: ErrInvalidSensitiveContentImage},
		{name: "over height limit", format: "png", width: 1, height: 2049, want: ErrInvalidSensitiveContentImage},
		{name: "unsupported GIF", format: "gif", width: 1, height: 1, want: ErrInvalidSensitiveContentImage},
		{name: "invalid header", want: ErrInvalidSensitiveContentImage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := []byte("invalid image")
			if tt.format != "" {
				input = encodeSensitiveContentTestImage(t, tt.format, tt.width, tt.height)
			}
			_, got := decisionImagePart(input)
			if !errors.Is(got, tt.want) {
				t.Errorf("decisionImagePart() error = %v, want %v", got, tt.want)
			}
		})
	}
}
