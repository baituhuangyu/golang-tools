package common

import (
	"regexp"
	"strings"
)

// 只保留汉子、日文、韩文、平假名（日语）、片假名（日语）字符大小写数字全角半角括号 全角半角括号
var patternName, _ = regexp.Compile(`[^[:alnum:]()（）\p{Han}\p{Hangul}\p{Hiragana}\p{Katakana}０-９Ａ-Ｚａ-ｚ 　]`)

// 全角转半角
func fullToHalfWidth(rawStr string) string {
	var rst = ""
	for _, ir := range rawStr {
		// 全角空格直接转换
		if ir == 12288{
			ir = 32
			//全角字符（除空格）根据关系转化
		}else {
			if ir <= 65374 && ir >= 65281{
				ir -= 65248
			}
		}
		rst = rst + string(ir)
	}
	return rst
}

// 清洗公司名字
func CleanName(rawStr string) string {
	if rawStr == "" {
		return rawStr
	}
	// 全角转半角
	rstStr :=fullToHalfWidth(rawStr)
	// 去除不能识别的
	rstStr = patternName.ReplaceAllLiteralString(rstStr, "")
	// 括号全部转为中文括号
	rstStr = strings.Replace(rstStr, "(", "（", -1)
	rstStr = strings.Replace(rstStr, ")", "）", -1)
	// 去除两端的空格
	return strings.TrimSpace(rstStr)
}

