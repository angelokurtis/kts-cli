package images

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

// Test_parseOutput
func Test_parseOutput(t *testing.T) {
	type args struct {
		ctx context.Context
		buf bytes.Buffer
	}

	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "single group of duplicates",
			args: args{
				ctx: context.Background(),
				buf: *bytes.NewBufferString("image1.jpg image2.jpg\n"),
			},
			want: [][]string{{"image1.jpg", "image2.jpg"}},
		},
		{
			name: "multiple groups of duplicates",
			args: args{
				ctx: context.Background(),
				buf: *bytes.NewBufferString("image1.jpg image2.jpg\nimage3.jpg image4.jpg\n"),
			},
			want: [][]string{
				{"image1.jpg", "image2.jpg"},
				{"image3.jpg", "image4.jpg"},
			},
		},
		{
			name: "empty output",
			args: args{
				ctx: context.Background(),
				buf: *bytes.NewBufferString(""),
			},
			want: [][]string{},
		},
		{
			name: "output with empty lines and spaces",
			args: args{
				ctx: context.Background(),
				buf: *bytes.NewBufferString("\n  image1.jpg   image2.jpg  \n\nimage3.jpg image4.jpg\n"),
			},
			want: [][]string{
				{"image1.jpg", "image2.jpg"},
				{"image3.jpg", "image4.jpg"},
			},
		},
		{
			name: "single file no duplicates",
			args: args{
				ctx: context.Background(),
				buf: *bytes.NewBufferString("image1.jpg\n"),
			},
			want: [][]string{
				{"image1.jpg"},
			},
		},
		{
			name: "realistic full file paths with spaces and parentheses",
			args: args{
				ctx: context.Background(),
				buf: *bytes.NewBufferString(
					`/home/kurtis/Pictures/2025-05-04_18-24_1.png /home/kurtis/Pictures/2025-05-04_18-24_1 (Copy).png
/home/kurtis/Pictures/2025-05-04_18-24 (Copy).png /home/kurtis/Pictures/2025-05-04_18-24.png /home/kurtis/Pictures/2025-05-04_18-24 (Copy 2).png
`),
			},
			want: [][]string{
				{
					"/home/kurtis/Pictures/2025-05-04_18-24_1.png",
					"/home/kurtis/Pictures/2025-05-04_18-24_1 (Copy).png",
				},
				{
					"/home/kurtis/Pictures/2025-05-04_18-24 (Copy).png",
					"/home/kurtis/Pictures/2025-05-04_18-24.png",
					"/home/kurtis/Pictures/2025-05-04_18-24 (Copy 2).png",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOutput(tt.args.ctx, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				got = [][]string{}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOutput() got = %v, want %v", got, tt.want)
			}
		})
	}
}
