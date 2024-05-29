package configs

type Mode string

const (
	ModeS3 Mode = "s3"
)

type S3 struct {
	Mode     Mode   `yaml:"mode"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
