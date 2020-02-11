module github.com/cosmos/modules/incubator/group

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200211145837-56c586897525
	github.com/cosmos/modules/incubator/orm v0.0.0-20200117100147-88228b5fa693
	github.com/gogo/protobuf v1.3.1
	google.golang.org/grpc v1.26.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.1

//replace github.com/cosmos/modules/incubator/orm => github.com/regen-network/cosmos-modules/incubator/orm v0.0.0-20200206151518-3155fe39bfb9
replace github.com/cosmos/modules/incubator/orm => ../orm