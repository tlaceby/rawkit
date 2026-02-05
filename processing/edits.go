package processing

import (
	"reflect"
	"strings"

	"github.com/tlaceby/rawkit/core"
)

// ----------------------------------------------------------------------------
// Image Edit
// ----------------------------------------------------------------------------

// ImageEdit represents core adjustments to apply during processing
type ImageEdit struct {
	// Crop
	CropTop    int
	CropLeft   int
	CropWidth  int
	CropHeight int

	// White Balance
	Temperature float32 // -100 to 100 (cool to warm)
	Tint        float32 // -100 to 100 (green to magenta)

	// Tone
	Exposure   float32 // -100 to 100
	Contrast   float32 // -100 to 100
	Highlights float32 // -100 to 100
	Shadows    float32 // -100 to 100
	Whites     float32 // -100 to 100
	Blacks     float32 // -100 to 100

	// Presence
	Clarity    float32 // -100 to 100
	Vibrance   float32 // -100 to 100
	Saturation float32 // -100 to 100

	// Detail
	Sharpness float32 // 0 to 100
}

// ----------------------------------------------------------------------------
// Shader Uniforms
// ----------------------------------------------------------------------------

// UniformType represents the type of a shader uniform
type UniformType int

const (
	UNIFORM_FLOAT UniformType = iota
	UNIFORM_INT
	UNIFORM_VEC2 // float
	UNIFORM_VEC3 // float
	UNIFORM_VEC4 // float
)

// holds a typed uniform value
type Uniforms map[string]UniformValue
type UniformValue struct {
	Type UniformType
	Data [4]float32 // float, vec2, vec3, or vec4 uniforms
	Int  int32      // int uniforms
}

func IntUniform(v int32) UniformValue {
	return UniformValue{Type: UNIFORM_INT, Int: v}
}

func FloatUniform(v float32) UniformValue {
	return UniformValue{Type: UNIFORM_FLOAT, Data: [4]float32{v, 0, 0, 0}}
}

func Vec2Uniform(x, y float32) UniformValue {
	return UniformValue{Type: UNIFORM_VEC2, Data: [4]float32{x, y, 0, 0}}
}

func Vec3Uniform(x, y, z float32) UniformValue {
	return UniformValue{Type: UNIFORM_VEC3, Data: [4]float32{x, y, z, 0}}
}

func Vec4Uniform(x, y, z, w float32) UniformValue {
	return UniformValue{Type: UNIFORM_VEC4, Data: [4]float32{x, y, z, w}}
}

func NewUniforms() Uniforms {
	return make(map[string]UniformValue)
}

func (u Uniforms) AddFloat(name string, v float32) Uniforms {
	u[name] = FloatUniform(v)
	return u
}

func (u Uniforms) AddInt(name string, i int32) Uniforms {
	u[name] = IntUniform(i)
	return u
}

func (u Uniforms) AddVec2(name string, x, y float32) Uniforms {
	u[name] = Vec2Uniform(x, y)
	return u
}

func (u Uniforms) AddVec3(name string, x, y, z float32) Uniforms {
	u[name] = Vec3Uniform(x, y, z)
	return u
}

func (u Uniforms) AddVec4(name string, x, y, z, w float32) Uniforms {
	u[name] = Vec4Uniform(x, y, z, w)
	return u
}

func (u Uniforms) AddUniforms(newUniforms Uniforms) {
	for name, uniform := range newUniforms {
		switch uniform.Type {
		case UNIFORM_INT:
			u.AddInt(name, uniform.Int)
		case UNIFORM_FLOAT:
			u.AddFloat(name, uniform.Data[0])
		case UNIFORM_VEC2:
			u.AddVec2(name, uniform.Data[0], uniform.Data[1])
		case UNIFORM_VEC3:
			u.AddVec3(name, uniform.Data[0], uniform.Data[1], uniform.Data[2])
		case UNIFORM_VEC4:
			u.AddVec4(name, uniform.Data[0], uniform.Data[1], uniform.Data[2], uniform.Data[3])
		}
	}
}

// represents a single shader to run with its uniforms
type ShaderPass struct {
	Name     string   // For debugging/logging
	Shader   string   // GLSL fragment shader string
	Uniforms Uniforms // Custom uniforms
}

// interface for headless GPU processing
type IGPUContext interface {
	Upload(data []uint16, width, height int, channels core.Channels) error
	RunShader(shader string, uniforms Uniforms) error
	Download() ([]uint16, error)
	Cleanup()
	Reset()
}

func (edit ImageEdit) GetUniforms() Uniforms {
	uniforms := make(Uniforms)
	v := reflect.ValueOf(edit)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// ignore crop fields as no uniforms are ever created
		if strings.HasPrefix(field.Name, "Crop") {
			continue
		}

		name := "u_" + strings.ToLower(field.Name)

		switch value.Kind() {
		case reflect.Float32:
			uniforms.AddFloat(name, float32(value.Float()))
		case reflect.Int:
			uniforms.AddInt(name, int32(value.Int()))
		}
	}

	return uniforms
}
