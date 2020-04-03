module github.com/cosmos/modules/incubator/group

go 1.13

require (
	github.com/99designs/keyring v1.1.4 // indirect
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200401144647-f7c0578fea8f
	github.com/cosmos/modules/incubator/orm v0.0.0-20200117100147-88228b5fa693
	github.com/gibson042/canonicaljson-go v1.0.3 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.4.3 // indirect
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/otiai10/copy v1.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7 // indirect
	github.com/regen-network/cosmos-proto v0.1.1-0.20200213154359-02baa11ea7c2
	github.com/spf13/cobra v0.0.7
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/tendermint v0.33.2
	github.com/tendermint/tm-db v0.5.1
	google.golang.org/protobuf v1.20.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.1

//replace github.com/cosmos/modules/incubator/orm => github.com/regen-network/cosmos-modules/incubator/orm v0.0.0-20200206151518-3155fe39bfb9
replace github.com/cosmos/modules/incubator/orm => ../orm
