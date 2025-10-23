module github.com/lucasvillarinho/restql/examples/echo-gorm

go 1.25.1

require (
	github.com/labstack/echo/v4 v4.12.0
	github.com/lucasvillarinho/restql v0.0.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.12
)

require (
	github.com/alecthomas/participle/v2 v2.1.4 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)

replace github.com/lucasvillarinho/restql => ../..
