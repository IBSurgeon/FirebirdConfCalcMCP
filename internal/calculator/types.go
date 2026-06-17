package calculator

import (
	"fmt"
	"strings"
)

const DefaultBaseURL = "https://cc.ib-aid.com"

var (
	validArchitectures = map[string]struct{}{
		"Classic":      {},
		"SuperClassic": {},
		"SuperServer":  {},
	}
	validPageSizes = map[int]struct{}{
		4096:  {},
		8192:  {},
		16384: {},
		32768: {},
	}
)

type Request struct {
	MailLogin          string `json:"mailLogin"`
	PassAPI            string `json:"passApi"`
	ServerVersion      string `json:"serverVersion"`
	ServerArchitecture string `json:"serverArchitecture"`
	Cores              *int   `json:"cores,omitempty"`
	CountUsers         *int   `json:"countUsers,omitempty"`
	SizeDB             *int   `json:"sizeDb,omitempty"`
	PageSize           *int   `json:"pageSize,omitempty"`
	RAM                *int   `json:"ram,omitempty"`
	NameMainDB         string `json:"nameMainDb,omitempty"`
	PathToMainDB       string `json:"pathToMainDb,omitempty"`
	OSType             string `json:"osType"`
	HWType             string `json:"hwType"`
}

type CalculateParams struct {
	ServerVersion      string      `json:"server_version" jsonschema:"Firebird/HQbird version (fb2.5, fb3, fb4, fb5, hq2.5, hq5, etc.)"`
	ServerArchitecture string      `json:"server_architecture" jsonschema:"Classic, SuperClassic, or SuperServer"`
	Cores              OptionalInt `json:"cores,omitzero" jsonschema:"Number of CPU cores (1-100)"`
	CountUsers         OptionalInt `json:"count_users,omitzero" jsonschema:"Number of users (1-30000)"`
	SizeDB             OptionalInt `json:"size_db,omitzero" jsonschema:"Database size in GB"`
	PageSize           OptionalInt `json:"page_size,omitzero" jsonschema:"Page size: 4096, 8192, 16384, or 32768"`
	RAM                OptionalInt `json:"ram,omitzero" jsonschema:"Server RAM in GB (4-10000)"`
	NameMainDB         string      `json:"name_main_db,omitempty" jsonschema:"Main database name (max 100 chars)"`
	PathToMainDB       string      `json:"path_to_main_db,omitempty" jsonschema:"Path to main database (max 200 chars)"`
	OSType             string      `json:"os_type,omitempty" jsonschema:"Windows, Linux, or Universal (default Universal)"`
	HWType             string      `json:"hw_type,omitempty" jsonschema:"Hardware, Virtual, or Universal (default Universal)"`
}

type Response struct {
	InputParameters       string `json:"inputParameters"`
	ConfigurationFirebird string `json:"configurationFirebird"`
	ConfigurationDatabase string `json:"configurationDatabase"`
	MessageError          string `json:"messageError"`
}

type Result struct {
	InputParameters string `json:"input_parameters"`
	FirebirdConf    string `json:"firebird_conf"`
	DatabasesConf   string `json:"databases_conf"`
	APIVersion      string `json:"api_version"`
}

func (p CalculateParams) Validate() error {
	if strings.TrimSpace(p.ServerVersion) == "" {
		return fmt.Errorf("server_version is required")
	}
	if strings.TrimSpace(p.ServerArchitecture) == "" {
		return fmt.Errorf("server_architecture is required")
	}
	arch := normalizeArchitecture(p.ServerArchitecture)
	if _, ok := validArchitectures[arch]; !ok {
		return fmt.Errorf("server_architecture must be Classic, SuperClassic, or SuperServer")
	}
	if p.Cores != 0 && (p.Cores < 1 || p.Cores > 100) {
		return fmt.Errorf("cores must be between 1 and 100")
	}
	if p.CountUsers != 0 && (p.CountUsers < 1 || p.CountUsers > 30000) {
		return fmt.Errorf("count_users must be between 1 and 30000")
	}
	if p.RAM != 0 && (p.RAM < 4 || p.RAM > 10000) {
		return fmt.Errorf("ram must be between 4 and 10000")
	}
	if p.PageSize != 0 {
		if _, ok := validPageSizes[p.PageSize.Int()]; !ok {
			return fmt.Errorf("page_size must be 4096, 8192, 16384, or 32768")
		}
	}
	if len(p.NameMainDB) > 100 {
		return fmt.Errorf("name_main_db must not exceed 100 characters")
	}
	if len(p.PathToMainDB) > 200 {
		return fmt.Errorf("path_to_main_db must not exceed 200 characters")
	}
	return nil
}

func (p CalculateParams) ToRequest(credsMail, credsPass string) Request {
	osType := p.OSType
	if osType == "" {
		osType = "Universal"
	}
	hwType := p.HWType
	if hwType == "" {
		hwType = "Universal"
	}
	return Request{
		MailLogin:          credsMail,
		PassAPI:            credsPass,
		ServerVersion:      strings.TrimSpace(p.ServerVersion),
		ServerArchitecture: normalizeArchitecture(p.ServerArchitecture),
		Cores:              optionalIntPtr(p.Cores),
		CountUsers:         optionalIntPtr(p.CountUsers),
		SizeDB:             optionalIntPtr(p.SizeDB),
		PageSize:           optionalIntPtr(p.PageSize),
		RAM:                optionalIntPtr(p.RAM),
		NameMainDB:         p.NameMainDB,
		PathToMainDB:       p.PathToMainDB,
		OSType:             osType,
		HWType:             hwType,
	}
}

func normalizeArchitecture(arch string) string {
	switch strings.ToLower(strings.TrimSpace(arch)) {
	case "classic":
		return "Classic"
	case "superclassic":
		return "SuperClassic"
	case "superserver":
		return "SuperServer"
	default:
		return strings.TrimSpace(arch)
	}
}
