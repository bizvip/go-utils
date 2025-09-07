package num

import (
	"errors"
	"math"
	"math/big"
	"regexp"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/sqids/sqids-go"
)

// 定义常见错误消息
var (
	ErrInvalidSecPwdLength      = errors.New("security password must be 6 digits")
	ErrInvalidSecPwdConsecutive = errors.New("security password must not contain three consecutive numbers")
	ErrEmptyID                  = errors.New("id is empty")
	ErrInvalidHashID            = errors.New("decoded hash id is invalid")
	ErrParseFloat               = errors.New("unable to parse string to float")
	ErrDecimalConversion        = errors.New("unable to convert string to decimal")
)

// ValidateSecPwd 验证交易密码
func ValidateSecPwd(secPwd string) error {
	matched, err := regexp.MatchString(`^[0-9]{6}$`, secPwd)
	if err != nil || !matched {
		return ErrInvalidSecPwdLength
	}
	// 禁止三个连续数字
	for i := 0; i < len(secPwd)-2; i++ {
		if secPwd[i+1] == secPwd[i]+1 && secPwd[i+2] == secPwd[i]+2 {
			return ErrInvalidSecPwdConsecutive
		}
	}
	return nil
}

// Int64ToHashId 将 int64 转换为 Sqids 编码的 hash ID
func Int64ToHashId(number int64, minLen uint8) string {
	var ids []uint64
	ids = append(ids, uint64(number))
	s, _ := sqids.New(sqids.Options{MinLength: minLen})
	id, _ := s.Encode(ids)
	return id
}

// HashIdToInt64 将 hash ID 转换为 int64
func HashIdToInt64(id string, minLen uint8) (int64, error) {
	if id == "" {
		return 0, ErrEmptyID
	}
	s, _ := sqids.New(sqids.Options{MinLength: minLen})
	u := s.Decode(id)
	if len(u) == 0 {
		return 0, ErrInvalidHashID
	}
	return int64(u[0]), nil
}

// MergeToDecimal 如果输入的number是100000，dec是10，那么：将100000的小数点向左移动10位，得到的结果是0.00001
func MergeToDecimal(number *big.Int, dec int) decimal.Decimal {
	decimalNumber := decimal.NewFromBigInt(number, 0)
	divisor := decimal.NewFromFloat(math.Pow(10, float64(dec)))
	result := decimalNumber.Div(divisor)
	return result
}

// FormatNumStrToDecimalAndShift 输入1000，4 ，那么会输出 0.1
func FormatNumStrToDecimalAndShift(number string, decimals uint) decimal.Decimal {
	a, err := decimal.NewFromString(number)
	if err != nil {
		panic(err)
	}
	a = a.Shift(-int32(decimals))
	return a
}

// CheckNumStrInRange 检查一个字符串数字，大小是否在指定的范围内
func CheckNumStrInRange(s string, min float64, max float64) (bool, error) {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false, ErrParseFloat
	}

	return num >= min && num < max, nil
}

// StrToDecimalTruncate 将字符串数字变成decimal类型，保留指定小数位数，多余的全部放弃，不做四舍五入
func StrToDecimalTruncate(s string, precision int32) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, ErrDecimalConversion
	}
	return d.Truncate(precision), nil
}

// DecimalFormatBanker 使用银行家舍入法格式化decimal类型值为两位小数
func DecimalFormatBanker(value decimal.Decimal) string {
	valueFixed := value.StringFixedBank(2)
	if value.Mod(decimal.NewFromInt(1)).IsZero() {
		return value.Truncate(0).String()
	}
	return valueFixed
}

// GetMaxNum 返回整数切片中的最大值
func GetMaxNum(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}

	maxVal := vals[0]
	for _, v := range vals {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// Calc  输入字符串数学表达式，将计算出结果
func Calc(exp string) (string, error) {
	return EvaluateToString(exp)
}
