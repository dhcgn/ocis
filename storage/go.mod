module github.com/owncloud/ocis/storage

go 1.15

require (
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/cs3org/reva v1.6.1-0.20210218130854-a954e7242e7a
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/mitchellh/copystructure v1.1.1 // indirect
	github.com/oklog/run v1.1.0
	github.com/owncloud/flaex v0.0.0-20200411150708-dce59891a203
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
