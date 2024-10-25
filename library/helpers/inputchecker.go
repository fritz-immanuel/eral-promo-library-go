package helpers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/go-playground/validator/v10"
)

func MultiValueFilterCheck(input string) (string, *types.Error) {
	result := ""

	if input != "" {
		sanitizeSpace := strings.ReplaceAll(input, " ", "")
		splitComa := strings.Split(sanitizeSpace, ",")

		for _, s := range splitComa {
			num, err := strconv.Atoi(s)
			if err != nil {
				return "", &types.Error{
					Message:    "Unknown input data type",
					Error:      fmt.Errorf("unknown input data type"),
					StatusCode: http.StatusBadRequest,
				}
			}

			if reflect.TypeOf(num).String() == "int" {
				if result == "" {
					result = fmt.Sprintf("%d", num)
				} else {
					result = fmt.Sprintf("%s,%d", result, num)
				}
			}
		}
	}

	return result, nil
}

func ValidateStruct(input interface{}) *types.Error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(input)
	if errValidation != nil {
		missingTag := strings.Split(strings.Split(errValidation.Error(), ".")[1], "'")[0]
		msg := fmt.Sprintf(`'%s' is required`, missingTag)

		return &types.Error{
			Path:       ".Helpers->validateStruct()",
			Message:    msg,
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	return nil
}
