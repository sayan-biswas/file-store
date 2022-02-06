package config

type Config struct {
	Server   server
	CORS     cors
	Database database
}

type server struct {
	Host        string
	Port        int
	Debug       bool
	TLS         bool
	Log         string
	Certificate string
	PrivateKey  string
}

type cors struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

type database struct {
	Diskless   bool
	Encryption bool
	CacheSize  int
	Path       string
}
