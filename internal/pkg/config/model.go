package config

type Configuration struct {
	Server  Server  `mapstructure:"server" json:"server" yaml:"server"`
	Crawler Crawler `mapstructure:"crawler" json:"crawler" yaml:"crawler"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	SMTP    SMTP    `mapstructure:"smtp" json:"smtp" yaml:"smtp"`
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
	BrowserSocketUrl      string `mapstructure:"browser_socket_url" json:"browser_socket_url,omitempty" yaml:"browser_socket_url"`
	ProxyBrowserSocketUrl string `mapstructure:"proxy_browser_socket_url" json:"proxy_browser_socket_url,omitempty" yaml:"proxy_browser_socket_url"`
	FetchInterval         int    `mapstructure:"fetch_interval" json:"fetch_interval,omitempty" yaml:"fetch_interval"`
	MaxParallel           int    `mapstructure:"max_parallel" json:"max_parallel,omitempty" yaml:"max_parallel"`
}

type Redis struct {
	Host     string `mapstructure:"host" json:"host,omitempty" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	Password string `mapstructure:"password" json:"password,omitempty" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db,omitempty" yaml:"db"`
}

type SMTP struct {
	Host string `mapstructure:"host" json:"host,omitempty" yaml:"host"`
	Port int    `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	Key  string `mapstructure:"key" json:"key,omitempty" yaml:"key"`
	Mail string `mapstructure:"mail" json:"mail,omitempty" yaml:"mail"`
}
