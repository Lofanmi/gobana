package constant

const (
	AtTimestamp    = "@timestamp"
	DefaultValue   = "default_value"
	MaxChartPoints = 60
)

type ClientType = string

const (
	ClientTypeElasticsearch ClientType = "elasticsearch"
	ClientTypeKibanaProxy   ClientType = "kibana-proxy"
	ClientTypeSLS           ClientType = "sls"
)

type (
	ParserFieldType   = string
	ParserFieldReturn = string
)

const (
	ParserFieldTypeReplacements ParserFieldType = "replacements"
	ParserFieldTypeLua          ParserFieldType = "lua"

	ParserFieldReturnString ParserFieldReturn = "string"
	ParserFieldReturnNumber ParserFieldReturn = "number"
)
