package calculator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/clc/calculation-params" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(Response{
			InputParameters:       `{"ram":16}`,
			ConfigurationFirebird: "ServerMode = Classic",
			ConfigurationDatabase: "{ DefaultDbCachePages = 250 }",
			MessageError:          "",
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	cores := 8
	users := 100
	ram := 16
	page := 4096
	params := CalculateParams{
		ServerVersion:      "fb3",
		ServerArchitecture: "Classic",
		Cores:              &cores,
		CountUsers:         &users,
		RAM:                &ram,
		PageSize:           &page,
	}

	result, err := client.Calculate(context.Background(), params.ToRequest("user@test.com", "pass"))
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	if result.FirebirdConf == "" {
		t.Fatal("expected firebird conf")
	}
}

func TestCalculateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Response{
			MessageError: "Count of users is null. Please, set value.",
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Calculate(context.Background(), Request{
		MailLogin:          "u",
		PassAPI:            "p",
		ServerVersion:      "fb3",
		ServerArchitecture: "Classic",
		OSType:             "Universal",
		HWType:             "Universal",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateParams(t *testing.T) {
	bad := CalculateParams{ServerVersion: "", ServerArchitecture: "Classic"}
	if err := bad.Validate(); err == nil {
		t.Fatal("expected validation error")
	}

	good := CalculateParams{ServerVersion: "fb3", ServerArchitecture: "classic"}
	if err := good.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
