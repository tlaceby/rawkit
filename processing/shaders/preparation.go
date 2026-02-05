package shaders

import "strings"

const (
	webGLHeader string = `#version 300 es
precision highp float;
`

	desktopHeader string = `#version 410 core
`
)

// converts a common shader to WebGL format
func PrepareForWebGL(shader string) string {
	return webGLHeader + removeDirectives(shader)
}

// converts a common shader to Desktop OpenGL format
func PrepareForDesktop(shader string) string {
	return desktopHeader + removeDirectives(shader)
}

func removeDirectives(shader string) string {
	lines := strings.Split(shader, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#version") {
			continue
		}
		if strings.HasPrefix(trimmed, "precision ") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
