package num_test

import (
	"math/big"
	"testing"

	"github.com/bizvip/go-utils/base/num"
	"github.com/shopspring/decimal"
)

func TestValidateSecPwd(t *testing.T) {
	tests := []struct {
		name    string
		secPwd  string
		wantErr error
	}{
		{"valid password", "135792", nil},
		{"valid password 2", "987654", nil},
		{"too short", "12345", num.ErrInvalidSecPwdLength},
		{"too long", "1234567", num.ErrInvalidSecPwdLength},
		{"contains letters", "12345a", num.ErrInvalidSecPwdLength},
		{"consecutive numbers", "123456", num.ErrInvalidSecPwdConsecutive}, // 123, 234, 345, 456 are consecutive
		{"consecutive numbers 2", "234567", num.ErrInvalidSecPwdConsecutive},
		{"consecutive numbers 3", "567890", num.ErrInvalidSecPwdConsecutive},
		{"non-consecutive", "135792", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := num.ValidateSecPwd(tt.secPwd)
			if err != tt.wantErr {
				t.Errorf("ValidateSecPwd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt64ToHashId(t *testing.T) {
	tests := []struct {
		name   string
		number int64
		minLen uint8
		want   string
	}{
		{"basic test", 12345, 6, "jR5aW0"},
		{"zero", 0, 4, "gY2z"},
		{"negative", -123, 5, "0XQ7j"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.Int64ToHashId(tt.number, tt.minLen)
			if len(got) < int(tt.minLen) {
				t.Errorf("Int64ToHashId() length = %d, want >= %d", len(got), tt.minLen)
			}
		})
	}
}

func TestHashIdToInt64(t *testing.T) {
	// First generate a valid ID to test round-trip
	testNumber := int64(12345)
	testMinLen := uint8(6)
	validId := num.Int64ToHashId(testNumber, testMinLen)

	tests := []struct {
		name    string
		id      string
		minLen  uint8
		want    int64
		wantErr error
	}{
		{"valid id", validId, testMinLen, testNumber, nil},
		{"empty id", "", 4, 0, num.ErrEmptyID},
		{"truly invalid id", "!@#$%", 4, 0, num.ErrInvalidHashID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.HashIdToInt64(tt.id, tt.minLen)
			if err != tt.wantErr {
				t.Errorf("HashIdToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got != tt.want {
				t.Errorf("HashIdToInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeToDecimal(t *testing.T) {
	tests := []struct {
		name   string
		number *big.Int
		dec    int
		want   string
	}{
		{"basic test", big.NewInt(100000), 10, "0.00001"},
		{"zero decimal", big.NewInt(12345), 0, "12345"},
		{"negative", big.NewInt(-100000), 5, "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.MergeToDecimal(tt.number, tt.dec)
			if got.String() != tt.want {
				t.Errorf("MergeToDecimal() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestFormatNumStrToDecimalAndShift(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		decimals uint
		want     string
	}{
		{"basic test", "1000", 4, "0.1"},
		{"zero shift", "123", 0, "123"},
		{"large shift", "1234567", 6, "1.234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.FormatNumStrToDecimalAndShift(tt.number, tt.decimals)
			if got.String() != tt.want {
				t.Errorf("FormatNumStrToDecimalAndShift() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestCheckNumStrInRange(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		min     float64
		max     float64
		want    bool
		wantErr error
	}{
		{"in range", "5.5", 1.0, 10.0, true, nil},
		{"out of range low", "0.5", 1.0, 10.0, false, nil},
		{"out of range high", "15.0", 1.0, 10.0, false, nil},
		{"invalid string", "abc", 1.0, 10.0, false, num.ErrParseFloat},
		{"edge case min", "1.0", 1.0, 10.0, true, nil},
		{"edge case max", "10.0", 1.0, 10.0, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.CheckNumStrInRange(tt.s, tt.min, tt.max)
			if err != tt.wantErr {
				t.Errorf("CheckNumStrInRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckNumStrInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToDecimalTruncate(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		precision int32
		want      string
		wantErr   error
	}{
		{"basic truncate", "123.456789", 2, "123.45", nil},
		{"no truncate needed", "123.45", 4, "123.45", nil},
		{"invalid string", "abc", 2, "0", num.ErrDecimalConversion},
		{"zero precision", "123.456", 0, "123", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.StrToDecimalTruncate(tt.s, tt.precision)
			if err != tt.wantErr {
				t.Errorf("StrToDecimalTruncate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got.String() != tt.want {
				t.Errorf("StrToDecimalTruncate() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestDecimalFormatBanker(t *testing.T) {
	tests := []struct {
		name  string
		value decimal.Decimal
		want  string
	}{
		{"integer", decimal.NewFromInt(123), "123"},
		{"two decimals", decimal.NewFromFloat(123.45), "123.45"},
		{"banker rounding", decimal.NewFromFloat(123.125), "123.12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.DecimalFormatBanker(tt.value)
			if got != tt.want {
				t.Errorf("DecimalFormatBanker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMaxNum(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want int
	}{
		{"empty slice", []int{}, 0},
		{"single value", []int{42}, 42},
		{"multiple values", []int{1, 5, 3, 9, 2}, 9},
		{"negative values", []int{-1, -5, -3}, -1},
		{"mixed values", []int{-5, 0, 10, -2}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.GetMaxNum(tt.vals...)
			if got != tt.want {
				t.Errorf("GetMaxNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		name    string
		exp     string
		want    string
		wantErr bool
	}{
		{"simple addition", "2 + 3", "5", false},
		{"simple subtraction", "10 - 4", "6", false},
		{"simple multiplication", "3 * 4", "12", false},
		{"simple division", "15 / 3", "5", false},
		{"complex expression", "(2 + 3) * 4", "20", false},
		{"decimal result", "7 / 2", "3.5", false},
		{"negative result", "3 - 5", "-2", false},
		{"division by zero", "5 / 0", "", true},
		{"invalid expression", "2 + + 3", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.Calc(tt.exp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
