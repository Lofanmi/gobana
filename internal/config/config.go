package config

type Config struct {
	BackendList BackendSlice `json:"backend_list"`
}

type BackendSlice []Backend

func (s BackendSlice) Match(name, typ string) *Backend {
	if len(s) <= 0 {
		return nil
	}
	for _, backend := range s {
		if backend.Name == name && backend.Type == typ {
			return &backend
		}
	}
	return nil
}

type Backend struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Addr        string                 `json:"addr"`
	Auth        Auth                   `json:"auth"`
	MultiSearch map[string]MultiSearch `json:"multi_search"`
}

type Auth struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type MultiSearch struct {
	Project string   `json:"project"`
	Storage []string `json:"storage"`
}
