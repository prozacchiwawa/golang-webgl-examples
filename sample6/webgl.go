package main

import (
  "math"
  "github.com/gopherjs/gopherjs/js"
  "github.com/gopherjs/webgl"
  "golang.org/x/image/math/f32"
  "github.com/prozacchiwawa/golang-webgl-examples/glUtils"
)

type Scene struct {
    start float64;
    gl *webgl.Context;

    cubeVerticesBuffer *js.Object;
    cubeVerticesTextureCoordBuffer *js.Object;
    cubeVerticesIndexBuffer *js.Object;
    cubeRotation float64;
    lastCubeUpdateTime *float64;

    cubeImage *js.Object;
    cubeTexture *js.Object;

    mvMatrix f32.Mat4;
    shaderProgram *js.Object;
    vertexPositionAttribute int;
    textureCoordAttribute int;
    perspectiveMatrix f32.Mat4;

    mvMatrixStack []f32.Mat4;
}

//
// start
//
// Called when the canvas is created to get the ball rolling.
//
func main() {
    document := js.Global.Get("document")
    canvas := document.Call("createElement", "canvas")
    document.Get("body").Call("appendChild", canvas)

    attrs := webgl.DefaultAttributes()
    attrs.Alpha = false

    gl, err := webgl.NewContext(canvas, attrs)
    if err != nil {
        panic("Error: "+err.Error());
    }

    s := Scene { gl: gl, mvMatrixStack: make([]f32.Mat4, 0, 0) };

    gl.ClearColor(0.0, 0.0, 0.0, 1.0);  // Clear to black, fully opaque
    gl.ClearDepth(1.0);                 // Clear everything
    gl.Enable(gl.DEPTH_TEST);           // Enable depth testing
    gl.DepthFunc(gl.LEQUAL);            // Near things obscure far things

    // Initialize the shaders; this is where all the lighting for the
    // vertices and so forth is established.

    s.initShaders();

    // Here's where we call the routine that builds all the objects
    // we'll be drawing.

    s.initBuffers();

    // Next, load and set up the textures we'll be using.

    s.initTextures();

    // Set up to draw the scene periodically.

    s.animate();
}

func (s *Scene) animate() {
    js.Global.Get("window").Call("requestAnimationFrame", func () { s.animate() });
    s.drawScene();
}

//
// initBuffers
//
// Initialize the buffers we'll need. For this demo, we just have
// one object -- a simple two-dimensional cube.
//
func (s *Scene) initBuffers() {
    gl := s.gl

    // Create a buffer for the cube's vertices.

    s.cubeVerticesBuffer = gl.CreateBuffer();

    // Select the cubeVerticesBuffer as the one to apply vertex
    // operations to from here out.

    gl.BindBuffer(gl.ARRAY_BUFFER, s.cubeVerticesBuffer);

    // Now create an array of vertices for the cube.

    vertices := []float32 {
        // Front face
        -1.0, -1.0,  1.0,
         1.0, -1.0,  1.0,
         1.0,  1.0,  1.0,
        -1.0,  1.0,  1.0,

        // Back face
        -1.0, -1.0, -1.0,
        -1.0,  1.0, -1.0,
         1.0,  1.0, -1.0,
         1.0, -1.0, -1.0,

        // Top face
        -1.0,  1.0, -1.0,
        -1.0,  1.0,  1.0,
         1.0,  1.0,  1.0,
         1.0,  1.0, -1.0,

        // Bottom face
        -1.0, -1.0, -1.0,
         1.0, -1.0, -1.0,
         1.0, -1.0,  1.0,
        -1.0, -1.0,  1.0,

        // Right face
         1.0, -1.0, -1.0,
         1.0,  1.0, -1.0,
         1.0,  1.0,  1.0,
         1.0, -1.0,  1.0,

        // Left face
        -1.0, -1.0, -1.0,
        -1.0, -1.0,  1.0,
        -1.0,  1.0,  1.0,
        -1.0,  1.0, -1.0,
    };

    // Now pass the list of vertices into WebGL to build the shape. We
    // do this by creating a Float32Array from the JavaScript array,
    // then use it to fill the current vertex buffer.

    gl.BufferData(gl.ARRAY_BUFFER, js.Global.Get("Float32Array").New(vertices), gl.STATIC_DRAW);

    // Map the texture onto the cube's faces.

    s.cubeVerticesTextureCoordBuffer = gl.CreateBuffer();
    gl.BindBuffer(gl.ARRAY_BUFFER, s.cubeVerticesTextureCoordBuffer);

    textureCoordinates := [] float32 {
        // Front
        0.0,  0.0,
        1.0,  0.0,
        1.0,  1.0,
        0.0,  1.0,
        // Back
        0.0,  0.0,
        1.0,  0.0,
        1.0,  1.0,
        0.0,  1.0,
        // Top
        0.0,  0.0,
        1.0,  0.0,
        1.0,  1.0,
        0.0,  1.0,
        // Bottom
        0.0,  0.0,
        1.0,  0.0,
        1.0,  1.0,
        0.0,  1.0,
        // Right
        0.0,  0.0,
        1.0,  0.0,
        1.0,  1.0,
        0.0,  1.0,
        // Left
        0.0,  0.0,
        1.0,  0.0,
        1.0,  1.0,
        0.0,  1.0,
    };

    gl.BufferData(gl.ARRAY_BUFFER, js.Global.Get("Float32Array").New(textureCoordinates),
                  gl.STATIC_DRAW);

    // Build the element array buffer; this specifies the indices
    // into the vertex array for each face's vertices.

    s.cubeVerticesIndexBuffer = gl.CreateBuffer();
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.cubeVerticesIndexBuffer);

    // This array defines each face as two triangles, using the
    // indices into the vertex array to specify each triangle's
    // position.

    cubeVertexIndices := [] float32 {
        0,  1,  2,      0,  2,  3,    // front
        4,  5,  6,      4,  6,  7,    // back
        8,  9,  10,     8,  10, 11,   // top
        12, 13, 14,     12, 14, 15,   // bottom
        16, 17, 18,     16, 18, 19,   // right
        20, 21, 22,     20, 22, 23,   // left
    };

    // Now send the element array to GL

    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,
        js.Global.Get("Uint16Array").New(cubeVertexIndices), gl.STATIC_DRAW);
}

//
// initTextures
//
// Initialize the textures we'll be using, then initiate a load of
// the texture images. The handleTextureLoaded() callback will finish
// the job; it gets called each time a texture finishes loading.
//
func (s *Scene) initTextures() {
    gl := s.gl;
    s.cubeTexture = gl.CreateTexture();
    s.cubeImage = js.Global.Get("Image").New();
    s.cubeImage.Set("onload", func() { s.handleTextureLoaded(s.cubeImage, s.cubeTexture); });
    s.cubeImage.Set("src", "cubetexture.png");
}

func (s *Scene) handleTextureLoaded(image *js.Object, texture *js.Object) {
    gl := s.gl;
    js.Global.Get("console").Call("log", "handleTextureLoaded, image = " + image.String());
    gl.BindTexture(gl.TEXTURE_2D, texture);
    gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA,
          gl.UNSIGNED_BYTE, image);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST);
    gl.GenerateMipmap(gl.TEXTURE_2D);
    gl.BindTexture(gl.TEXTURE_2D, nil);
}

//
// drawScene
//
// Draw the scene.
//
func (s *Scene) drawScene() {
    gl := s.gl;

    // Clear the canvas before we start drawing on it.

    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

    // Establish the perspective with which we want to view the
    // scene. Our field of view is 45 degrees, with a width/height
    // ratio of 640:480, and we only want to see objects between 0.1 units
    // and 100 units away from the camera.

    s.perspectiveMatrix = glUtils.MakePerspective(45, 640.0/480.0, 0.1, 100.0);

    // Set the drawing position to the "identity" point, which is
    // the center of the scene.

    s.loadIdentity();

    // Now move the drawing position a bit to where we want to start
    // drawing the cube.

    s.mvTranslate(f32.Vec3 {-0.0, 0.0, -6.0});

    // Save the current matrix, then rotate before we draw.

    s.mvPushMatrix(nil);
    s.mvRotate(s.cubeRotation, f32.Vec3 {1, 0, 1});

    // Draw the cube by binding the array buffer to the cube's vertices
    // array, setting attributes, and pushing it to GL.

    gl.BindBuffer(gl.ARRAY_BUFFER, s.cubeVerticesBuffer);
    gl.VertexAttribPointer(s.vertexPositionAttribute, 3, gl.FLOAT, false, 0, 0);

    // Set the texture coordinates attribute for the vertices.

    gl.BindBuffer(gl.ARRAY_BUFFER, s.cubeVerticesTextureCoordBuffer);
    gl.VertexAttribPointer(s.textureCoordAttribute, 2, gl.FLOAT, false, 0, 0);

    // Specify the texture to map onto the faces.

    gl.ActiveTexture(gl.TEXTURE0);
    gl.BindTexture(gl.TEXTURE_2D, s.cubeTexture);
    gl.Uniform1i(gl.GetUniformLocation(s.shaderProgram, "uSampler"), 0);

    // Draw the cube.

    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.cubeVerticesIndexBuffer);
    s.setMatrixUniforms();
    gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_SHORT, 0);

    // Restore the original matrix

    s.mvPopMatrix();

    // Update the rotation for the next draw, if it's time to do so.

    currentTime := js.Global.Get("Date").New().Call("getTime").Float();
    if s.lastCubeUpdateTime != nil {
        delta := currentTime - *s.lastCubeUpdateTime;

        s.cubeRotation += (30 * delta) / 1000.0;
    }

    s.lastCubeUpdateTime = &currentTime;
}

//
// initShaders
//
// Initialize the shaders, so WebGL knows how to light our scene.
//
func (s *Scene) initShaders() {
    gl := s.gl;
    fragmentShader := s.getShader("shader-fs");
    vertexShader := s.getShader("shader-vs");

    // Create the shader program

    s.shaderProgram = gl.CreateProgram();
    gl.AttachShader(s.shaderProgram, vertexShader);
    gl.AttachShader(s.shaderProgram, fragmentShader);
    gl.LinkProgram(s.shaderProgram);

    // If creating the shader program failed, alert

    if !gl.GetProgramParameterb(s.shaderProgram, gl.LINK_STATUS) {
        panic("Unable to initialize the shader program.");
    }

    gl.UseProgram(s.shaderProgram);

    s.vertexPositionAttribute = gl.GetAttribLocation(s.shaderProgram, "aVertexPosition");
    gl.EnableVertexAttribArray(s.vertexPositionAttribute);

    s.textureCoordAttribute = gl.GetAttribLocation(s.shaderProgram, "aTextureCoord");
    gl.EnableVertexAttribArray(s.textureCoordAttribute);
}

//
// getShader
//
// Loads a shader program by scouring the current document,
// looking for a script with the specified ID.
//
func (s *Scene) getShader(id string) *js.Object {
    gl := s.gl;
    shaderScript := js.Global.Get("document").Call("getElementById",id);

    // Didn't find an element with the specified ID; abort.

    if shaderScript == nil {
        panic("Couldn't get shader script " + id);
    }

    // Walk through the source element's children, building the
    // shader source string.

    theSource := shaderScript.Get("innerText").String();

    // Now figure out what type of shader script we have,
    // based on its MIME type.

    var shader *js.Object;
    shaderType := shaderScript.Get("type").String();

    if (shaderType == "x-shader/x-fragment") {
        shader = gl.CreateShader(gl.FRAGMENT_SHADER);
    } else if (shaderType == "x-shader/x-vertex") {
        shader = gl.CreateShader(gl.VERTEX_SHADER);
    } else {
        panic("Unknown shader type " + shaderType);  // Unknown shader type
    }

    // Send the source to the shader object

    gl.ShaderSource(shader, theSource);

    // Compile the shader program

    gl.CompileShader(shader);

    // See if it compiled successfully

    if !gl.GetShaderParameter(shader, gl.COMPILE_STATUS).Bool() {
        panic("An error occurred compiling the shaders: " + gl.GetShaderInfoLog(shader));
    }

    return shader;
}

//
// Matrix utility functions
//

func (s *Scene) loadIdentity() {
    s.mvMatrix = glUtils.Identity();
}

func (s *Scene) multMatrix(m *f32.Mat4) {
    s.mvMatrix = glUtils.X4(s.mvMatrix, *m);
}

func (s *Scene) mvTranslate(v f32.Vec3) {
    i := glUtils.Identity();
    tx := glUtils.TranslateMatrix(i, v);
    s.multMatrix(&tx);
}

func (s *Scene) setMatrixUniforms() {
    gl := s.gl;
    pUniform := gl.GetUniformLocation(s.shaderProgram, "uPMatrix");
    gl.UniformMatrix4fv(pUniform, false, glUtils.Flatten(&s.perspectiveMatrix));

    mvUniform := gl.GetUniformLocation(s.shaderProgram, "uMVMatrix");
    gl.UniformMatrix4fv(mvUniform, false, glUtils.Flatten(&s.mvMatrix));
}

func (s *Scene) mvPushMatrix(m *f32.Mat4) {
    if m != nil {
        s.mvMatrixStack = append(s.mvMatrixStack, *m);
        s.mvMatrix = *m;
    } else {
        s.mvMatrixStack = append(s.mvMatrixStack, s.mvMatrix);
    }
}

func (s *Scene) mvPopMatrix() f32.Mat4 {
    if len(s.mvMatrixStack) == 0 {
        panic("Can't pop from an empty matrix stack.");
    }

    s.mvMatrix = s.mvMatrixStack[len(s.mvMatrixStack)-1];
    s.mvMatrixStack = s.mvMatrixStack[0:len(s.mvMatrixStack)-1];
    return s.mvMatrix;
}

func (s *Scene) mvRotate(angle float64, v f32.Vec3) {
    inRadians := angle * math.Pi / 180.0;
    m := glUtils.RotateMatrix(inRadians, v);
    s.multMatrix(&m);
}
