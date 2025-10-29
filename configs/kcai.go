package configs

import (
	"github.com/songquanpeng/one-api/controller/oauth2"
	"github.com/spf13/viper"
)

type Kcai struct {
	Oauth2      Oauth2 `json:"oauth2"`
	CallbackUrl string `json:"callbackUrl"`
	Nextchat    string `json:"nextchat"`
}
type Oauth2 struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Domain       string `json:"domain"`
}

func (k *Kcai) GetNextchatUrl() string {
	return k.Nextchat
}

var (
	KcaiConfig Kcai
)

func OpenKcaiConfig(inputPath string) error {
	VIPER := viper.New()
	VIPER.AddConfigPath(inputPath)
	VIPER.SetConfigName("kcai")
	VIPER.SetConfigType("yaml")
	var dest Kcai
	if err := VIPER.ReadInConfig(); err != nil {
		return err
	}
	if err := VIPER.Unmarshal(&dest); err != nil {
		return err
	}
	KcaiConfig = dest
	return oauth2.SetConfig(dest.Oauth2.ClientId, dest.Oauth2.ClientSecret, dest.Oauth2.Domain, dest.CallbackUrl)
}
