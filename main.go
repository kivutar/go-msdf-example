package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.GetError() // ignore INVALID_ENUM that GLEW raises when using OpenGL 3.2+

	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("Texture\x00")), 0)

	// Configure the vertex data
	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	image, err := NewImage("foo.png")
	if err != nil {
		panic(err)
	}

	fmt.Println(image)

	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Do OpenGL stuff.
		va := vertexArray(window, 100, 100, 512, 512, 1)

		gl.UseProgram(program)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.BindVertexArray(vao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, image)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(va)*4, gl.Ptr(va), gl.STATIC_DRAW)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		gl.BindVertexArray(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.UseProgram(0)
		gl.Disable(gl.BLEND)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
