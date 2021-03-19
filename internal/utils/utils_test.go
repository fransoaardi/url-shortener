package utils

import "testing"

func TestGenerateID(t *testing.T) {
	type args struct {
		rawURL string
		base   int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				rawURL: "https://www.google.com/",
				base:   0,
			},
			want:    "fNgvnA6",
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GenerateID(tt.args.rawURL, tt.args.base)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GenerateID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
