package shaders

import _ "embed"

//go:embed core_adjustments.glsl
var CORE_ADJUSTMENT_SHADER string

//go:embed core_vertex_shader.glsl
var CORE_VERTEX_SHADER string
