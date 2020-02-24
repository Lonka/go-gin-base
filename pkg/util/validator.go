package util

import (
	"fmt"
	"go_gin_base/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func processErr(err error, name string) (bool, *models.ValidField) {
	vf := &models.ValidField{}
	if err != nil {
		vf.Field = name
		for _, err := range err.(validator.ValidationErrors) {
			vf.Tag = err.Tag()
			vf.Param = err.Param()
		}
		return false, vf
	}
	return true, nil
}

func HasValue(field interface{}, name string) (bool, *models.ValidField) {
	err := validate.Var(field, "required")
	return processErr(err, name)
}

func IsLte(field interface{}, name string, param string) (bool, *models.ValidField) {
	err := validate.Var(field, fmt.Sprintf("lte=%s", param))
	return processErr(err, name)
}

func HasMaxOf(field interface{}, name string, param string) (bool, *models.ValidField) {
	err := validate.Var(field, fmt.Sprintf("max=%s", param))
	return processErr(err, name)
}

func HasMinOf(field interface{}, name string, param string) (bool, *models.ValidField) {
	err := validate.Var(field, fmt.Sprintf("min=%s", param))
	return processErr(err, name)
}

func Contains(field string, name string, param string) (bool, *models.ValidField) {
	vf := &models.ValidField{}
	if ok := strings.Contains(param, field); !ok {
		vf.Field = name
		vf.Tag = "in"
		vf.Param = param
		return false, vf
	}
	return true, nil
}

func ValidStruct(s interface{}) (bool, *[]models.ValidField) {
	err := validate.Struct(s)
	if err != nil {
		var vfs []models.ValidField
		for _, e := range err.(validator.ValidationErrors) {
			vfs = append(vfs, models.ValidField{
				Field: e.Field(),
				Tag:   e.Tag(),
				Param: e.Param(),
			})
		}
		return false, &vfs

	} else {
		return true, nil
	}
}
