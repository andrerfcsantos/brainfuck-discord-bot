package brainfuck

import "testing"

func TestShorten(t *testing.T) {
	type args struct {
		program string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "teste1",
			args: args{
				program: "------",
			},
			want: "-6",
		},
		{
			name: "teste1",
			args: args{
				program: ".,---[]",
			},
			want: ".,-3[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Shorten(tt.args.program); got != tt.want {
				t.Errorf("Shorten() = %v, want %v", got, tt.want)
			}
		})
	}
}
