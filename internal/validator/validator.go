package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = map[string]string{}
	}

	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = message
	}
}

func (v *Validator) CheckError(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(n int, permitted ...int) bool {
	for i := range permitted {
		if n == permitted[i] {
			return true
		}
	}
	return false
}
