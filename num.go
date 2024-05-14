/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/sqids/sqids-go"

	"github.com/bizvip/go-utils/logs"
)

type NumUtils struct {
}

func NewNumUtils() *NumUtils { return &NumUtils{} }

func (n *NumUtils) Int64ToHashId(number int64, minLen uint8) string {
	var ids []uint64
	ids = append(ids, uint64(number))
	s, _ := sqids.New(sqids.Options{MinLength: minLen})
	id, _ := s.Encode(ids)
	return id
}
func (n *NumUtils) HashIdToInt64(id string, minLen uint8) (int64, error) {
	if id == "" {
		return 0, errors.New("id is empty")
	}
	s, _ := sqids.New(sqids.Options{MinLength: minLen})
	u := s.Decode(id)
	return int64(u[0]), nil
}
func (n *NumUtils) RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

// MergeToDecimal 如果输入的number是100000，dec是10，那么：将100000除以10000000000 (即10的10次方) 得到的结果是0.00001
func (n *NumUtils) MergeToDecimal(number *big.Int, dec int) decimal.Decimal {
	decimalNumber := decimal.NewFromBigInt(number, 0)
	divisor := decimal.NewFromFloat(math.Pow(10, float64(dec)))
	result := decimalNumber.Div(divisor)
	return result
}

// FormatNumStrToDecimalAndShift 输入1000，4 ，那么会输出 0.1
func (n *NumUtils) FormatNumStrToDecimalAndShift(number string, decimals uint) decimal.Decimal {
	a, e := decimal.NewFromString(number)
	if e != nil {
		panic(e)
	}
	a = a.Shift(-int32(decimals))
	return a
}

// CheckNumStrInRange 检查一个字符串数字，大小是否在指定的范围内
func (n *NumUtils) CheckNumStrInRange(s string, min float64, max float64) (bool, error) {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false, err
	}

	return num >= min && num < max, nil
}

// StrToDecimalTruncate 将字符串数字变成decimal类型，保留指定小数位数，多余的全部放弃，不做四舍五入
func (n *NumUtils) StrToDecimalTruncate(s string, precision int32) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		logs.Logger().Error(err)
		return decimal.Zero
	}
	return d.Truncate(precision)
}

// DecimalFormatBanker 使用银行家舍入法格式化decimal类型值为两位小数
func (n *NumUtils) DecimalFormatBanker(value decimal.Decimal) string {
	valueFixed := value.StringFixedBank(2)
	if value.Mod(decimal.NewFromInt(1)).IsZero() {
		return value.Truncate(0).String()
	}
	return valueFixed
}

func (n *NumUtils) GetMaxNum(vals ...int) int {
	maxVal := vals[0]
	for _, v := range vals {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}
