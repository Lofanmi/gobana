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
	GobanaTimestampFunc  = "gobanaTimestamp"
	GobanaTimestampField = "gobanaTimestamp"
	MaxChartPoints       = 60
)

// ---------------------------------------------------------------------------------------------------------------------

type ParserName = string

type (
	ParserFieldType   = string
	ParserFieldReturn = string
)

// ---------------------------------------------------------------------------------------------------------------------

const (
	ParserFieldTypeTimestamp    ParserFieldType = "timestamp"
	ParserFieldTypeReplacements ParserFieldType = "replacements"
	ParserFieldTypeJavaScript   ParserFieldType = "javascript"

	ParserFieldReturnString ParserFieldReturn = "string"
	ParserFieldReturnNumber ParserFieldReturn = "number"
)

// ---------------------------------------------------------------------------------------------------------------------
