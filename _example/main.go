package main

import (
	"fmt"

	"github.com/evenboee/env"
	_ "github.com/evenboee/env/autoload"
)

type Config struct {
	API struct {
		Port string `env:"PORT,default=8080"`
		Cors struct {
			// Can set key explicitly, but it will be inferred from the field name if not
			AllowedOrigins string `env:"ALLOWED_ORIGINS,default=*"`
			// Default values are separated by spaces
			AllowedMethods string `env:",default=GET POST PUT DELETE"`
			// Add custom default value separator
			AllowedHeaders string `env:",default=Content-Type|Authorization,sep=|"`
		}
	}

	Database struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	} `env:"DB"`
}

func main() {
	env.MustLoad("_example/.env")

	conf := env.MustBind[Config]()
	fmt.Printf("%+v\n", conf)
}
