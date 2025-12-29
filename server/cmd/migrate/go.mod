module github.com/kk/kkartifact-server/cmd/migrate

go 1.21

replace github.com/kk/kkartifact-server => ../..

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/kk/kkartifact-server v0.0.0-00010101000000-000000000000
)

