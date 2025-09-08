package configuration

type Config struct {
	CurrentBuild string            `mapstructure:"current_build"`
	Builds       map[string]string `mapstructure:"builds"`
}
