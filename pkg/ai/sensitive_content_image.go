package ai

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

const maxSensitiveContentImageEdge = 2048

var ErrInvalidSensitiveContentImage = errors.New("image must be JPEG or PNG with dimensions between 1 and 2048 pixels")

func decisionImagePart(data []byte) (decisionInputPart, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "jpeg" && format != "png") ||
		min(config.Width, config.Height) < 1 || max(config.Width, config.Height) > maxSensitiveContentImageEdge {
		return decisionInputPart{}, ErrInvalidSensitiveContentImage
	}
	return decisionInputPart{
		Type:     "input_image",
		ImageURL: "data:image/" + format + ";base64," + base64.StdEncoding.EncodeToString(data),
	}, nil
}
