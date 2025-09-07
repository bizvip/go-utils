/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package num

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strconv"
	"strings"
)

// Calculator 高级表达式计算器
type Calculator struct{}

// NewCalculator 创建计算器实例
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Evaluate 计算数学表达式（使用Go AST）
func (c *Calculator) Evaluate(expression string) (float64, error) {
	// 清理表达式
	expr := c.cleanExpression(expression)

	// 使用Go的AST解析器来计算表达式
	return c.evaluateWithAST(expr)
}

// cleanExpression 清理表达式
func (c *Calculator) cleanExpression(expression string) string {
	// 移除所有空格
	expr := strings.ReplaceAll(expression, " ", "")

	// 标准化运算符
	expr = strings.ReplaceAll(expr, "×", "*")
	expr = strings.ReplaceAll(expr, "÷", "/")

	return expr
}

// evaluateWithAST 使用Go AST解析和计算表达式
func (c *Calculator) evaluateWithAST(expression string) (float64, error) {
	// 解析表达式
	expr, err := parser.ParseExpr(expression)
	if err != nil {
		return 0, fmt.Errorf("无效的表达式")
	}

	// 计算表达式
	return c.evalNode(expr)
}

// evalNode 递归计算AST节点
func (c *Calculator) evalNode(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.BasicLit:
		// 基本字面量（数字）
		if n.Kind == token.INT || n.Kind == token.FLOAT {
			return strconv.ParseFloat(n.Value, 64)
		}
		return 0, fmt.Errorf("不支持的字面量类型")

	case *ast.BinaryExpr:
		// 二元表达式
		left, err := c.evalNode(n.X)
		if err != nil {
			return 0, err
		}
		right, err := c.evalNode(n.Y)
		if err != nil {
			return 0, err
		}

		switch n.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, fmt.Errorf("除零错误")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("不支持的运算符: %s", n.Op.String())
		}

	case *ast.UnaryExpr:
		// 一元表达式（负号）
		operand, err := c.evalNode(n.X)
		if err != nil {
			return 0, err
		}

		switch n.Op {
		case token.SUB:
			return -operand, nil
		case token.ADD:
			return operand, nil
		default:
			return 0, fmt.Errorf("不支持的一元运算符: %s", n.Op.String())
		}

	case *ast.ParenExpr:
		// 括号表达式
		return c.evalNode(n.X)

	default:
		return 0, fmt.Errorf("不支持的表达式类型")
	}
}

// formatNumber 格式化数字显示
func (c *Calculator) formatNumber(num float64) string {
	// 如果是整数，显示为整数
	if num == math.Floor(num) {
		return fmt.Sprintf("%.0f", num)
	}

	// 否则显示小数，最多6位小数
	return fmt.Sprintf("%.6g", num)
}

// 全局实例和便捷函数
var globalCalculator = NewCalculator()

// Evaluate 全局计算函数
func Evaluate(expression string) (float64, error) {
	return globalCalculator.Evaluate(expression)
}

// EvaluateToString 计算并返回字符串结果
func EvaluateToString(expression string) (string, error) {
	result, err := globalCalculator.Evaluate(expression)
	if err != nil {
		return "", err
	}

	return globalCalculator.formatNumber(result), nil
}
