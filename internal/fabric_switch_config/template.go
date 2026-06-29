package fabric_switch_config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\{([^{}]+)\}`)

// ordinalByFunc returns the (1-based) ordinal for a tag key, used by the
// {ordinalBy:<key>} / {ordinalBy0:<key>} placeholders. nil if unsupported.
type ordinalByFunc func(key string) (int, error)

// expandTemplate expands placeholders in template. No eval; unknown placeholders
// return a ConfigError. Supported forms:
//   - {tag:<key>[:fmt]}        value from tags (error if missing)
//   - {ordinalBy:<key>[:fmt]}  1-based ordinal among same-position devices
//     sharing the device's value for tag <key>
//   - {ordinalBy0:<key>[:fmt]} same but 0-based
//   - {name}/{name:fmt}        from values, with optional format spec
func expandTemplate(template string, tags map[string]string, ordinalBy ordinalByFunc, values map[string]any) (string, error) {
	var firstErr error
	out := placeholderRe.ReplaceAllStringFunc(template, func(match string) string {
		inner := match[1 : len(match)-1]
		s, err := expandOne(inner, template, tags, ordinalBy, values)
		if err != nil && firstErr == nil {
			firstErr = err
		}
		return s
	})
	if firstErr != nil {
		return "", firstErr
	}
	return out, nil
}

func expandOne(inner, template string, tags map[string]string, ordinalBy ordinalByFunc, values map[string]any) (string, error) {
	// Longest prefix first (ordinalBy0 before ordinalBy).
	for _, prefix := range []string{"tag", "ordinalBy0", "ordinalBy"} {
		if !strings.HasPrefix(inner, prefix+":") {
			continue
		}
		rest := inner[len(prefix)+1:]
		key, fmtSpec := rpartition(rest, ":")

		if prefix == "tag" {
			value, ok := tags[key]
			if !ok {
				return "", configErrorf("tag '%s' not present in tagsMap", key)
			}
			return pyFormat(value, fmtSpec)
		}
		// ordinalBy / ordinalBy0
		if ordinalBy == nil {
			return "", configErrorf("'{%s}' is not supported in template %q", inner, template)
		}
		ordinal, err := ordinalBy(key)
		if err != nil {
			return "", err
		}
		if prefix == "ordinalBy0" {
			ordinal--
		}
		return pyFormat(int64(ordinal), fmtSpec)
	}

	// Plain {name} / {name:fmt}.
	name, fmtSpec := partition(inner, ":")
	value, ok := values[name]
	if !ok {
		return "", configErrorf("unknown placeholder '{%s}' in template %q", inner, template)
	}
	if fmtSpec == "" && !strings.Contains(inner, ":") {
		return pyStr(value), nil
	}
	return pyFormat(value, fmtSpec)
}

// rpartition splits on the LAST sep; if absent, the whole string is the key and
// fmt is empty (mirrors Python str.rpartition then the "if not sep" fallback).
func rpartition(s, sep string) (key, rest string) {
	idx := strings.LastIndex(s, sep)
	if idx < 0 {
		return s, ""
	}
	return s[:idx], s[idx+len(sep):]
}

// partition splits on the FIRST sep.
func partition(s, sep string) (head, tail string) {
	idx := strings.Index(s, sep)
	if idx < 0 {
		return s, ""
	}
	return s[:idx], s[idx+len(sep):]
}

func pyStr(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// pyFormat mirrors Python's format(value, fmt). An empty fmt is str(value). A
// fmt whose presentation type is numeric (b/c/d/o/x/X/n) coerces the value to
// int first (so string tag values can be zero-padded), then formats as an
// integer; otherwise the value is formatted as a string.
func pyFormat(value any, fmtSpec string) (string, error) {
	if fmtSpec == "" {
		return pyStr(value), nil
	}
	last := fmtSpec[len(fmtSpec)-1]
	if strings.IndexByte("bcdoxXn", last) >= 0 {
		n, ok := toInt(value)
		if !ok {
			return "", configErrorf("value %v for format %q is not numeric", value, fmtSpec)
		}
		return formatInt(n, fmtSpec)
	}
	return formatStr(pyStr(value), fmtSpec)
}

func toInt(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

// formatInt implements the subset of Python's integer format mini-language used
// by the reference templates: [[fill]align][0][width][type] with type in
// d/x/X/o/b/n (n treated as d). Sign/grouping/# flags are not used and ignored.
func formatInt(n int64, spec string) (string, error) {
	fill, align, rest := parseFillAlign(spec)
	if strings.HasPrefix(rest, "0") {
		// A leading '0' implies zero-padding with '=' alignment.
		if align == 0 {
			fill, align = '0', '='
		}
		rest = rest[1:]
	}
	width, typ, err := parseWidthType(rest, "bcdoxXn")
	if err != nil {
		return "", err
	}

	neg := n < 0
	mag := n
	if neg {
		mag = -n
	}
	var body string
	switch typ {
	case 'x':
		body = strconv.FormatInt(mag, 16)
	case 'X':
		body = strings.ToUpper(strconv.FormatInt(mag, 16))
	case 'o':
		body = strconv.FormatInt(mag, 8)
	case 'b':
		body = strconv.FormatInt(mag, 2)
	default: // 'd', 'n', 'c'(unused) -> decimal
		body = strconv.FormatInt(mag, 10)
	}
	sign := ""
	if neg {
		sign = "-"
	}

	if align == 0 {
		align = '>' // numbers default to right-aligned
	}
	return pad(sign, body, fill, align, width), nil
}

// formatStr implements the subset of Python's string format mini-language used:
// [[fill]align][width]. Strings default to left-aligned.
func formatStr(s, spec string) (string, error) {
	fill, align, rest := parseFillAlign(spec)
	width, _, err := parseWidthType(rest, "s")
	if err != nil {
		return "", err
	}
	if align == 0 {
		align = '<'
	}
	if fill == 0 {
		fill = ' '
	}
	return pad("", s, fill, align, width), nil
}

func parseFillAlign(spec string) (fill byte, align byte, rest string) {
	if len(spec) >= 2 && isAlign(spec[1]) {
		return spec[0], spec[1], spec[2:]
	}
	if len(spec) >= 1 && isAlign(spec[0]) {
		return ' ', spec[0], spec[1:]
	}
	return 0, 0, spec
}

func isAlign(b byte) bool { return b == '<' || b == '>' || b == '^' || b == '=' }

func parseWidthType(rest, types string) (width int, typ byte, err error) {
	if rest == "" {
		return 0, 0, nil
	}
	last := rest[len(rest)-1]
	if last >= '0' && last <= '9' {
		// No type char; the whole rest is width.
		w, e := strconv.Atoi(rest)
		if e != nil {
			return 0, 0, configErrorf("invalid format width %q", rest)
		}
		return w, 0, nil
	}
	if strings.IndexByte(types, last) < 0 {
		return 0, 0, configErrorf("unsupported format type %q", string(last))
	}
	digits := rest[:len(rest)-1]
	if digits == "" {
		return 0, last, nil
	}
	w, e := strconv.Atoi(digits)
	if e != nil {
		return 0, 0, configErrorf("invalid format width %q", digits)
	}
	return w, last, nil
}

func pad(sign, body string, fill, align byte, width int) string {
	if fill == 0 {
		fill = ' '
	}
	total := len(sign) + len(body)
	if total >= width {
		return sign + body
	}
	padLen := width - total
	f := string(fill)
	switch align {
	case '<':
		return sign + body + strings.Repeat(f, padLen)
	case '^':
		left := padLen / 2
		right := padLen - left
		return strings.Repeat(f, left) + sign + body + strings.Repeat(f, right)
	case '=':
		return sign + strings.Repeat(f, padLen) + body
	default: // '>'
		return strings.Repeat(f, padLen) + sign + body
	}
}
