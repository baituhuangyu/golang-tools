package common

import (
	"testing"
	"fmt"
)

func TestCleanName(t *testing.T){
	rstStr := CleanName("   ａｓｄＡＳＦ　ＤＦ!@#$%％…………＆×＃＠＃^*aaaaqq（）qbc（）.ghjcgjwerADFGG453247687０１２３４５６７８９ＡＳＦＤＦasdsdaａｓｄ	```hghfjgj我是谁窩中華民國　aaaaaaaa全角全角。。。。；；；；“＠＠￥％＆％＆……％×＃。  ")
	fmt.Println(rstStr)
}
