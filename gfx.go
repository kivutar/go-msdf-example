package main

import (
	"errors"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func textureLoad(rgba *image.RGBA) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

// NewImage opens an image file, upload it the the GPU and returns the texture id
func NewImage(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, errors.New("incorrect stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return textureLoad(rgba), nil
}

// XYWHTo4points converts coordinates from (x, y, width, height) to (x1, y1, x2, y2, x3, y3, x4, y4)
func XYWHTo4points(x, y, w, h, fbh float32) (x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	x1 = x
	x2 = x
	x3 = x + w
	x4 = x + w
	y1 = fbh - (y + h)
	y2 = fbh - y
	y3 = fbh - (y + h)
	y4 = fbh - y
	return
}

func vertexArray(window *glfw.Window, x, y, w, h, scale float32) []float32 {
	fbw, fbh := window.GetFramebufferSize()
	ffbw := float32(fbw)
	ffbh := float32(fbh)

	w *= scale
	h *= scale

	x1, y1, x2, y2, x3, y3, x4, y4 := XYWHTo4points(x, y, w, h, ffbh)

	return []float32{
		//  X, Y, U, V
		x1/ffbw*2 - 1, y1/ffbh*2 - 1, 0, 1, // left-bottom
		x2/ffbw*2 - 1, y2/ffbh*2 - 1, 0, 0, // left-top
		x3/ffbw*2 - 1, y3/ffbh*2 - 1, 1, 1, // right-bottom
		x4/ffbw*2 - 1, y4/ffbh*2 - 1, 1, 0, // right-top
	}
}

var vertices = []float32{
	//  X, Y, U, V
	-1.0, -1.0, 0.0, 1.0, // left-bottom
	-1.0, 1.0, 0.0, 0.0, // left-top
	1.0, -1.0, 1.0, 1.0, // right-bottom
	1.0, 1.0, 1.0, 0.0, // right-top
}
