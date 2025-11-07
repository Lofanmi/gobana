package config

// ---------------------------------------------------------------------------------------------------------------------

type QQWry struct {
	IPv4FilePath string `json:"ipv4_file_path"`
	IPv4Data     []byte `json:"-"`
	IPv6FilePath string `json:"ipv6_file_path"`
	IPv6Data     []byte `json:"-"`
}

// ---------------------------------------------------------------------------------------------------------------------
