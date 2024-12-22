package rpn

import (
	"fmt"
	"strconv"
)

func Calc(expression string) (float64, error) {
	part := extractParts(expression)
	return parseExpression(&part)
}

func extractParts(expression string) []string {
	var part []string
	current := ""

	for _, r := range expression {
		char := string(r)

		if char == " " {
			continue
		}

		if char == "+" || char == "-" || char == "*" || char == "/" || char == "(" || char == ")" {
			if current != "" {
				part = append(part, current)
				current = ""
			}
			part = append(part, char)
		} else {
			current += char
		}
	}

	if current != "" {
		part = append(part, current)
	}

	return part
}

func parseExpression(part *[]string) (float64, error) {
	result, err := parseTerm(part)
	if err != nil {
		return 0, err
	}

	for len(*part) > 0 {
		op := (*part)[0]
		if op != "+" && op != "-" {
			break
		}
		*part = (*part)[1:]

		nextTerm, err := parseTerm(part)
		if err != nil {
			return 0, err
		}

		if op == "+" {
			result += nextTerm
		} else {
			result -= nextTerm
		}
	}

	return result, nil
}

func parseTerm(part *[]string) (float64, error) {
	result, err := parseFactor(part)
	if err != nil {
		return 0, err
	}

	for len(*part) > 0 {
		op := (*part)[0]
		if op != "*" && op != "/" {
			break
		}
		*part = (*part)[1:]

		nextFactor, err := parseFactor(part)
		if err != nil {
			return 0, err
		}

		if op == "*" {
			result *= nextFactor
		} else {
			if nextFactor == 0 {
				return 0, fmt.Errorf("деление на ноль")
			}
			result /= nextFactor
		}
	}

	return result, nil
}

func parseFactor(part *[]string) (float64, error) {
	if len(*part) == 0 {
		return 0, fmt.Errorf("неожиданный конец выражения")
	}
	token := (*part)[0]
	*part = (*part)[1:]

	if token == "(" {
		result, err := parseExpression(part)
		if err != nil {
			return 0, err
		}
		if len(*part) == 0 || (*part)[0] != ")" {
			return 0, fmt.Errorf("отсутствует закрывающая скобка")
		}
		*part = (*part)[1:]
		return result, nil
	}

	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, fmt.Errorf("неверное число: %s", token)
	}

	return value, nil
}
