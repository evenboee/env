package env

import (
	"fmt"
	"reflect"
	"strings"
)

// BindValue binds a value to an environment variable.
// Since binding is based on struct tags, we mock the struct tags by using the `tags` parameter.
// The `tags` parameter is a list of strings in the format `key=value`.
// The `tg` parameter is the `env` tag e.g. "PORT,default=8080,required".
func BindValue(v any, tg string, tags ...string) error {
	tgs, err := encodeTags(append(tags, "env="+tg))
	if err != nil {
		return err
	}

	params := DefualtSetParams.Copy()
	params.tags = tgs
	params.Tag = "env"
	return params.SetValue(v)
}

// Create and bind value
func Get[T any](tg string, tags ...string) (T, error) {
	var v T
	return v, BindValue(&v, tg, tags...)
}

// Alias for Must(Get)
func MustGet[T any](tg string, tags ...string) T {
	v, err := Get[T](tg, tags...)
	if err != nil {
		panic(err)
	}
	return v
}

// Create and bind value
func Bind[T any](opts ...setParamsOpt) (T, error) {
	var v T
	return v, SetValue(&v, opts...)
}

// Alias for Must(Bind)
func MustBind[T any](opts ...setParamsOpt) T {
	v, err := Bind[T](opts...)
	if err != nil {
		panic(err)
	}
	return v
}

// Alias for Must(Get[string])
func GetString(tg string, tags ...string) string {
	s := MustGet[string](tg, tags...)
	return s
}

func encodeTags(ts []string) (reflect.StructTag, error) {
	tgs := make([]string, len(ts))

	for i, t := range ts {
		parts := strings.SplitN(t, "=", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid tag %q", t)
		}

		tgs[i] = fmt.Sprintf("%s:%q", parts[0], parts[1])
	}

	return reflect.StructTag(strings.Join(tgs, " ")), nil
}
