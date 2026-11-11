package resizer

import (
	"reflect"
	"testing"
)

func TestParseGeometry(t *testing.T) {
	tests := []struct {
		geo  string
		want Geometry
	}{
		{"w=200", Geometry{Width: 200}},
		{"h=100", Geometry{Height: 100}},
		{"w=200,h=100,q=70", Geometry{Width: 200, Height: 100, Quality: 70}},
		{"w=200,h=100,c=true", Geometry{Width: 200, Height: 100, NeedsAutoCrop: true}},
		{"w=100,h=100,mc=true,wo=10,ho=20,cw=50,ch=60,aw=500", Geometry{
			Width: 100, Height: 100, NeedsManualCrop: true,
			CropWidthOffset: 10, CropHeightOffset: 20, CropWidth: 50, CropHeight: 60, AssumptionWidth: 500,
		}},
		{"o=true", Geometry{NeedsOriginalImage: true}},
		{"m=true", Geometry{MiddleImageSize: "1000"}},
		{"m=2000", Geometry{MiddleImageSize: "2000"}},
	}

	for _, tt := range tests {
		got, err := ParseGeometry(tt.geo)
		if err != nil {
			t.Errorf("ParseGeometry(%q) returned error: %v", tt.geo, err)
			continue
		}
		if !reflect.DeepEqual(*got, tt.want) {
			t.Errorf("ParseGeometry(%q) = %+v, want %+v", tt.geo, *got, tt.want)
		}
	}
}

func TestParseGeometryInvalid(t *testing.T) {
	tests := []string{
		"",                    // no key=value
		"w",                   // no value
		"w=abc",               // not numeric
		"h=100,w=200",         // wrong order
		"w=200,q=101",         // quality over max
		"w=200,q=-1",          // quality under min
		"w=200,c=false",       // c must be true
		"q=80",                // no width or height
		"m=500",               // invalid middle size
		"w=100,h=100,mc=true", // manual crop without crop params
	}

	for _, geo := range tests {
		if got, err := ParseGeometry(geo); err == nil {
			t.Errorf("ParseGeometry(%q) = %+v, want error", geo, *got)
		}
	}
}

func TestGeometryResizeMode(t *testing.T) {
	tests := []struct {
		geo  Geometry
		want int
	}{
		{Geometry{Width: 100}, NORMAL_RESIZE},
		{Geometry{Width: 100, NeedsAutoCrop: true}, AUTO_CROP},
		{Geometry{NeedsOriginalImage: true}, ORIGINAL},
	}

	for _, tt := range tests {
		if got := tt.geo.ResizeMode(); got != tt.want {
			t.Errorf("%+v.ResizeMode() = %d, want %d", tt.geo, got, tt.want)
		}
	}
}
