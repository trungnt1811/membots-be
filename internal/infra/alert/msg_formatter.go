package alert

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

type MsgFormatter struct {
	title string
	resp  string
}

func (f *MsgFormatter) FormatTitle(title string) *MsgFormatter {
	f.title = title
	return f
}

func (f *MsgFormatter) FormatKeyValueMsg(key string, value interface{}) *MsgFormatter {
	if valStr, ok := value.(string); ok && valStr == "" {
		return f
	}

	f.resp = f.resp + fmt.Sprintf("%v: %v\n", strcase.ToCamel(key), value)
	return f
}

func (f *MsgFormatter) FormatMsg(msg string) *MsgFormatter {
	f.resp = f.resp + fmt.Sprintf("%v\n", msg)
	return f
}

func (f *MsgFormatter) String() string {
	resp := fmt.Sprintf("========== [%v] ==========\n", strings.ToUpper(f.title))
	resp = resp + f.resp
	return resp
}
