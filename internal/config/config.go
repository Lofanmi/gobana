package config

type Config struct {
	Application Application `json:"application"`
	QQWry       QQWry       `json:"-"`

	Providers Providers `json:"providers"`
	Backends  Backends  `json:"backends"`
	Indexes   Indexes   `json:"indexes"`
}

type Application struct {
	Timezone   string `json:"timezone"`
	Production bool   `json:"production"`
	ListenAddr string `json:"listen_addr"`
}

type QQWry struct {
	IPv4FilePath string `json:"ipv4_file_path"`
	IPv4Data     []byte `json:"-"`
	IPv6FilePath string `json:"ipv6_file_path"`
	IPv6Data     []byte `json:"-"`
}
