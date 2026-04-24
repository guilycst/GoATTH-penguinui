package schemafield

import (
	"fmt"
	"sort"
	"strings"
)

// Kind is the rendered input kind.
type Kind string

const (
	KindString  Kind = "string"
	KindNumber  Kind = "number"
	KindInteger Kind = "integer"
	KindBoolean Kind = "boolean"
	KindEnum    Kind = "enum"
	KindObject  Kind = "object"   // group of nested fields
	KindArray   Kind = "array"    // simple array of scalars (renders as tagslist)
	KindUnknown Kind = "unknown"  // fallback — render read-only JSON
)

// Field describes one renderable input derived from a JSON Schema node + defaults.
type Field struct {
	// Path is a dotted JSONPath from the root of the values object (e.g. "auth.password").
	Path string
	// Name is the input name attribute; usually Path.
	Name string
	// Label is the human-readable field label.
	Label string
	// Description is an optional helper text.
	Description string
	// Kind determines rendering.
	Kind Kind
	// Required marks the field with an asterisk.
	Required bool
	// Managed means AllowList matches this path — rendered read-only with a lock.
	Managed bool
	// Default is the JSON-serialized default value (quoted scalars preserved).
	Default string
	// Value is the current form value (JSON-serialized).
	Value string
	// Enum lists valid options when Kind == KindEnum.
	Enum []string
	// Children are nested fields when Kind == KindObject.
	Children []Field
	// Errors hold per-field validation errors.
	Errors []string
}

// Walk builds a flat-ish list of Fields from a JSON Schema subset.
// - schema: the JSON Schema root (typically an object with "properties").
// - defaults: default values map (from AddonVersion.DefaultValues).
// - values: current form values (user-submitted so far).
// - allowList: set of paths the platform forbids overriding — rendered read-only.
//
// Returns the list of top-level fields. Nested objects produce fields with Children.
// Unknown / unsupported shapes produce a KindUnknown field — the caller can decide
// to fall back to a raw YAML editor for the whole form.
func Walk(schema map[string]any, defaults map[string]any, values map[string]any, allowList map[string]bool) []Field {
	if schema == nil {
		return nil
	}
	props, ok := schema["properties"].(map[string]any)
	if !ok || len(props) == 0 {
		return nil
	}
	required := stringSet(schema["required"])

	// Deterministic ordering: schema-declared order isn't preserved by map, so sort.
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]Field, 0, len(keys))
	for _, k := range keys {
		node, _ := props[k].(map[string]any)
		if node == nil {
			continue
		}
		out = append(out, buildField(k, k, node, defaults, values, allowList, required[k]))
	}
	return out
}

func buildField(path, name string, node map[string]any, defaults, values map[string]any, allowList map[string]bool, required bool) Field {
	f := Field{
		Path:     path,
		Name:     name,
		Label:    humanize(name),
		Required: required,
		Managed:  allowList[path],
	}
	if t, ok := getString(node, "title"); ok && t != "" {
		f.Label = t
	}
	if d, ok := getString(node, "description"); ok {
		f.Description = d
	}

	// Enum takes precedence over raw type.
	if enumRaw, ok := node["enum"].([]any); ok && len(enumRaw) > 0 {
		f.Kind = KindEnum
		f.Enum = toStringSlice(enumRaw)
	} else {
		tp, _ := getString(node, "type")
		f.Kind = kindFromJSONType(tp)
	}

	f.Default = jsonScalar(lookupDotted(defaults, path))
	f.Value = jsonScalar(lookupDotted(values, path))

	if f.Kind == KindObject {
		childProps, _ := node["properties"].(map[string]any)
		if childProps != nil {
			required := stringSet(node["required"])
			keys := make([]string, 0, len(childProps))
			for k := range childProps {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				childNode, _ := childProps[k].(map[string]any)
				if childNode == nil {
					continue
				}
				childPath := path + "." + k
				f.Children = append(f.Children, buildField(childPath, childPath, childNode, defaults, values, allowList, required[k]))
			}
		}
	}

	if f.Kind == KindArray {
		// Only support simple arrays of scalars in MVP. Anything more complex downgrades.
		if items, ok := node["items"].(map[string]any); ok {
			if tp, _ := getString(items, "type"); tp == "" || tp == "object" || tp == "array" {
				f.Kind = KindUnknown
			}
		}
	}

	return f
}

// FallbackFromDefaults produces fields when no schema is available.
// Every top-level default becomes a KindString field (the user can still edit nested
// values via the YAML escape hatch). Returns nil if no defaults.
func FallbackFromDefaults(defaults map[string]any, values map[string]any, allowList map[string]bool) []Field {
	if len(defaults) == 0 {
		return nil
	}
	keys := make([]string, 0, len(defaults))
	for k := range defaults {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]Field, 0, len(keys))
	for _, k := range keys {
		val := defaults[k]
		// Nested objects/arrays downgrade to KindUnknown so the caller surfaces a
		// YAML editor rather than rendering a misleading input.
		kind := kindFromGo(val)
		f := Field{
			Path:    k,
			Name:    k,
			Label:   humanize(k),
			Kind:    kind,
			Managed: allowList[k],
			Default: jsonScalar(val),
			Value:   jsonScalar(lookupDotted(values, k)),
		}
		out = append(out, f)
	}
	return out
}

// HasOnlySimpleScalars reports whether every field in the list is a simple scalar —
// no objects, arrays, or unknowns. Used to decide whether the caller can render a
// pure form or needs the YAML escape hatch.
func HasOnlySimpleScalars(fields []Field) bool {
	for _, f := range fields {
		switch f.Kind {
		case KindObject, KindUnknown:
			return false
		case KindArray:
			if len(f.Children) > 0 {
				return false
			}
		}
	}
	return true
}

// --- helpers ---

func kindFromJSONType(tp string) Kind {
	switch tp {
	case "string":
		return KindString
	case "number":
		return KindNumber
	case "integer":
		return KindInteger
	case "boolean":
		return KindBoolean
	case "object":
		return KindObject
	case "array":
		return KindArray
	default:
		return KindUnknown
	}
}

func kindFromGo(v any) Kind {
	switch v.(type) {
	case string:
		return KindString
	case bool:
		return KindBoolean
	case int, int32, int64, uint, uint32, uint64:
		return KindInteger
	case float32, float64:
		return KindNumber
	case map[string]any:
		return KindObject
	case []any:
		return KindArray
	default:
		return KindUnknown
	}
}

func getString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func stringSet(v any) map[string]bool {
	out := map[string]bool{}
	arr, ok := v.([]any)
	if !ok {
		return out
	}
	for _, a := range arr {
		if s, ok := a.(string); ok {
			out[s] = true
		}
	}
	return out
}

func toStringSlice(raw []any) []string {
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		out = append(out, fmt.Sprint(v))
	}
	return out
}

func jsonScalar(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return fmt.Sprintf("%v", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}

// lookupDotted walks a dotted path through a nested map. Missing segments return nil.
func lookupDotted(m map[string]any, path string) any {
	if m == nil || path == "" {
		return nil
	}
	segs := strings.Split(path, ".")
	var cur any = m
	for _, s := range segs {
		mp, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = mp[s]
	}
	return cur
}

// humanize converts a snake_or-camel path segment into a Title Case label.
func humanize(s string) string {
	// Take only the last segment for the label.
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	// Split camelCase boundaries.
	var b strings.Builder
	for i, r := range s {
		if i > 0 && isUpper(r) && !isUpper(prevRune(s, i)) {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
	}
	out := b.String()
	if out == "" {
		return s
	}
	// Capitalize first letter.
	return strings.ToUpper(out[:1]) + out[1:]
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }

func prevRune(s string, i int) rune {
	if i == 0 {
		return 0
	}
	for j := i - 1; j >= 0; j-- {
		r := rune(s[j])
		if r != 0 {
			return r
		}
	}
	return 0
}
