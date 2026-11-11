package resizer

import "testing"

func calc(t *testing.T, imageWidth, imageHeight int, option *ResizeOption) Coodinates {
	t.Helper()
	calculator, err := NewCoodinatesCalculator(option)
	if err != nil {
		t.Fatalf("NewCoodinatesCalculator(%+v) returned error: %v", option, err)
	}
	calculator.SetImageSize(imageWidth, imageHeight)
	return *calculator.Calc(option)
}

func TestNewCoodinatesCalculatorInvalid(t *testing.T) {
	if _, err := NewCoodinatesCalculator(&ResizeOption{}); err == nil {
		t.Error("NewCoodinatesCalculator without width and height should return error")
	}
}

func TestCoodinatesCalculatorResize(t *testing.T) {
	tests := []struct {
		option *ResizeOption
		want   Coodinates
	}{
		{&ResizeOption{Width: 200}, Coodinates{ResizeWidth: 200, ResizeHeight: 100}},
		{&ResizeOption{Height: 100}, Coodinates{ResizeWidth: 200, ResizeHeight: 100}},
		{&ResizeOption{Width: 300, Height: 300}, Coodinates{ResizeWidth: 300, ResizeHeight: 150}},
	}

	for _, tt := range tests {
		if got := calc(t, 1000, 500, tt.option); got != tt.want {
			t.Errorf("Resize %+v = %+v, want %+v", tt.option, got, tt.want)
		}
	}
}

func TestCoodinatesCalculatorAutoCrop(t *testing.T) {
	tests := []struct {
		imageWidth, imageHeight int
		want                    Coodinates
	}{
		// landscape: fit height, crop width center
		{1000, 500, Coodinates{ResizeWidth: 600, ResizeHeight: 300, CropWidth: 300, CropHeight: 300, WidthOffset: 150}},
		// portrait: fit width, crop height center
		{500, 1000, Coodinates{ResizeWidth: 300, ResizeHeight: 600, CropWidth: 300, CropHeight: 300, HeightOffset: 150}},
	}

	option := &ResizeOption{Width: 300, Height: 300, NeedsAutoCrop: true}
	for _, tt := range tests {
		if got := calc(t, tt.imageWidth, tt.imageHeight, option); got != tt.want {
			t.Errorf("AutoCrop %dx%d = %+v, want %+v", tt.imageWidth, tt.imageHeight, got, tt.want)
		}
	}
}

func TestCoodinatesCalculatorManualCrop(t *testing.T) {
	option := &ResizeOption{
		Width: 100, Height: 100, NeedsManualCrop: true,
		CropWidthOffset: 10, CropHeightOffset: 20, CropWidth: 100, CropHeight: 50, AssumptionWidth: 500,
	}
	want := Coodinates{ResizeWidth: 100, ResizeHeight: 100, CropWidth: 100, CropHeight: 100, WidthOffset: 70, HeightOffset: 40}

	if got := calc(t, 1000, 500, option); got != want {
		t.Errorf("ManualCrop = %+v, want %+v", got, want)
	}
}
