package processing

import (
	"github.com/disintegration/imaging"
	"github.com/tlaceby/rawkit/core"
	"github.com/tlaceby/rawkit/processing/shaders"
)

// ApplyEdit applies core adjustments only. Simplest way to use rawkit.
//
// Example:
//
//	edit := processing.ImageEdit{
//		Exposure:   10,
//		Contrast:   5,
//		Saturation: -10,
//	}
//
//	result, err := processing.ApplyEdit(img, edit)
func ApplyEdit(img *core.ImageData, edit ImageEdit) (*core.ImageData, error) {
	coreUniforms := edit.GetUniforms()
	opengl, err := CreateOpenGLContext()
	if err != nil {
		return nil, err
	}

	defer opengl.Cleanup()
	newImg := CropImage(edit, img)
	if err = opengl.Upload(newImg.Data, img.Width, img.Height, img.Channels); err != nil {
		return nil, err
	}

	if err = opengl.RunShader(shaders.PrepareForDesktop(shaders.CORE_ADJUSTMENT_SHADER), coreUniforms); err != nil {
		return nil, err
	}

	data, err := opengl.Download()
	if err != nil {
		return nil, err
	}

	return &core.ImageData{
		Width:      img.Width,
		Height:     img.Height,
		Channels:   img.Channels,
		Colorspace: img.Colorspace,
		Data:       data,
	}, nil
}

// Shaders must be written WITHOUT #version or precision directives.
// These headers are automatically prepended to shader source
//
// Example usage:
//
//	brightnessShader := `
//		in vec2 uv;
//		out vec4 frag_color;
//		uniform sampler2D u_texture;
//		uniform float brightness;
//
//		void main() {
//			vec4 color = texture(u_texture, uv);
//			color.rgb += brightness;
//			frag_color = color;
//		}
//	`
//
//	passes := []processing.ShaderPass{
//		{
//			Name:     "brightness",
//			Shader:   brightnessShader,
//			Uniforms: processing.NewUniforms().AddFloat("brightness", 0.1),
//		},
//	}
//
//	result, err := processing.ApplyShaders(img, passes)
func ApplyShaders(img *core.ImageData, shaderPasses []ShaderPass) (*core.ImageData, error) {
	opengl, err := CreateOpenGLContext()
	if err != nil {
		return nil, err
	}

	defer opengl.Cleanup()
	if err = opengl.Upload(img.Data, img.Width, img.Height, img.Channels); err != nil {
		return nil, err
	}

	for _, pass := range shaderPasses {
		if err = opengl.RunShader(shaders.PrepareForDesktop(pass.Shader), pass.Uniforms); err != nil {
			return nil, err
		}
	}

	data, err := opengl.Download()
	if err != nil {
		return nil, err
	}

	return &core.ImageData{
		Width:      img.Width,
		Height:     img.Height,
		Channels:   img.Channels,
		Colorspace: img.Colorspace,
		Data:       data,
	}, nil
}

// applies core adjustments followed by custom shader passes. core uniforms (u_exposure, u_contrast, etc.) are injected into each shader pass.
//
// Shaders must be written WITHOUT #version or precision directives.
// These headers are automatically added for the target platform.
//
// Example usage:
//
//	edit := processing.ImageEdit{
//		Exposure: 15,
//		Contrast: 10,
//	}
//
//	denoiseShader := `
//		in vec2 uv;
//		out vec4 frag_color;
//		uniform sampler2D u_texture;
//		uniform float u_exposure;
//		uniform float strength;
//
//		void main() {
//			vec4 color = texture(u_texture, uv);
//			// denoise logic using strength and u_exposure
//			frag_color = color;
//		}
//	`
//
//	passes := []processing.ShaderPass{
//		{
//			Name:     "denoise",
//			Shader:   denoiseShader,
//			Uniforms: processing.NewUniforms().AddFloat("strength", 50.0),
//		},
//	}
//
//	result, err := processing.Process(img, edit, passes)
func Process(img *core.ImageData, edit ImageEdit, shaderPasses []ShaderPass) (*core.ImageData, error) {
	coreUniforms := edit.GetUniforms()
	opengl, err := CreateOpenGLContext()
	if err != nil {
		return nil, err
	}

	defer opengl.Cleanup()
	newImg := CropImage(edit, img)

	if err = opengl.Upload(newImg.Data, img.Width, img.Height, img.Channels); err != nil {
		return nil, err
	}

	if err = opengl.RunShader(shaders.PrepareForDesktop(shaders.CORE_ADJUSTMENT_SHADER), coreUniforms); err != nil {
		return nil, err
	}

	for _, pass := range shaderPasses {
		uniforms := make(Uniforms)
		uniforms.AddUniforms(pass.Uniforms)
		uniforms.AddUniforms(coreUniforms)

		if err = opengl.RunShader(shaders.PrepareForDesktop(pass.Shader), uniforms); err != nil {
			return nil, err
		}
	}

	data, err := opengl.Download()
	if err != nil {
		return nil, err
	}

	return &core.ImageData{
		Width:      img.Width,
		Height:     img.Height,
		Channels:   img.Channels,
		Colorspace: img.Colorspace,
		Data:       data,
	}, nil
}

// extracts a rectangular region from the image and returns the original image if crop dimensions are invalid or match original size.
func CropImage(edit ImageEdit, img *core.ImageData) *core.ImageData {
	if edit.CropWidth <= 0 || edit.CropHeight <= 0 {
		return img
	}

	if edit.CropLeft < 0 || edit.CropTop < 0 {
		return img
	}

	if edit.CropLeft+edit.CropWidth > img.Width {
		return img
	}

	if edit.CropTop+edit.CropHeight > img.Height {
		return img
	}

	if edit.CropWidth == img.Width && edit.CropHeight == img.Height && edit.CropLeft == 0 && edit.CropTop == 0 {
		return img
	}

	newData := make([]uint16, edit.CropWidth*edit.CropHeight*int(img.Channels))
	for y := 0; y < edit.CropHeight; y++ {
		srcY := edit.CropTop + y
		srcOffset := (srcY*img.Width + edit.CropLeft) * int(img.Channels)
		dstOffset := y * edit.CropWidth * int(img.Channels)

		// copy row
		copy(newData[dstOffset:dstOffset+edit.CropWidth*int(img.Channels)], img.Data[srcOffset:srcOffset+edit.CropWidth*int(img.Channels)])
	}

	return &core.ImageData{
		Width:      edit.CropWidth,
		Height:     edit.CropHeight,
		Channels:   img.Channels,
		Colorspace: img.Colorspace,
		Data:       newData,
	}
}

// specifies the resampling filter to use during image resizing.
type ResizeFilter int

const (
	// FilterNearest uses nearest-neighbor interpolation (fastest, lowest quality)
	FilterNearest ResizeFilter = iota
	// FilterBox uses box (average) sampling (fast, simple antialiasing)
	FilterBox
	// FilterLinear uses bilinear interpolation (balanced speed/quality)
	FilterLinear
	// FilterCubic uses bicubic interpolation (higher quality, slower)
	FilterCubic
	// FilterLanczos uses Lanczos resampling (best quality, slowest)
	FilterLanczos
)

// toImagingFilter converts ResizeFilter to imaging.ResampleFilter.
func (f ResizeFilter) toImagingFilter() imaging.ResampleFilter {
	switch f {
	case FilterNearest:
		return imaging.NearestNeighbor
	case FilterBox:
		return imaging.Box
	case FilterLinear:
		return imaging.Linear
	case FilterCubic:
		return imaging.CatmullRom
	case FilterLanczos:
		return imaging.Lanczos
	default:
		return imaging.Lanczos // default to best quality
	}
}

// Resize an image
func Resize(img *core.ImageData, width, height int, filter ResizeFilter) *core.ImageData {
	resized := imaging.Resize(img, width, height, filter.toImagingFilter())
	bounds := resized.Bounds()
	newWidth := bounds.Dx()
	newHeight := bounds.Dy()

	newData := make([]uint16, newWidth*newHeight*int(img.Channels))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			r, g, b, a := resized.At(x, y).RGBA()
			idx := (y*newWidth + x) * int(img.Channels)
			newData[idx] = uint16(r)
			newData[idx+1] = uint16(g)
			newData[idx+2] = uint16(b)

			if img.Channels == core.LIBRAW_CHANNELS_RGBA {
				newData[idx+3] = uint16(a)
			}
		}
	}

	return &core.ImageData{
		Width:      newWidth,
		Height:     newHeight,
		Channels:   img.Channels,
		Colorspace: img.Colorspace,
		Data:       newData,
	}
}
