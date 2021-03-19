package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"github.com/fransoaardi/url-shortener/internal/server"
)

func init() {
	viper.SetConfigFile("./internal/config/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	s := http.Server{
		Addr:    ":9000",
		Handler: server.New(),
	}

	fmt.Println("serve on :9000")
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
