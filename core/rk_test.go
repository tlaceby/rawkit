package core

import "testing"

func TestLibrawVersion(t *testing.T) {
	v := LibrawVersion()
	if v == "" {
		t.Fatal("expected version string, got empty")
	}
	t.Log("libraw version:", v)
}

func TestLoadingFiles(t *testing.T) {
	filepaths := []struct {
		RAW  bool
		Path string
	}{
		{
			RAW:  true,
			Path: "../testdata/_VED1070.ARW",
		},
		{
			RAW:  true,
			Path: "../testdata/_VED1242.ARW",
		},
		{
			RAW:  false,
			Path: "../testdata/_VED1242.jpg",
		},
		{
			RAW:  false,
			Path: "../testdata/_YEL4145.jpg",
		},
	}

	for _, tc := range filepaths {
		t.Run(tc.Path, func(t *testing.T) {
			img, err := ReadAll(tc.Path)
			if err != nil {
				t.Errorf("failed to load %s: %v", tc.Path, err)
				return
			}

			if img == nil {
				t.Errorf("got nil image for %s", tc.Path)
				return
			}

			if img.IsRaw() && img.Meta.ISO > 0 && len(img.Meta.CameraMake) > 0 {

			}

			// Log success info

			t.Logf("[raw=%v] (%s %s %s) loaded\n", img.IsRaw(), img.Path, img.Type, img.RawType)
		})
	}
}
