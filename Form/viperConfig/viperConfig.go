package viperConfig

import (
	"fmt"

	"github.com/spf13/viper"
)

func CheckSetConfig() error {
	viper.SetConfigName("config")

	viper.AddConfigPath("./Form/viperConfig")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Конфиг не подключился: %s \n", err)
	}
	return nil
}
