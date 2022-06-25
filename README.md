# gorm-filter

## Field Lookups

like [django filter](https://docs.djangoproject.com/en/4.0/ref/models/querysets/#field-lookups)

## Use

```shell
go get -u github.com/pandalzy/gorm-filter
```

```go
import (
    filter "github.com/pandalzy/gorm-filter"
)

type UserFilter struct {
    Username *string       `form:"username" filter:"field:username;expr:contains"`
    Name     string        `form:"name" filter:"field:name;expr:contains"`
    Age      []interface{} `form:"age" filter:"field:age;expr:in"`
    Email    string        `form:"email" filter:"field:email;expr:exact"`
}
```

## Example

See `example/main.go`

```shell
cd example
go run main.go
```

## output

```
[2.411ms] [rows:1] SELECT * FROM `users` WHERE username LIKE '%a%' AND name LIKE '%l%' AND age IN (20,19) AND `users`.`deleted_at` IS NULL
2022/06/25 09:24:56 [{{1 2022-06-20 21:29:19 +0800 CST 2022-06-20 21:29:21 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} panda lzy 20 }]
```