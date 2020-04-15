module github.com/cosmos/modules/incubator/group

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200211145837-56c586897525
	github.com/cosmos/modules/incubator/orm v0.0.0-20200117100147-88228b5fa693
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/mux v1.7.3
	github.com/pkg/errors v0.9.1
	github.com/regen-network/cosmos-proto v0.1.1-0.20200213154359-02baa11ea7c2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.33.0
	github.com/tendermint/tm-db v0.4.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.1

//replace github.com/cosmos/modules/incubator/orm => github.com/regen-network/cosmos-modules/incubator/orm v0.0.0-20200206151518-3155fe39bfb9
replace github.com/cosmos/modules/incubator/orm => ../orm
