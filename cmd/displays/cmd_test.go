package displays

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

func Test_parseXrandrOutput(t *testing.T) {
	type args struct {
		ctx    context.Context
		stdout bytes.Buffer
	}
	trueVal := true
	refresh60 := int64(60)
	tests := []struct {
		name    string
		args    args
		want    DisplayList
		wantErr bool
	}{
		{
			name: "Single connected display with primary and mode",
			args: args{
				ctx: context.TODO(),
				stdout: *bytes.NewBufferString(`
HDMI-1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis)
   1920x1080     60.00*+  59.94  
`),
			},
			want: DisplayList{
				{
					Name:      "HDMI-1",
					Connected: true,
					Primary:   &trueVal,
					Resolution: &Resolution{
						Width:  1920,
						Height: 1080,
					},
					Position: &Position{
						X: 0,
						Y: 0,
					},
					RefreshRate: &refresh60,
					Modes: []Mode{
						{
							Width:        1920,
							Height:       1080,
							RefreshRates: []float64{60.00, 59.94},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Disconnected display",
			args: args{
				ctx: context.TODO(),
				stdout: *bytes.NewBufferString(`
HDMI-2 disconnected (normal left inverted right x axis y axis)
`),
			},
			want: DisplayList{
				{
					Name:      "HDMI-2",
					Connected: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple displays connected",
			args: args{
				ctx: context.TODO(),
				stdout: *bytes.NewBufferString(`
DP-1 connected primary 2560x1440+0+0 (normal left inverted right x axis y axis)
   2560x1440     60.00*+ 59.95  
HDMI-1 connected 1920x1080+2560+0 (normal left inverted right x axis y axis)
   1920x1080     60.00*+ 59.94  
`),
			},
			want: DisplayList{
				{
					Name:      "DP-1",
					Connected: true,
					Primary:   &trueVal,
					Resolution: &Resolution{
						Width:  2560,
						Height: 1440,
					},
					Position: &Position{
						X: 0,
						Y: 0,
					},
					RefreshRate: &refresh60,
					Modes: []Mode{
						{
							Width:        2560,
							Height:       1440,
							RefreshRates: []float64{60.00, 59.95},
						},
					},
				},
				{
					Name:      "HDMI-1",
					Connected: true,
					Resolution: &Resolution{
						Width:  1920,
						Height: 1080,
					},
					Position: &Position{
						X: 2560,
						Y: 0,
					},
					RefreshRate: &refresh60,
					Modes: []Mode{
						{
							Width:        1920,
							Height:       1080,
							RefreshRates: []float64{60.00, 59.94},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseXrandrOutput(tt.args.ctx, tt.args.stdout)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXrandrOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseXrandrOutput() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
