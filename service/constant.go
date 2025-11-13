package service

// ---------------------------------------------------------------------------------------------------------------------

type ProviderName = string

type ProviderType = string

const (
	ProviderTypeSls           ProviderType = "sls"
	ProviderTypeElasticsearch ProviderType = "elasticsearch"
	ProviderTypeKibanaProxy   ProviderType = "kibana_proxy"
)

// ---------------------------------------------------------------------------------------------------------------------

const (
	AtTimestamp    = "@timestamp"
	DefaultValue   = "default_value"
	MaxChartPoints = 60
)

// ---------------------------------------------------------------------------------------------------------------------

type ParserName = string

type (
	ParserFieldType   = string
	ParserFieldReturn = string
)

// ---------------------------------------------------------------------------------------------------------------------

const (
	ParserFieldTypeReplacements ParserFieldType = "replacements"
	ParserFieldTypeJavaScript   ParserFieldType = "javascript"

	ParserFieldReturnString ParserFieldReturn = "string"
	ParserFieldReturnNumber ParserFieldReturn = "number"
)

// ---------------------------------------------------------------------------------------------------------------------
