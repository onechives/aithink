package controllers

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/gin-gonic/gin/binding"
)

// ContextUserIDKey 用户ID在上下文中的 key
const ContextUserIDKey = "userID"
const ContextUserRoleKey = "userRole"

var trans ut.Translator

// InitTrans 初始化校验器的语言翻译器（用于参数校验错误提示）。
func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		trans, _ = uni.GetTranslator(locale)
		if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("validator engine type is not *validator.Validate")
}
