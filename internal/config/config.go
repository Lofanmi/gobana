package config

type Config struct {
	Application Application `json:"application"`
	QQWry       QQWry       `json:"qq_wry"`
	Providers   Providers   `json:"providers"`
	Backends    Backends    `json:"backends"`
	Indexes     Indexes     `json:"indexes"`
}

type Application struct {
	Production bool   `json:"production"`
	ListenAddr string `json:"listen_addr"`
}
