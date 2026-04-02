package viperConfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Тест: Успешная загрузка кофига
func TestCheckSetConfig_Success(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "Form", "viperConfig")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatal("Не удалось создать директорию:", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `
app:
  name: "test-app"
  port: "8080"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal("Не удалось создать файл конфига:", err)
	}

	viper.AddConfigPath(configDir)

	err = CheckSetConfig()
	assert.NoError(t, err)

	assert.Equal(t, "test-app", viper.GetString("app.name"))
	assert.Equal(t, "8080", viper.GetString("app.port"))
	viper.Reset()
}
