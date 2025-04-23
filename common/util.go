package common

import (
	"go-project/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func BindAndValidate(ctx *gin.Context, params interface{}) error {
	if err := ctx.ShouldBindJSON(params); err != nil {
		ResponseError(ctx, INVALID_PARAMS)
		return err
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		// 使用不同的变量名避免冲突
		for _, validationErr := range err.(validator.ValidationErrors) {
			logger.Logger.Error("Validation error",
				zap.String("field", validationErr.Field()),
				zap.String("tag", validationErr.Tag()),
			)
		}
		ResponseError(ctx, INVALID_PARAMS)
		return err
	}

	return nil
}
