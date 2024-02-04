package config

type Configuration struct {
	Server     Server     `mapstructure:"server" json:"server" yaml:"server"`
	Crawler    Crawler    `mapstructure:"crawler" json:"crawler" yaml:"crawler"`
	Postgresql Postgresql `mapstructure:"postgresql" json:"postgresql" yaml:"postgresql"`
	Redis      Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
	SMTP       SMTP       `mapstructure:"smtp" json:"smtp" yaml:"smtp"`
	GPT        GPT        `mapstructure:"gpt" json:"gpt" yaml:"gpt"`
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
	ProxyPort     int `mapstructure:"proxy_port" json:"proxy_port,omitempty" yaml:"proxy_port"`
	FetchInterval int `mapstructure:"fetch_interval" json:"fetch_interval,omitempty" yaml:"fetch_interval"`
	MaxParallel   int `mapstructure:"max_parallel" json:"max_parallel,omitempty" yaml:"max_parallel"`
}

type Postgresql struct {
	DevHost  string `mapstructure:"dev_host" json:"dev_host,omitempty" yaml:"dev_host"`
	Host     string `mapstructure:"host" json:"host,omitempty" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	DB       string `mapstructure:"db" json:"db,omitempty" yaml:"db"`
	User     string `mapstructure:"user" json:"user,omitempty" yaml:"user"`
	Password string `mapstructure:"password" json:"password,omitempty" yaml:"password"`
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

type GPT struct {
	Enable   bool   `mapstructure:"enable" json:"enable,omitempty" yaml:"enable"`
	URL      string `mapstructure:"url" json:"url,omitempty" yaml:"url"`
	Username string `mapstructure:"username" json:"username,omitempty" yaml:"username"`
	Password string `mapstructure:"password" json:"password,omitempty" yaml:"password"`
}
