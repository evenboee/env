package env

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Debug = false

type StringUnmarshaler interface {
	UnmarshalString(string) error
}

type TextUnmarshaler interface {
	encoding.TextUnmarshaler
}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "unsupported type: " + e.Type.String()
}

type RequiredFieldError struct {
	Field string
}

func (e *RequiredFieldError) Error() string {
	return "required field: " + e.Field
}

type BindError struct {
	Field string
	Type  string
	Err   error
}

func (e *BindError) Error() string {
	return fmt.Sprintf("field %q (%s): %s", e.Field, e.Type, e.Err)
}

func (e *BindError) Unwrap() error {
	return e.Err
}

type SetParams struct {
	Prefix                string
	Sep                   string
	Tag                   string
	ArrSep                string
	AutoFormatMissingKeys bool

	tags reflect.StructTag
}

var DefualtSetParams = &SetParams{
	Prefix:                "",
	Sep:                   "_",
	Tag:                   "env",
	ArrSep:                ",",
	AutoFormatMissingKeys: true,
	tags:                  "",
}

func (p *SetParams) Copy() *SetParams {
	return &SetParams{
		Prefix:                p.Prefix,
		Sep:                   p.Sep,
		Tag:                   p.Tag,
		ArrSep:                p.ArrSep,
		AutoFormatMissingKeys: p.AutoFormatMissingKeys,
		tags:                  p.tags,
	}
}

type setParamsOpt func(*SetParams)

func WithPrefix(prefix string) setParamsOpt {
	return func(p *SetParams) {
		p.Prefix = prefix
	}
}

func WithSep(sep string) setParamsOpt {
	return func(p *SetParams) {
		p.Sep = sep
	}
}

func WithTag(tag string) setParamsOpt {
	return func(p *SetParams) {
		p.Tag = tag
	}
}

func WithArrSep(arrSep string) setParamsOpt {
	return func(p *SetParams) {
		p.ArrSep = arrSep
	}
}

func WithAutoFormatMissingKeys(auto bool) setParamsOpt {
	return func(p *SetParams) {
		p.AutoFormatMissingKeys = auto
	}
}

func With(f setParamsOpt) setParamsOpt {
	return func(p *SetParams) {
		f(p)
	}
}

func SetValue(obj any, opts ...setParamsOpt) error {
	params := DefualtSetParams.Copy()
	for _, opt := range opts {
		opt(params)
	}

	return params.SetValue(obj)
}

func (c *SetParams) SetValue(obj any) error {
	v := reflect.ValueOf(obj).Elem()

	err, _ := c.setValue(v, reflect.StructField{
		Tag: c.tags,
	}, c.Prefix, "")
	return err
}

func (c *SetParams) setValue(value reflect.Value, field reflect.StructField, prefix string, v string) (error, bool) {
	var (
		valS   = v
		ok     = v != ""
		tg     tag
		envKey string
		err    error
	)

	tagValue := field.Tag.Get(c.Tag)
	if tagValue == "" {
		tagValue = field.Name
	} else if tagValue == "-" {
		return nil, false
	}

	tg, err = parseTag(tagValue)
	if err != nil {
		return err, false
	}

	if tg.Name == "" && c.AutoFormatMissingKeys {
		tg.Name = formatName(field.Name)
	}

	if Debug {
		fmt.Printf("tag for %q: %+v\n", field.Name, tg)
	}

	envKey = prefix
	if tg.Name != "" {
		if envKey != "" {
			envKey += c.Sep
		}
		envKey += tg.Name
	}

	valS, ok = os.LookupEnv(envKey)
	if Debug {
		if !ok {
			fmt.Println("env key not found: ", envKey)
		} else {
			fmt.Printf("found env key: %s\n\t>> %s\n", envKey, valS)
		}
	}
	if !ok && tg.Default == "" {
		if tg.Required {
			return &RequiredFieldError{
				Field: envKey,
			}, false
		}
		if tg.SkipOnNoValue {
			if Debug {
				fmt.Printf("skipping %q\n", envKey)
			}
			return nil, false
		}
	}

	// c.ArrSep should be tag option
	val := strings.Split(valS, c.ArrSep)

	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			val = strings.Split(tg.Default, tg.DefaultValueSeparator)
		}
		return c.setSlice(value, field, tg, val), true
	case reflect.Array:
		if !ok {
			val = strings.Split(tg.Default, tg.DefaultValueSeparator)
		}
		if len(val) != value.Len() {
			return &BindError{
				Field: envKey,
				Type:  value.Type().String(),
				Err:   fmt.Errorf("%q is not a valid length for %s", val, value.Type().String()),
			}, false
		}
		return c.setArray(value, field, tg, val), true
	default:
		if !ok {
			valS = tg.Default
		}
		return c.setWithType(value, field, tg, valS)
	}
}

func (c *SetParams) setWithType(value reflect.Value, field reflect.StructField, tg tag, val string) (error, bool) {
	switch value.Addr().Interface().(type) {
	case StringUnmarshaler:
		return value.Addr().Interface().(StringUnmarshaler).UnmarshalString(val), true
	case TextUnmarshaler:
		return value.Addr().Interface().(TextUnmarshaler).UnmarshalText([]byte(val)), true
	}

	switch value.Kind() {
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value), true
		}
		return c.setStruct(val, value, field, tg), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value), true
		}
		return setIntField(val, value, field), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintField(val, value, field), true
	case reflect.Bool:
		return setBoolField(val, value, field), true
	case reflect.Float32, reflect.Float64:
		return setFloatField(val, value, field), true
	case reflect.String:
		value.SetString(val)
		return nil, true
	case reflect.Pointer:
		if !value.Elem().IsValid() {
			value.Set(reflect.New(value.Type().Elem()))
		}

		// return c.setWithType(value.Elem(), field, tg, val)
		return c.setValue(value.Elem(), field, "", val)
	}

	return &UnsupportedTypeError{Type: value.Type()}, false
}

func (c *SetParams) setStruct(val string, value reflect.Value, field reflect.StructField, tg tag) error {
	nPrefix := c.Prefix
	if tg.Name != "" {
		nPrefix += c.Sep + tg.Name
	}

	t := value.Type()

	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		fieldT := t.Field(i)
		if !fieldT.IsExported() {
			continue
		}

		fieldV := value.Field(i)

		err, _ := c.setValue(fieldV, fieldT, nPrefix, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *SetParams) setArray(value reflect.Value, field reflect.StructField, tg tag, vals []string) error {
	for i, s := range vals {
		err, _ := c.setWithType(value.Index(i), field, tg, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *SetParams) setSlice(value reflect.Value, field reflect.StructField, tg tag, vals []string) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := c.setArray(slice, field, tg, vals)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setIntField(val string, value reflect.Value, field reflect.StructField) error {
	if val == "" {
		val = "0"
	}

	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}

	value.SetInt(n)
	return nil
}

func setUintField(val string, value reflect.Value, field reflect.StructField) error {
	if val == "" {
		val = "0"
	}

	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return err
	}

	value.SetUint(n)
	return nil
}

func setBoolField(val string, value reflect.Value, field reflect.StructField) error {
	if val == "" {
		val = "false"
	}

	b, err := parseBool(val)
	if err != nil {
		return err
	}

	value.SetBool(b)
	return nil
}

func setFloatField(val string, value reflect.Value, field reflect.StructField) error {
	if val == "" {
		val = "0.0"
	}

	n, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}

	value.SetFloat(n)
	return nil
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

var DefaultTimeFormat = time.RFC3339

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = DefaultTimeFormat
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := parseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

// same as strconv.ParseBool but with extra values
func parseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Y", "y":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "N", "n":
		return false, nil
	}
	return false, &strconv.NumError{Func: "env.parseBool", Num: str, Err: strconv.ErrSyntax}
}
