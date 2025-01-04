package core

import (
	"testing"
)

func TestValidateConstraint(t *testing.T) {
	svc := New(nil)

	tests := []struct {
		name       string
		code       string
		constraint Constraint
		want       bool
	}{
		{
			name: "max constraint valid",
			code: "short code",
			constraint: Constraint{
				Type:  "max",
				Value: 20,
			},
			want: true,
		},
		{
			name: "max constraint invalid",
			code: "this is a very long piece of code",
			constraint: Constraint{
				Type:  "max",
				Value: 10,
			},
			want: false,
		},
		{
			name: "regex constraint valid",
			code: "func NewUser()",
			constraint: Constraint{
				Type:  "regex",
				Value: "^func New",
			},
			want: true,
		},
		{
			name: "regex constraint invalid",
			code: "func CreateUser()",
			constraint: Constraint{
				Type:  "regex",
				Value: "^func New",
			},
			want: false,
		},
		{
			name: "forbidden constraint valid",
			code: "good code",
			constraint: Constraint{
				Type:  "forbidden",
				Value: []string{"bad", "wrong"},
			},
			want: true,
		},
		{
			name: "forbidden constraint invalid",
			code: "bad code",
			constraint: Constraint{
				Type:  "forbidden",
				Value: []string{"bad", "wrong"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.validateConstraint(tt.code, tt.constraint)
			if got != tt.want {
				t.Errorf("validateConstraint() = %v, want %v", got, tt.want)
			}
		})
	}
}
