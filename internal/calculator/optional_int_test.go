package calculator

import (
	"encoding/json"
	"testing"
)

func TestOptionalIntUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    OptionalInt
		wantErr bool
	}{
		{name: "integer", input: `8`, want: 8},
		{name: "string integer", input: `"16"`, want: 16},
		{name: "null", input: `null`, want: 0},
		{name: "empty string", input: `""`, want: 0},
		{name: "invalid string", input: `"abc"`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got OptionalInt
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalculateParamsUnmarshalStringNumbers(t *testing.T) {
	raw := `{
		"server_version": "fb3",
		"server_architecture": "Classic",
		"cores": "8",
		"ram": "16",
		"count_users": "100",
		"size_db": "100",
		"page_size": "4096"
	}`

	var params CalculateParams
	if err := json.Unmarshal([]byte(raw), &params); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if params.Cores != 8 || params.RAM != 16 || params.CountUsers != 100 || params.SizeDB != 100 || params.PageSize != 4096 {
		t.Fatalf("unexpected params: %+v", params)
	}
	if err := params.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
