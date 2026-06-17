package calculator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// OptionalInt is an optional integer tool parameter. It marshals as a plain
// integer in JSON Schema (not null|integer union) and accepts string values
// from clients that send numeric fields as strings.
type OptionalInt int

func (o *OptionalInt) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "" || s == "null" {
		*o = 0
		return nil
	}
	if strings.HasPrefix(s, `"`) {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		str = strings.TrimSpace(str)
		if str == "" {
			*o = 0
			return nil
		}
		v, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid integer %q: %w", str, err)
		}
		*o = OptionalInt(v)
		return nil
	}
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid integer: %w", err)
	}
	*o = OptionalInt(v)
	return nil
}

func (o OptionalInt) Int() int {
	return int(o)
}

func optionalIntPtr(v OptionalInt) *int {
	if v == 0 {
		return nil
	}
	i := int(v)
	return &i
}
