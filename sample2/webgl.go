package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/prozacchiwawa/golang-webgl-examples/glUtils"
	"golang.org/x/image/math/f32"
)

type Scene struct {
	start float64
	gl    *webgl.Context

	squareVerticesBuffer    *js.Object
	mvMatrix                f32.Mat4
	shaderProgram           *js.Object
	vertexPositionAttribute int
	perspectiveMatrix       f32.Mat4
}

//
// start
//
// Called when the canvas is created to get the ball rolling.
// Figuratively, that is. There's nothing moving in this demo.
//
func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		panic("Error: " + err.Error())
	}

	scene := Scene{gl: gl}

	gl.ClearColor(0, 0, 0, 1)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// Initialize the shaders; this is where all the lighting for the
	// vertices and so forth is established.

	scene.initShaders()

	// Here's where we call the routine that builds all the objects
	// we'll be drawing.

	scene.initBuffers()

	// Set up to draw the scene periodically.

	scene.animate()
}

//
// initBuffers
//
// Initialize the buffers we'll need. For this demo, we just have
// one object -- a simple two-dimensional square.
//
func (s *Scene) initBuffers() {
	gl := s.gl

	// Create a buffer for the square's vertices.

	s.squareVerticesBuffer = gl.CreateBuffer()

	// Select the squareVerticesBuffer as the one to apply vertex
	// operations to from here out.

	gl.BindBuffer(gl.ARRAY_BUFFER, s.squareVerticesBuffer)

	// Now create an array of vertices for the square. Note that the Z
	// coordinate is always 0 here.

	vertices := []float32{
		1.0, 1.0, 0.0,
		-1.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, -1.0, 0.0,
	}

	// Now pass the list of vertices into WebGL to build the shape. We
	// do this by creating a Float32Array from the JavaScript array,
	// then use it to fill the current vertex buffer.

	gl.BufferData(gl.ARRAY_BUFFER, js.Global.Get("Float32Array").New(vertices), gl.STATIC_DRAW)
}

func (s *Scene) animate() {
	js.Global.Get("window").Call("requestAnimationFrame", func() { s.animate() })
	s.drawScene()
}

func (s *Scene) drawScene() {
	gl := s.gl

	// Clear the canvas before we start drawing on it.

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Establish the perspective with which we want to view the
	// scene. Our field of view is 45 degrees, with a width/height
	// ratio of 640:480, and we only want to see objects between 0.1 units
	// and 100 units away from the camera.

	s.perspectiveMatrix = glUtils.MakePerspective(45, 640.0/480.0, 0.1, 100.0)

	// Set the drawing position to the "identity" point, which is
	// the center of the scene.

	s.loadIdentity()

	// Now move the drawing position a bit to where we want to start
	// drawing the square.

	s.mvTranslate(f32.Vec3{0, 0, -6})

	// Draw the square by binding the array buffer to the square's vertices
	// array, setting attributes, and pushing it to GL.

	gl.BindBuffer(gl.ARRAY_BUFFER, s.squareVerticesBuffer)
	gl.VertexAttribPointer(s.vertexPositionAttribute, 3, gl.FLOAT, false, 0, 0)
	s.setMatrixUnforms()
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

//
// initShaders
//
// Initialize the shaders, so WebGL knows how to light our scene.
//
func (s *Scene) initShaders() {
	gl := s.gl
	fragmentShader := s.getShader("shader-fs")
	vertexShader := s.getShader("shader-vs")

	// Create the shader program

	s.shaderProgram = gl.CreateProgram()
	gl.AttachShader(s.shaderProgram, vertexShader)
	gl.AttachShader(s.shaderProgram, fragmentShader)
	gl.LinkProgram(s.shaderProgram)

	// If creating the shader program failed, alert

	status := gl.GetProgramParameterb(s.shaderProgram, gl.LINK_STATUS)
	if !status {
		panic("Could not initialize the shader program")
	}

	gl.UseProgram(s.shaderProgram)
	s.vertexPositionAttribute = gl.GetAttribLocation(s.shaderProgram, "aVertexPosition")
	gl.EnableVertexAttribArray(s.vertexPositionAttribute)
}

//
// getShader
//
// Loads a shader program by scouring the current document,
// looking for a script with the specified ID.
//
func (s *Scene) getShader(id string) *js.Object {
	gl := s.gl
	shaderScript := js.Global.Get("document").Call("getElementById", id)

	// Didn't find an element with the specified ID; abort.

	if shaderScript == nil {
		panic("Could not find element id " + id)
	}

	// Walk through the source element's children, building the
	// shader source string.

	theSource := shaderScript.Get("innerText")
	shaderType := shaderScript.Call("getAttribute", "type").String()

	// Now figure out what type of shader script we have,
	// based on its MIME type.

	var shader *js.Object
	if shaderType == "x-shader/x-fragment" {
		shader = gl.CreateShader(gl.FRAGMENT_SHADER)
	} else if shaderType == "x-shader/x-vertex" {
		shader = gl.CreateShader(gl.VERTEX_SHADER)
	} else {
		panic(fmt.Sprintf("Unknown shader type %v", shaderType)) // Unknown shader type
	}

	// Send the source to the shader object

	gl.ShaderSource(shader, theSource.String())

	// Compile the shader program

	gl.CompileShader(shader)

	status := gl.GetShaderParameter(shader, gl.COMPILE_STATUS)

	// See if it compiled successfully

	if !status.Bool() {
		panic(fmt.Sprintf("error compiling shader: %v", gl.GetShaderInfoLog(shader)))
	}

	return shader
}

func (self *Scene) loadIdentity() {
	self.mvMatrix = glUtils.Identity()
}

func (self *Scene) multMatrix(m f32.Mat4) {
	self.mvMatrix = glUtils.X4(self.mvMatrix, m)
}

func (self *Scene) mvTranslate(v f32.Vec3) {
	self.multMatrix(glUtils.TranslateMatrix(glUtils.Identity(), v))
}

func (self *Scene) setMatrixUnforms() {
	gl := self.gl
	pUniform := gl.GetUniformLocation(self.shaderProgram, "uPMatrix")
	gl.UniformMatrix4fv(pUniform, false, glUtils.Flatten(&self.perspectiveMatrix))

	mvUniform := gl.GetUniformLocation(self.shaderProgram, "uMVMatrix")
	gl.UniformMatrix4fv(mvUniform, false, glUtils.Flatten(&self.mvMatrix))
}
