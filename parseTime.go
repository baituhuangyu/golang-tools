package common

import (
	"regexp"
	"fmt"
	"strconv"
	"strings"
)
// 中文数字
var chineseNums  = `〇一二三四五六七八九十`
// 纯数字格式
var patternNumberTime, _ = regexp.Compile(`\d{8}`)
// 标准时间格式的正则
var patternStandardTime, _ = regexp.Compile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
// 中文时间格式的正则
var patternChineseTime, _ = regexp.Compile(fmt.Sprintf("([%s]{4})年([%s]+)月([%s]+)日", chineseNums, chineseNums, chineseNums))
// 其它时间格式的正则
var patternGeneralTime, _ = regexp.Compile(`(\d{4})[^\d](\d{1,2})[^\d](\d{1,2})[^\d]?(\d{1,2})?:?(\d{1,2})?:?(\d{1,2})?`)


func noErrorExp(repStr string) *regexp.Regexp {
	reg, err := regexp.Compile(repStr)
	if err != nil{
		fmt.Println(err)
	}
	return reg
}

// 根据解析出的时间字符串关键字计算标准时间表示格式的字符串
func buildStrFromTimeItems(items []string) string{
	var tmps = []int{}
	if len(items) == 0{
		return ""
	}
	for _, v := range items{
		if v != ""{
			vi, _ := strconv.Atoi(v)
			tmps = append(tmps, vi)
		}
	}
	for range items[len(tmps):] {
		tmps = append(tmps, 0)
	}
	return fmt.Sprintf(`%d-%02d-%02d %02d:%02d:%02d`, tmps[0], tmps[1], tmps[2], tmps[3], tmps[4], tmps[5])
}

// 根据数字字符串关键字计算标准时间表示格式的字符串
func buildStrFromNum(timeStr string) string {
	year, _ := strconv.Atoi(timeStr[:4])
	month, _ := strconv.Atoi(timeStr[4:6])
	day, _ := strconv.Atoi(timeStr[6:])
	if 1 <= month && month <= 12 && (1 <= day && day<= 31){
		return fmt.Sprintf("%04d-%02d-%02d 00:00:00", year, month, day)
	}
	return ""
}

// 根据解析出的中文时间字符串的关键字返回对应的标准格式字符串
func buildStrFromChinese(chineseItems []string) string {
	year, month, day := cnToNum(chineseItems[0]), cnToNum(chineseItems[1]), cnToNum(chineseItems[2])
	return fmt.Sprintf("%04d-%02d-%02d 00:00:00", year, month, day)
}

// 中文数字转阿拉伯数字
func cnToNum(cnStrs string) int {
	var rst = 0
	for _, v := range cnStrs{
		index := strings.Index(chineseNums, string(v))/3
		if index >= 10{
			rst += index/10
		} else{
			rst = rst*10 + index
		}
	}
	return rst
}


func parseTime(timeStr string) string {
	// 8位数字
	if patternNumberTime.MatchString(timeStr){
		r := buildStrFromNum(patternNumberTime.FindString(timeStr))
		if r != ""{
			return r
		}
	}

	// 如果包含标准字符串形式
	if patternStandardTime.MatchString(timeStr){
		return patternStandardTime.FindString(timeStr)
	}

	// 如果是时间表达式（3个月以前，2年后...） todo

	// 中文
	timeStr = strings.Replace(timeStr, "零", "〇", -1)
	if patternChineseTime.MatchString(timeStr){
		return buildStrFromChinese(patternChineseTime.FindStringSubmatch(timeStr)[1:])
	}

	// 去除杂质
	timeStr = func(rawStr string) string{
		for _, v := range [][]string{{`日|号`, ""}, {`点|时|分`, ":"}}{
			rawStr = noErrorExp(v[0]).ReplaceAllLiteralString(rawStr, v[1])
		}
		return rawStr
	}(timeStr)

	// 查找时间
	if patternGeneralTime.MatchString(timeStr){
		return buildStrFromTimeItems(patternGeneralTime.FindStringSubmatch(timeStr)[1:])
	}

	return ""
}

func ParseTime(timeStr string) string {
	// check, change encode todo

	return parseTime(timeStr)
}