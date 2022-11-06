package validation

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Singleton used for caching struct tags
var validate *validator.Validate

type ValidationErrors = validator.ValidationErrors

var (
	IsDigestString50 = regexp.MustCompile(`^\d{1,50}$`)
	IsDigestString20 = regexp.MustCompile(`^\d{1,20}$`)
	IsUUID           = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
)

func Struct(s interface{}) error {
	if validate == nil {
		validate = validator.New()

		validate.RegisterTagNameFunc(
			func(fld reflect.StructField) string {
				n := 2
				name := strings.SplitN(fld.Tag.Get("json"), ",", n)[0]

				if name == "-" {
					return ""
				}

				return name
			},
		)
	}

	return validate.Struct(s)
}
