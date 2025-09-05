package config

import (
	"time"

	"github.com/higress-group/wasm-go/pkg/wrapper"
	"github.com/tidwall/gjson"
)

func init() {
	RegisterCtxOption(wrapper.ParseConfig(parseFederationConfig))
}

type FederationConfig struct {
	EnableQueryPlan  bool          `json:"enableQueryPlanning"`
	EnableCaching    bool          `json:"enableCaching"`
	MaxQueryDepth    int64         `json:"maxQueryDepth"`
	QueryTimeout     time.Duration `json:"queryTimeout"`
	EnableIntrospect bool          `json:"enableIntrospection"`
	GraphqlAddress   string        `json:"graphqlAddress"`
	DebugMode        bool          `json:"debugMode"`
}

func parseFederationConfig(json gjson.Result, config *FederationConfig) error {
	config.EnableQueryPlan = json.Get("enableQueryPlanning").Bool()
	config.EnableCaching = json.Get("enableCaching").Bool()
	config.MaxQueryDepth = json.Get("maxQueryDepth").Int()
	config.QueryTimeout = time.Duration(json.Get("queryTimeout").Int()) * time.Millisecond
	config.EnableIntrospect = json.Get("enableIntrospection").Bool()
	config.GraphqlAddress = json.Get("graphqlAddress").String()
	config.DebugMode = json.Get("debugMode").Bool()
	return nil
}
