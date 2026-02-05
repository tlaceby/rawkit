package processing_test

import (
	"testing"
	"time"

	"github.com/tlaceby/rawkit/core"
	"github.com/tlaceby/rawkit/processing"
	"github.com/tlaceby/rawkit/processing/shaders"
)

func TestExposureValues(t *testing.T) {
	filepaths := []struct {
		RAW  bool
		Path string
	}{
		{RAW: true, Path: "../testdata/_VED1070.ARW"},
		{RAW: true, Path: "../testdata/_VED1242.ARW"},
		{RAW: false, Path: "../testdata/_VED1242.jpg"},
		{RAW: false, Path: "../testdata/_YEL4145.jpg"},
	}

	for _, tc := range filepaths {
		t.Run(tc.Path, func(t *testing.T) {
			img, err := core.ReadAll(tc.Path)
			if err != nil {
				t.Fatalf("failed to load: %v", err)
			}

			originalAvg := calculateAverageBrightness(img.Data.Data)
			edit := processing.ImageEdit{Exposure: 85}

			ts := time.Now()
			result, err := processing.ApplyEdit(img.Data, edit)
			if err != nil {
				t.Fatalf("failed to apply edit: %v", err)
			}

			resultAvg := calculateAverageBrightness(result.Data)
			if resultAvg <= originalAvg {
				t.Errorf("brightness should increase. original: %.2f, result: %.2f", originalAvg, resultAvg)
			}

			percentIncrease := ((resultAvg - originalAvg) / originalAvg) * 100
			t.Logf("exposure +85: %.1f%% increase (%dms)", percentIncrease, time.Since(ts).Milliseconds())
		})
	}
}

func TestExposureValuesManual(t *testing.T) {
	filepaths := []struct {
		RAW  bool
		Path string
	}{
		{RAW: true, Path: "../testdata/_VED1070.ARW"},
		{RAW: true, Path: "../testdata/_VED1242.ARW"},
		{RAW: false, Path: "../testdata/_VED1242.jpg"},
		{RAW: false, Path: "../testdata/_YEL4145.jpg"},
	}

	for _, tc := range filepaths {
		t.Run(tc.Path, func(t *testing.T) {
			ts := time.Now()
			img, err := core.ReadAll(tc.Path)
			if err != nil {
				t.Fatal(err)
			}

			loadTime := time.Since(ts).Milliseconds()

			ts = time.Now()
			opengl, err := processing.CreateOpenGLContext()
			if err != nil {
				t.Fatal(err)
			}

			defer opengl.Cleanup()
			contextTime := time.Since(ts).Milliseconds()

			ts = time.Now()
			err = opengl.Upload(img.Data.Data, img.Data.Width, img.Data.Height, img.Data.Channels)
			if err != nil {
				t.Fatal(err)
			}

			uploadTime := time.Since(ts).Milliseconds()

			edit := processing.ImageEdit{Exposure: 85}
			uniforms := edit.GetUniforms()

			ts = time.Now()
			err = opengl.RunShader(shaders.PrepareForDesktop(shaders.CORE_ADJUSTMENT_SHADER), uniforms)
			if err != nil {
				t.Fatal(err)
			}

			shaderTime := time.Since(ts).Milliseconds()

			ts = time.Now()
			_, err = opengl.Download()
			if err != nil {
				t.Fatal(err)
			}

			downloadTime := time.Since(ts).Milliseconds()

			t.Logf("load: %dms, context: %dms, upload: %dms, shader: %dms, download: %dms", loadTime, contextTime, uploadTime, shaderTime, downloadTime)
		})
	}
}

func TestGPUSession(t *testing.T) {
	filepaths := []string{
		"../testdata/_VED1070.ARW",
		"../testdata/_VED1242.ARW",
		"../testdata/_VED1242.jpg",
		"../testdata/_YEL4145.jpg",
	}

	for _, filepath := range filepaths {
		t.Run(filepath, func(t *testing.T) {
			img, err := core.ReadAll(filepath)
			if err != nil {
				t.Fatal(err)
			}

			originalAvg := calculateAverageBrightness(img.Data.Data)

			ts := time.Now()
			ctx, err := processing.CreateOpenGLContext()
			if err != nil {
				t.Fatal(err)
			}

			session, err := processing.NewGPUSession(img.Data, ctx, processing.ImageEdit{}, []processing.ShaderPass{})
			if err != nil {
				t.Fatal(err)
			}

			defer session.Cleanup()

			setupTime := time.Since(ts).Milliseconds()

			exposureValues := []float32{-75, -50, -40, -30, -20, -10, 0, 10, 20, 30, 40, 50, 60, 70, 75}

			var totalUpdateTime int64
			var minTime int64 = 999999
			var maxTime int64 = 0

			for _, exposure := range exposureValues {
				ts = time.Now()
				result, err := session.UpdateEdit(processing.ImageEdit{Exposure: exposure})
				updateTime := time.Since(ts).Milliseconds()

				// t.Logf("GPUSession.Update() time: %dms\n", updateTime)

				if err != nil {
					t.Fatalf("update failed at exposure %.0f: %v", exposure, err)
				}

				totalUpdateTime += updateTime
				if updateTime < minTime {
					minTime = updateTime
				}

				if updateTime > maxTime {
					maxTime = updateTime
				}

				resultAvg := calculateAverageBrightness(result.Data)
				percentChange := ((resultAvg - originalAvg) / originalAvg) * 100

				if exposure == 0 && percentChange > 1.0 {
					t.Errorf("exposure 0 should not change brightness (got %.1f%% change)", percentChange)
				}
			}

			avgUpdateTime := totalUpdateTime / int64(len(exposureValues))
			t.Logf("setup: %dms, updates: %d, avg: %dms, min: %dms, max: %dms", setupTime, len(exposureValues), avgUpdateTime, minTime, maxTime)

			if avgUpdateTime > 50 {
				t.Logf("warning: average update time %dms exceeds 50ms target", avgUpdateTime)
			}
		})
	}
}

func TestGPUSessionWithPasses(t *testing.T) {
	filepaths := []string{
		"../testdata/_VED1070.ARW",
		"../testdata/_VED1242.ARW",
		"../testdata/_VED1242.jpg",
		"../testdata/_YEL4145.jpg",
	}

	vignetteShader := `
		in vec2 uv;
		out vec4 frag_color;
		uniform sampler2D u_texture;
		uniform float strength;
		
		void main() {
			vec4 color = texture(u_texture, uv);
			vec2 center = vec2(0.5, 0.5);
			float dist = distance(uv, center);
			float vignette = 1.0 - (dist * strength);
			vignette = clamp(vignette, 0.0, 1.0);
			color.rgb *= vignette;
			frag_color = color;
		}
	`

	tintShader := `
		in vec2 uv;
		out vec4 frag_color;
		uniform sampler2D u_texture;
		uniform vec3 tint_color;
		uniform float tint_strength;
		
		void main() {
			vec4 color = texture(u_texture, uv);
			color.rgb = mix(color.rgb, color.rgb * tint_color, tint_strength);
			frag_color = color;
		}
	`

	sharpenShader := `
		in vec2 uv;
		out vec4 frag_color;
		uniform sampler2D u_texture;
		uniform float strength;
		uniform vec2 u_resolution;
		
		void main() {
			vec2 texelSize = 1.0 / u_resolution;
			vec4 center = texture(u_texture, uv);
			vec4 top    = texture(u_texture, uv + vec2(0.0, texelSize.y));
			vec4 bottom = texture(u_texture, uv - vec2(0.0, texelSize.y));
			vec4 left   = texture(u_texture, uv - vec2(texelSize.x, 0.0));
			vec4 right  = texture(u_texture, uv + vec2(texelSize.x, 0.0));
			vec4 sharpened = center * (1.0 + 4.0 * strength) - (top + bottom + left + right) * strength;
			frag_color = clamp(sharpened, 0.0, 1.0);
		}
	`

	for _, filepath := range filepaths {
		t.Run(filepath, func(t *testing.T) {
			img, err := core.ReadAll(filepath)
			if err != nil {
				t.Fatal(err)
			}

			ctx, err := processing.CreateOpenGLContext()
			if err != nil {
				t.Fatal(err)
			}

			session, err := processing.NewGPUSession(img.Data, ctx, processing.ImageEdit{}, []processing.ShaderPass{})
			if err != nil {
				t.Fatal(err)
			}
			defer session.Cleanup()

			// baseline with no edits
			ts := time.Now()
			baseline, err := session.UpdateEdit(processing.ImageEdit{})
			if err != nil {
				t.Fatal(err)
			}
			baselineTime := time.Since(ts).Milliseconds()
			baselineAvg := calculateAverageBrightness(baseline.Data)

			// test edit only
			ts = time.Now()
			_, err = session.UpdateEdit(processing.ImageEdit{Exposure: 75})
			if err != nil {
				t.Fatal(err)
			}
			editTime := time.Since(ts).Milliseconds()

			// test vignette without exposure
			vignettePass := processing.ShaderPass{
				Name:     "vignette",
				Shader:   vignetteShader,
				Uniforms: processing.NewUniforms().AddFloat("strength", 1.2),
			}

			ts = time.Now()
			resultVignette, err := session.Update(processing.ImageEdit{}, []processing.ShaderPass{vignettePass})
			if err != nil {
				t.Fatal(err)
			}
			vignetteTime := time.Since(ts).Milliseconds()

			// test multiple shaders
			tintPass := processing.ShaderPass{
				Name:     "warm_tint",
				Shader:   tintShader,
				Uniforms: processing.NewUniforms().AddVec3("tint_color", 1.2, 1.0, 0.8).AddFloat("tint_strength", 0.3),
			}

			ts = time.Now()
			_, err = session.UpdatePasses([]processing.ShaderPass{vignettePass, tintPass})
			if err != nil {
				t.Fatal(err)
			}
			multiTime := time.Since(ts).Milliseconds()

			// test combined
			ts = time.Now()
			_, err = session.Update(
				processing.ImageEdit{Exposure: 85, Contrast: 10},
				[]processing.ShaderPass{vignettePass, tintPass},
			)
			if err != nil {
				t.Fatal(err)
			}
			combinedTime := time.Since(ts).Milliseconds()

			// test sharpen
			sharpenPass := processing.ShaderPass{
				Name:   "sharpen",
				Shader: sharpenShader,
				Uniforms: processing.NewUniforms().
					AddFloat("strength", 0.5).
					AddVec2("u_resolution", float32(img.Data.Width), float32(img.Data.Height)),
			}

			ts = time.Now()
			_, err = session.UpdatePasses([]processing.ShaderPass{sharpenPass})
			if err != nil {
				t.Fatal(err)
			}
			sharpenTime := time.Since(ts).Milliseconds()

			// verify vignette darkens image
			vignetteAvg := calculateAverageBrightness(resultVignette.Data)
			if vignetteAvg >= baselineAvg {
				t.Errorf("vignette should darken image: baseline %.2f, with vignette %.2f", baselineAvg, vignetteAvg)
			}

			brightnessReduction := ((baselineAvg - vignetteAvg) / baselineAvg) * 100

			t.Logf("baseline: %dms, edit: %dms, vignette: %dms, multi: %dms, combined: %dms, sharpen: %dms",
				baselineTime, editTime, vignetteTime, multiTime, combinedTime, sharpenTime)
			t.Logf("vignette darkened by %.1f%%", brightnessReduction)
		})
	}
}

func calculateAverageBrightness(data []uint16) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum uint64
	for _, val := range data {
		sum += uint64(val)
	}

	return float64(sum) / float64(len(data))
}
