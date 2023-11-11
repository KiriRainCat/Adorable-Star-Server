package config

type Configuration struct {
	Server  Server  `mapstructure:"server" json:"server" yaml:"server"`
	Crawler Crawler `mapstructure:"crawler" json:"crawler" yaml:"crawler"`
}

type Server struct {
	Port        string `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	EncryptSalt string `mapstructure:"encrypt_salt" json:"encrypt_salt" yaml:"encrypt_salt"`
	RequestAuth string `mapstructure:"request_auth" json:"request_auth" yaml:"request_auth"`
	JwtEncrypt  string `mapstructure:"jwt_encrypt" json:"jwt_encrypt" yaml:"jwt_encrypt"`
	JwtIssuer   string `mapstructure:"jwt_issuer" json:"jwt_issuer" yaml:"jwt_issuer"`
	AdminAuth   string `mapstructure:"admin_auth" json:"admin_auth" yaml:"admin_auth"`
}

type Crawler struct {
	BrowserSocketUrl string `mapstructure:"browser_socket_url" json:"browser_socket_url" yaml:"browser_socket_url"`
	FetchInterval    int    `mapstructure:"fetch_interval" json:"fetch_interval" yaml:"fetch_interval"`
	MaxParallel      int    `mapstructure:"max_parallel" json:"max_parallel" yaml:"max_parallel"`
}
