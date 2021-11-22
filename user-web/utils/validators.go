package utils

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobie(f1 validator.FieldLevel) bool {
	mobile := f1.Field().String()
	//正则
	if ok, _ := regexp.MatchString(`^1(3[0-9]|5[0-3,5-9]|7[1-3,5-8]|8[0-9])\d{8}$`, mobile); ok {
		return true
	}
	return false

}
