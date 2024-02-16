package env

import (
	"fmt"
	"strings"
)

type UnknownTagPartError struct {
	Part     string
	TagName  string
	TagValue string
}

func (e *UnknownTagPartError) Error() string {
	return fmt.Sprintf("unknown tag part %q in tag %q with value %q", e.Part, e.TagName, e.TagValue)
}

var DefaultValueSeparator = " "

type tag struct {
	Name                  string
	Default               string
	Required              bool
	DefaultValueSeparator string
	SkipOnNoValue         bool
}

func parseTag(s string) (tag, error) {
	tg := tag{
		DefaultValueSeparator: DefaultValueSeparator,
	}
	if s == "" {
		return tg, nil
	}

	parts := strings.Split(s, ",")
	tg.Name = parts[0]

	for _, part := range parts[1:] {
		switch part {
		case "required", "require", "req":
			tg.Required = true
		case "skip_on_no_value", "snv":
			tg.SkipOnNoValue = true
		default:
			vs := strings.SplitN(part, "=", 2)
			if len(vs) != 2 {
				continue
			}
			switch vs[0] {
			case "default":
				tg.Default = vs[1]
			case "separator", "sep":
				tg.DefaultValueSeparator = vs[1]
			default:
				panic(&UnknownTagPartError{
					Part:     vs[0],
					TagName:  tg.Name,
					TagValue: s,
				})
			}
		}
	}

	return tg, nil
}

func (tg tag) Encode() string {
	parts := []string{tg.Name}
	if tg.Required {
		parts = append(parts, "required")
	}
	if tg.Default != "" {
		parts = append(parts, fmt.Sprintf("default=%s", tg.Default))
	}
	if tg.DefaultValueSeparator != "" {
		parts = append(parts, fmt.Sprintf("separator=%s", tg.DefaultValueSeparator))
	}

	return strings.Join(parts, ",")
}
