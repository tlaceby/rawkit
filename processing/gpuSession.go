package processing

import (
	"github.com/tlaceby/rawkit/core"
	"github.com/tlaceby/rawkit/processing/shaders"
)

// maintains GPU state for interactive editing with custom shader passes
type GPUSession struct {
	ctx           IGPUContext
	originalData  *core.ImageData
	croppedData   *core.ImageData
	lastCrop      ImageEdit
	currentEdit   ImageEdit
	currentPasses []ShaderPass
}

// creates a session for interactive editing with custom shaders
func NewGPUSession(img *core.ImageData, ctx IGPUContext, edit ImageEdit, passes []ShaderPass) (*GPUSession, error) {
	croppedData := CropImage(edit, img)
	if err := ctx.Upload(croppedData.Data, croppedData.Width, croppedData.Height, croppedData.Channels); err != nil {
		ctx.Cleanup()
		return nil, err
	}

	session := &GPUSession{
		ctx:           ctx,
		originalData:  img,
		croppedData:   croppedData,
		lastCrop:      edit,
		currentEdit:   edit,
		currentPasses: passes,
	}

	return session, nil
}

// updates core adjustments in real-time
func (s *GPUSession) UpdateEdit(edit ImageEdit) (*core.ImageData, error) {
	return s.Update(edit, s.currentPasses)
}

// updates custom shader passes in real-time
func (s *GPUSession) UpdatePasses(passes []ShaderPass) (*core.ImageData, error) {
	return s.Update(s.currentEdit, passes)
}

// applies both core adjustments and custom shader passes
// Note: Changing crop settings triggers re-upload (~50-100ms overhead)
// It's recommended to simply crop before uploading the file (or waiting til after)
func (s *GPUSession) Update(edit ImageEdit, passes []ShaderPass) (*core.ImageData, error) {
	cropChanged := s.cropChanged(edit)

	if cropChanged {
		// re-upload if crop changed
		s.croppedData = CropImage(edit, s.originalData)
		if err := s.ctx.Upload(s.croppedData.Data, s.croppedData.Width, s.croppedData.Height, s.croppedData.Channels); err != nil {
			return nil, err
		}
		s.ctx.Reset()
		s.lastCrop = edit
	}

	s.currentEdit = edit
	s.currentPasses = passes

	// core adjustments (reads from original texture1)
	coreUniforms := edit.GetUniforms()
	if err := s.ctx.RunShader(shaders.PrepareForDesktop(shaders.CORE_ADJUSTMENT_SHADER), coreUniforms); err != nil {
		return nil, err
	}

	// custom shader passes
	for _, pass := range passes {
		uniforms := make(Uniforms)
		uniforms.AddUniforms(pass.Uniforms)
		uniforms.AddUniforms(coreUniforms)

		if err := s.ctx.RunShader(shaders.PrepareForDesktop(pass.Shader), uniforms); err != nil {
			return nil, err
		}
	}

	data, err := s.ctx.Download()
	if err != nil {
		return nil, err
	}

	return &core.ImageData{
		Width:      s.croppedData.Width,
		Height:     s.croppedData.Height,
		Channels:   s.croppedData.Channels,
		Colorspace: s.croppedData.Colorspace,
		Data:       data,
	}, nil
}

// checks if crop settings have changed
func (s *GPUSession) cropChanged(edit ImageEdit) bool {
	return s.lastCrop.CropTop != edit.CropTop ||
		s.lastCrop.CropLeft != edit.CropLeft ||
		s.lastCrop.CropWidth != edit.CropWidth ||
		s.lastCrop.CropHeight != edit.CropHeight
}

// releases GPU resources
func (s *GPUSession) Cleanup() {
	if s.ctx != nil {
		s.ctx.Cleanup()
		s.ctx = nil
	}
}
