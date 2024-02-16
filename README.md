# env

Wrapper for [godotenv](https://github.com/joho/godotenv) with extra functionality to bind to values. 

Binding functions are heavily inspired by the [gin](https://github.com/gin-gonic/gin) binding package. 

Supported types:
- int/uint (all types)
- bool (strconv.Parsebool + yes/no)
- float32/64
- string
- pointers
- arrays
- slices
- structs
- time.Time and time.Duration

## set parameters

there are options modifying the behaviour of the set function. 

`Prefix` is used to specify prefix to be prepended to all values. default is "". 

`Sep` is used when combining keys. e.g. prefix is "API" and key is "PORT" and separator is "_" the resulting key is "API_PORT". default is "_". 

`Tag` is used to specify the tag to use when getting key. default is "env". 

`ArrSep` is used to specify the separator of array and slice values from env. e.g. sep="," - "1,2,3" => [1, 2, 3]. default is ",". 

`AutoFormatMissingKeys` is used to signal is you want missing keys to be automatically formated. See [formating](#auto-key-formating) for more info. 


## struct tag

the main struct tag has options to modify the behaviour when setting the value. 
tag parts are separated by a comma. 
key-value parts have the format `key=value`, boolean parts only have a key. 

here are the available parts

### Name

it is the first value and does not have a key. this, together with the prefix, forms the key to get env variable. 

### Default

the value to set is there is no env variable. this value requires the key `default` with a value. 

### Required

if there is no env variable with key and this is present, the set function will return an error saying there is a missing key. this is a boolean part where the key can be one of: `required`, `require`, or `req`. 

### DefaultValueSeparator

in the case where the default value is used the value of the default will be split by a value. by default this value will be env.DefaultValueSeparator (default: " ") but this can be specified on a per field basis with the DefaultValueSeparator part. this is a key-value part with the key being either `separator` or `sep`. 


### SkipOnNoValue

This is primarily used for pointer types. 
When the type is a pointer, the set function will try to set the value if it is null and then try to set the underlying value. this is good for structs or other types where you want to set the underlying value, but not for primitive types such as *int where you want to have nil when there is no value and not 0. 
Use keys `skip_on_no_value` or `snv` to skip setting values when no env variable was found for the key. 


## structs

use env.Bind to set structs. 
when settings structs the substructs will get the parents key appended as a prefix. a substruct named DB will give all values in the struct a `DB` prefix. this will add up if there will be many levels of prefixes. 

See `_example/main.go` for use. 

## time.Time

time.Time is a struct but is treated differently. 
For time.Time fields there are a few more tags that can be set. this is taken from the gin binding function. 

`time_format` can be used to specify the time format. RFC3339 will be used if nothing is specified. 

`time_utc` can be used to specify time to be utc if not `time_location` can be used to specify location acording to time.LoadLocation


## single values

the Get function is used to first generate the values to emulate a struct field by generating tags before binding. 
first parameter is the normal format for the `env` struct tag. 
other tags can be specified in format `key=value`


## Unmarshal interface

before setting value the set function will take the address of the value and check if that implements either the StringUnmarshaler or TextUnmarshaler interfaces and will use that if it exists. 


## Auto Key formating

examples of formating:

- "DBName" -> "DB_NAME"
- "DB" -> "DB"
- "Name" -> "NAME"
- "AllowedOrigins" -> "ALLOWED_ORIGINS"
- "TestA" -> "TEST_A"
