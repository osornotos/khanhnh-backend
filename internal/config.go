package internal

type Config struct {
	Version  string      `yaml:"version"`
	Port     int         `yaml:"port"`
	Env      string      `yaml:"env"`
	Debug    bool        `yaml:"debug"`
	Hash     HashConfig  `json:"hash"`
	Postgres PgConfig    `yaml:"postgres"`
	Auth     AuthConfig  `yaml:"auth"`
	Redis    RedisConfig `yaml:"redis"`
}

func (c Config) Safe() Config {
	c.Postgres.User = "***"
	c.Postgres.Password = "***"
	return c
}

func (c Config) Validate() error {
	// Config validation here
	return nil
}

type HashConfig struct {
	// time represents the number of
	// passed over the specified memory.
	Time uint32 `yaml:"time"`
	// cpu memory to be used. MB
	Memory uint32 `yaml:"memory"`
	// threads for parallelism aspect
	// of the algorithm.
	Threads uint8 `yaml:"threads"`
	// keyLen of the generate hash key.
	KeyLen uint32 `yaml:"keyLen"`
	// saltLen the length of the salt used.
	SaltLen uint32 `yaml:"saltLen"`
}

type PgConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	PoolSize int    `yaml:"poolSize"`
}

type AuthConfig struct {
	JwtPublicKey  string `yaml:"jwtPublicKey"`
	JwtPrivateKey string `yaml:"JwtPrivateKey"`
}

type RedisConfig struct {
	UseSentinel      bool   `yaml:"useSentinel"`
	Host             string `yaml:"host"`
	MasterName       string `yaml:"masterName"`
	SentinelPassword string `yaml:"sentinelPassword"`
	Password         string `yaml:"password"`
	Database         int    `yaml:"database"`
}
