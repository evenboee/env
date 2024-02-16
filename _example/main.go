package main

import (
	"fmt"

	"github.com/evenboee/env"
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
	env.MustLoad("_example/.env") // Load since path is not .env

	conf := env.MustBind[Config]()
	fmt.Printf("%+v\n", conf)

	port := env.GetString("API_PORT")
	fmt.Println("Port", port) // "" (empty)

	maxConnections := env.MustGet[int]("MAX_CONNECTIONS,default=100")
	fmt.Println("Max connections:", maxConnections) // 100

	mystrings := env.MustGet[[]myString]("MY_STRING,default=hello world|welcome home,sep=|")
	fmt.Println(mystrings) // ["hello world!" "welcome home!"]

	pointer1 := env.MustGet[*int]("NUM1")
	fmt.Println("Pointer1:", *pointer1) // 0

	pointer2 := env.MustGet[*int]("NUM2,skip_on_no_value")
	fmt.Println("Pointer2:", pointer2) // nil
}

type myString string

var _ env.StringUnmarshaler = (*myString)(nil)

func (ms *myString) UnmarshalString(s string) error {
	*ms = myString(s) + "!"
	return nil
}
