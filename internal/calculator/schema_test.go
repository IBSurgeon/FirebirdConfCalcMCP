package calculator

import (
	"encoding/json"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
)

func TestCalculateParamsSchemaHasIntegerTypes(t *testing.T) {
	schema, err := jsonschema.For[CalculateParams](nil)
	if err != nil {
		t.Fatalf("For[CalculateParams]() error = %v", err)
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Marshal schema: %v", err)
	}
	t.Logf("schema: %s", data)

	intFields := []string{"cores", "count_users", "size_db", "page_size", "ram"}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal schema: %v", err)
	}
	props, _ := parsed["properties"].(map[string]any)
	if props == nil {
		t.Fatal("missing properties")
	}
	for _, field := range intFields {
		prop, ok := props[field].(map[string]any)
		if !ok {
			t.Fatalf("field %q missing from properties", field)
		}
		if !schemaPropertyIsInteger(prop) {
			t.Fatalf("field %q schema = %#v, want integer type", field, prop)
		}
	}
}

func schemaPropertyIsInteger(prop map[string]any) bool {
	if typ, ok := prop["type"].(string); ok && typ == "integer" {
		return true
	}
	return false
}
