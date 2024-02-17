package main

import (
	"fmt"
	"time"

	"github.com/evenboee/env"
)

type Config struct {
	API struct {
		Port string `env:"PORT,default=8080,osp"`
		Cors struct {
			// Can set key explicitly, but it will be inferred from the field name if not
			AllowedOrigins []string `env:"ALLOWED_ORIGINS,default=*"`
			// Default values are separated by spaces
			AllowedMethods []string `env:",default=GET POST PUT DELETE"`
			// Add custom default value separator
			AllowedHeaders []string `env:",default=Content-Type|Authorization,dsep=|"`
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
	conf := env.MustBind[Config]()
	fmt.Printf("%+v\n", conf)

	port := env.GetString("API_PORT")
	fmt.Println("Port: ", port) // "" (empty)

	maxConnections := env.MustGet[int]("MAX_CONNECTIONS,default=100")
	fmt.Println("Max connections:", maxConnections) // 100

	mystrings := env.MustGet[[]myString]("MY_STRING,default=hello world|welcome home,dsep=|")
	fmt.Println("[]MyString: ", mystrings) // ["hello world!" "welcome home!"]

	pointer1 := env.MustGet[*int]("NUM1")
	fmt.Println("Pointer1:", *pointer1) // 0

	pointer2 := env.MustGet[*int]("NUM2,skip_on_no_value")
	fmt.Println("Pointer2:", pointer2) // nil

	t := env.MustGet[time.Time]("DEADLINE,default=23:59:59 2020-12-31", "time_format=15:04:05 2006-01-02")
	fmt.Println("Deadline:", t) // 2020-12-31 23:59:59 +0000 UTC
}

type myString string

var _ env.StringUnmarshaler = (*myString)(nil)

func (ms *myString) UnmarshalString(s string) error {
	*ms = myString(s) + "!"
	return nil
}
