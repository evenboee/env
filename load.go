package env

import (
	"fmt"
	"io/fs"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

func Must[T any](d T, err error) T {
	if err != nil {
		panic(err)
	}
	return d
}

func Guard(err error) {
	if err != nil {
		panic(err)
	}
}

var BaseLoadConfig = &LoadConfig{
	IgnoreMissingFile: true,
}

type LoadConfig struct {
	IgnoreMissingFile bool
}

func NewLoadConfig() *LoadConfig {
	return BaseLoadConfig.Copy()
}

func (c *LoadConfig) Copy() *LoadConfig {
	return &LoadConfig{
		IgnoreMissingFile: c.IgnoreMissingFile,
	}
}

func WithIgnoreMissingFile(b bool) *LoadConfig {
	return BaseLoadConfig.WithIgnoreMissingFile(b)
}

func (c *LoadConfig) WithIgnoreMissingFile(b bool) *LoadConfig {
	c.IgnoreMissingFile = b
	return c
}

func (c *LoadConfig) Load(filenames ...string) error {
	err := godotenv.Load(filenames...)
	if err != nil && c.IgnoreMissingFile {
		if _, ok := err.(*fs.PathError); ok {
			return nil
		}
	}

	return err
}

func Load(filenames ...string) (err error) {
	return BaseLoadConfig.Load(filenames...)
}

func (c *LoadConfig) MustLoad(filenames ...string) {
	Guard(c.Load(filenames...))
}

func MustLoad(filenames ...string) {
	Guard(Load(filenames...))
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
	return Must(Get[T](tg, tags...))
}

// Create and bind value
func Bind[T any](opts ...setParamsOpt) (T, error) {
	var v T
	return v, SetValue(&v, opts...)
}

// Alias for Must(Bind)
func MustBind[T any](opts ...setParamsOpt) T {
	return Must(Bind[T](opts...))
}

// Alias for Must(Get[string])
func GetString(tg string, tags ...string) string {
	s := MustGet[string](tg, tags...)
	return s
}
