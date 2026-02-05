package processing

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/tlaceby/rawkit/core"
	"github.com/tlaceby/rawkit/processing/shaders"
)

type OpenGLContext struct {
	window *glfw.Window

	program  uint32 // id of the program
	width    int
	height   int
	channels int

	useTexture2ForOutput bool // whether to set the framebuffer (output) to texture 2

	texture1          uint32
	texture2          uint32
	frameBufferObject uint32

	vertexArrayBuffer uint32
}

// Based on what pass we are on, returns the correct (input, output) textures
func (ctx *OpenGLContext) getTextures() (uint32, uint32) {
	if ctx.useTexture2ForOutput {
		return ctx.texture1, ctx.texture2
	}

	return ctx.texture2, ctx.texture1
}

func CreateOpenGLContext() (*OpenGLContext, error) {
	// opengl MUST live on the main thread
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(1, 1, "", nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}

	// create and link program with vertex shader
	program := gl.CreateProgram()

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	preparedVertexShader := shaders.PrepareForDesktop(shaders.CORE_VERTEX_SHADER) // ADD THIS
	vPtr, vFree := gl.Strs(preparedVertexShader + "\x00")
	defer vFree()

	gl.ShaderSource(vertexShader, 1, vPtr, nil)
	gl.CompileShader(vertexShader)

	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(vertexShader, logLength, nil, &log[0])
		return nil, errors.New("vertex shader compilation failed: " + string(log))
	}

	gl.AttachShader(program, vertexShader)
	gl.LinkProgram(program)

	// check linking
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])
		return nil, errors.New("program linking failed: " + string(log))
	}

	gl.UseProgram(program)
	gl.DeleteShader(vertexShader)

	var vao uint32
	gl.GenVertexArrays(1, &vao) // fake vao

	ctx := &OpenGLContext{
		window:               window,
		program:              program,
		useTexture2ForOutput: true,
		vertexArrayBuffer:    vao,
	}

	return ctx, nil
}
func (ctx *OpenGLContext) Upload(data []uint16, width, height int, channels core.Channels) error {
	if len(data) == 0 {
		return errors.New("empty image data")
	}

	ctx.cleanupTextures()

	ctx.width = width
	ctx.height = height
	ctx.channels = int(channels)

	// texture formats
	internalFormat := int32(gl.RGB16)
	format := uint32(gl.RGB)

	if channels == core.LIBRAW_CHANNELS_RGBA {
		internalFormat = int32(gl.RGBA16)
		format = gl.RGBA
	}

	// input texture
	gl.GenTextures(1, &ctx.texture1)
	gl.BindTexture(gl.TEXTURE_2D, ctx.texture1)
	setTextureParams()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		internalFormat,
		int32(width),
		int32(height),
		0,
		format,
		gl.UNSIGNED_SHORT,
		unsafe.Pointer(&data[0]),
	)

	// output texture
	gl.GenTextures(1, &ctx.texture2)
	gl.BindTexture(gl.TEXTURE_2D, ctx.texture2)
	setTextureParams()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		internalFormat,
		int32(width),
		int32(height),
		0,
		format,
		gl.UNSIGNED_SHORT,
		nil, // allocate only
	)

	gl.BindTexture(gl.TEXTURE_2D, 0) // un-bind

	// frame buffer (attach the output texture as the location to render too)
	gl.GenFramebuffers(1, &ctx.frameBufferObject)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ctx.frameBufferObject)
	gl.FramebufferTexture2D(
		gl.FRAMEBUFFER,
		gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D,
		ctx.texture2,
		0,
	)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return errors.New("framebuffer incomplete")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return nil
}

func (ctx *OpenGLContext) RunShader(shader string, uniforms Uniforms) error {
	inputTexture, outputTexture := ctx.getTextures()

	gl.BindFramebuffer(gl.FRAMEBUFFER, ctx.frameBufferObject)
	gl.FramebufferTexture2D(
		gl.FRAMEBUFFER,
		gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D,
		outputTexture,
		0,
	)

	// specify drawbuffer: out vec4 fragColor
	drawBuffers := []uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &drawBuffers[0])
	gl.Viewport(0, 0, int32(ctx.width), int32(ctx.height))

	// input texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, inputTexture)

	// compile fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	shaderPtr, shaderFree := gl.Strs(shader + "\x00")
	defer shaderFree()

	gl.ShaderSource(fragmentShader, 1, shaderPtr, nil)
	gl.CompileShader(fragmentShader)

	// compilation status
	var status int32
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, &log[0])
		gl.DeleteShader(fragmentShader)
		return errors.New("shader compilation failed: " + string(log))
	}

	gl.AttachShader(ctx.program, fragmentShader)
	defer func() {
		gl.DetachShader(ctx.program, fragmentShader)
		gl.DeleteShader(fragmentShader)
	}()

	gl.LinkProgram(ctx.program)

	// enure link status GUD
	gl.GetProgramiv(ctx.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(ctx.program, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(ctx.program, logLength, nil, &log[0])
		return errors.New("program linking failed: " + string(log))
	}

	gl.UseProgram(ctx.program)

	// set texture sampler uniform
	textureLocation := gl.GetUniformLocation(ctx.program, gl.Str("u_texture\x00"))
	if textureLocation != -1 {
		gl.Uniform1i(textureLocation, 0)
	}

	// custom uniforms
	for uniformName, uniform := range uniforms {
		namePtr := gl.Str(uniformName + "\x00")
		uniformID := gl.GetUniformLocation(ctx.program, namePtr)

		if uniformID == -1 {
			continue // uniform not used in shader or optimized out
		}

		floatVec := uniform.Data

		switch uniform.Type {
		case UNIFORM_FLOAT:
			gl.Uniform1f(uniformID, floatVec[0])
		case UNIFORM_INT:
			gl.Uniform1i(uniformID, uniform.Int)
		case UNIFORM_VEC2:
			gl.Uniform2f(uniformID, floatVec[0], floatVec[1])
		case UNIFORM_VEC3:
			gl.Uniform3f(uniformID, floatVec[0], floatVec[1], floatVec[2])
		case UNIFORM_VEC4:
			gl.Uniform4f(uniformID, floatVec[0], floatVec[1], floatVec[2], floatVec[3])
		}
	}

	// Draw full-screen quad
	gl.BindVertexArray(ctx.vertexArrayBuffer)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, 0, 6) // 2 triangle vertex array
	gl.BindVertexArray(0)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// set for next pass
	ctx.useTexture2ForOutput = !ctx.useTexture2ForOutput
	return nil
}

func (ctx *OpenGLContext) Reset() {
	ctx.useTexture2ForOutput = true
}
func (ctx *OpenGLContext) Download() ([]uint16, error) {
	inputTexture, _ := ctx.getTextures()

	gl.BindFramebuffer(gl.FRAMEBUFFER, ctx.frameBufferObject)
	gl.FramebufferTexture2D(
		gl.FRAMEBUFFER,
		gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D,
		inputTexture,
		0,
	)

	pixelCount := ctx.width * ctx.height * ctx.channels
	result := make([]uint16, pixelCount)

	format := gl.RGB
	if ctx.channels == 4 {
		format = gl.RGBA
	}

	gl.ReadPixels(
		0,
		0,
		int32(ctx.width),
		int32(ctx.height),
		uint32(format),
		gl.UNSIGNED_SHORT,
		unsafe.Pointer(&result[0]),
	)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	ctx.useTexture2ForOutput = true

	return result, nil
}

func (ctx *OpenGLContext) Cleanup() {
	ctx.cleanupTextures()

	if ctx.frameBufferObject != 0 {
		gl.DeleteFramebuffers(1, &ctx.frameBufferObject)
	}

	if ctx.program != 0 {
		gl.DeleteProgram(ctx.program)
	}

	ctx.window.Destroy()
	glfw.Terminate()
}

func (ctx *OpenGLContext) cleanupTextures() {
	if ctx.texture1 != 0 {
		gl.DeleteTextures(1, &ctx.texture1)
		ctx.texture1 = 0
	}

	if ctx.texture2 != 0 {
		gl.DeleteTextures(1, &ctx.texture2)
		ctx.texture2 = 0
	}
}

func setTextureParams() {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}
