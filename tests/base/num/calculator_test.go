package num_test

import (
	"math"
	"testing"

	"github.com/bizvip/go-utils/base/num"
)

func TestNewCalculator(t *testing.T) {
	calc := num.NewCalculator()
	if calc == nil {
		t.Error("NewCalculator() returned nil")
	}
}

func TestCalculator_cleanExpression(t *testing.T) {
	calc := num.NewCalculator()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"remove spaces", "1 + 2 * 3", "1+2*3"},
		{"replace ×", "2×3", "2*3"},
		{"replace ÷", "6÷2", "6/2"},
		{"mixed operators", "1 + 2 × 3 ÷ 4", "1+2*3/4"},
		{"no changes needed", "1+2*3", "1+2*3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: cleanExpression is not exported, so we test indirectly through Evaluate
			result, _ := calc.Evaluate(tt.input)
			expectedResult, _ := calc.Evaluate(tt.want)
			if math.Abs(result-expectedResult) > 1e-9 {
				t.Errorf("cleanExpression() processing failed for %s", tt.input)
			}
		})
	}
}

func TestCalculator_Evaluate(t *testing.T) {
	calc := num.NewCalculator()
	tests := []struct {
		name       string
		expression string
		want       float64
		wantErr    bool
	}{
		{"simple addition", "2+3", 5.0, false},
		{"simple subtraction", "10-4", 6.0, false},
		{"simple multiplication", "3*4", 12.0, false},
		{"simple division", "15/3", 5.0, false},
		{"decimal numbers", "2.5+1.5", 4.0, false},
		{"negative numbers", "-5+3", -2.0, false},
		{"parentheses", "(2+3)*4", 20.0, false},
		{"nested parentheses", "((2+3)*4)/2", 10.0, false},
		{"complex expression", "2+3*4-1", 13.0, false},
		{"division with remainder", "7/2", 3.5, false},
		{"unary plus", "+5", 5.0, false},
		{"unary minus", "-(-5)", 5.0, false},
		{"spaces in expression", " 2 + 3 ", 5.0, false},
		{"unicode operators", "6÷2×3", 9.0, false},
		{"division by zero", "5/0", 0.0, true},
		{"invalid expression", "2++3", 0.0, true},
		{"empty expression", "", 0.0, true},
		{"letters in expression", "2+a", 0.0, true},
		{"mismatched parentheses", "(2+3", 0.0, true},
		{"multiple operations", "1+2*3-4/2", 5.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Evaluate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculator_formatNumber(t *testing.T) {
	// Test formatNumber indirectly through EvaluateToString
	tests := []struct {
		name       string
		expression string
		want       string
	}{
		{"integer result", "42", "42"},
		{"negative integer", "-42", "-42"},
		{"decimal result", "3.14159", "3.14159"},
		{"zero", "0", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.EvaluateToString(tt.expression)
			if err != nil {
				t.Errorf("EvaluateToString() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
		wantErr    bool
	}{
		{"global function test", "2+2", 4.0, false},
		{"complex expression", "10*(2+3)", 50.0, false},
		{"invalid expression", "abc", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.Evaluate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateToString(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       string
		wantErr    bool
	}{
		{"integer result", "4+4", "8", false},
		{"decimal result", "7/2", "3.5", false},
		{"complex calculation", "(5+3)*2", "16", false},
		{"negative result", "3-5", "-2", false},
		{"parentheses", "(10-5)*3", "15", false},
		{"division by zero", "1/0", "", true},
		{"invalid expression", "2++2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.EvaluateToString(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCalculator_Evaluate(b *testing.B) {
	calc := num.NewCalculator()
	expressions := []string{
		"2+3",
		"(2+3)*4",
		"10*(5-2)/3",
		"((2+3)*4)/2+1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := expressions[i%len(expressions)]
		_, _ = calc.Evaluate(expr)
	}
}

func BenchmarkEvaluateToString(b *testing.B) {
	expressions := []string{
		"2+3",
		"(2+3)*4",
		"10*(5-2)/3",
		"((2+3)*4)/2+1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := expressions[i%len(expressions)]
		_, _ = num.EvaluateToString(expr)
	}
}
