package main

import "testing"

func Test_convert(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test convert resp protocol",
			args: args{
				input: "$5\r\nAhmed\r\n",
			},
			want:    "Ahmed",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
