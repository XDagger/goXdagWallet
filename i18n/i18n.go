package i18n

import (
	"encoding/json"
	"goXdagWallet/config"
	"io/ioutil"
	"path"
)

var i18n = make(map[string]string)

func stringRead(res string) []byte {
	bytes, err := ioutil.ReadFile(path.Join("data", res))
	if err != nil {
		return nil
	}
	return bytes
}

func LoadI18nStrings() error {
	lang := config.GetConfig().CultureInfo
	data := stringRead(lang + ".json")
	if len(data) == 0 {
		return nil
	}
	err := json.Unmarshal(data, &i18n)
	if err != nil {
		return err
	}
	return nil
}

func GetString(id string) string {
	str := i18n[id]
	if str != "" {
		return str
	}
	return id
}
