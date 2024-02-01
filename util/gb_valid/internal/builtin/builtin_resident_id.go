package builtin

import (
	"errors"
	gbregex "github.com/Ghostbb-io/gb/text/gb_regex"
	"strconv"
	"strings"
)

// RuleResidentId implements `resident-id` rule:
// Resident id number.
//
// Format: resident-id
type RuleResidentId struct{}

func init() {
	Register(RuleResidentId{})
}

func (r RuleResidentId) Name() string {
	return "resident-id"
}

func (r RuleResidentId) Message() string {
	return "The {field} value `{value}` is not a valid resident id number"
}

func (r RuleResidentId) Run(in RunInput) error {
	if r.checkResidentId(in.Value.String()) {
		return nil
	}
	return errors.New(in.Message)
}

// checkResidentId checks whether given id a china resident id number.
//
// xxxxxx yyyy MM dd 375 0  十八位
// xxxxxx   yy MM dd  75 0  十五位
//
// 地区：     [1-9]\d{5}
// 年的前两位：(18|19|([23]\d))  1800-2399
// 年的后两位：\d{2}
// 月份：     ((0[1-9])|(10|11|12))
// 天数：     (([0-2][1-9])|10|20|30|31) 闰年不能禁止29+
//
// 三位顺序码：\d{3}
// 两位顺序码：\d{2}
// 校验码：   [0-9Xx]
//
// 十八位：^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$
// 十五位：^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$
//
// 总：
// (^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)
func (r RuleResidentId) checkResidentId(id string) bool {
	id = strings.ToUpper(strings.TrimSpace(id))
	if len(id) != 18 {
		return false
	}
	var (
		weightFactor = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		checkCode    = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
		last         = id[17]
		num          = 0
	)
	for i := 0; i < 17; i++ {
		tmp, err := strconv.Atoi(string(id[i]))
		if err != nil {
			return false
		}
		num = num + tmp*weightFactor[i]
	}
	if checkCode[num%11] != last {
		return false
	}

	return gbregex.IsMatchString(
		`(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)`,
		id,
	)
}
