package viperConfig

import (
	"fmt"

	"github.com/spf13/viper"
)

func CheckSetConfig() {
	viper.SetConfigName("config")

	viper.AddConfigPath("./Form/viperConfig")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Конфиг не подключился: %s \n", err))
	}
}
