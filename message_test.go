package hub

import "testing"

func TestFields_String(t *testing.T) {
	tests := []struct {
		name string
		f    Fields
		want string
	}{
		{"empty fields", Fields{}, "Fields(<empty>)"},
		{"one field", Fields{"foo": "bar"}, "Fields( [foo]bar )"},
		{"two fields", Fields{"foo": "bar", "baz": true}, "Fields( [baz]true, [foo]bar )"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("Fields.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
