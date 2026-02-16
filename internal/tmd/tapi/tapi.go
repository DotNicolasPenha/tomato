package tapi

import (
	"encoding/json"
	"fmt"
	"os"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/core/distros"
	"com.dotvinci.tm/internal/tmd/tapi/bases"
	"com.dotvinci.tm/internal/tmd/tapi/router"
)

type Tapi struct{}

func (Tapi) NameDistro() string {
	return "tapi-1.0"
}

type RouteBaseConfigs struct{}

func (Tapi) Exec(ctx distros.DistroExecContext) error {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Error to read you cwd in tapi-1.0")
	}
	router.Router(cwd, ctx)
	return nil
}

type ResponseMsgJSONBase struct{}
type ResponseMsgJson struct {
	Msg string `json:"msg"`
}

func (ResponseMsgJSONBase) Exec(ctx *bases.BaseContext) error {
	msgRaw := ctx.Route.BaseConfigs["msg"]
	if msgRaw == nil {
		return fmt.Errorf("base config 'msg' not found in %s", ctx.Route.Path)
	}
	ctx.Writter.Header().Set("Content-Type", "application/json")
	msg, ok := msgRaw.(string)
	if !ok {
		return fmt.Errorf("base config 'msg' must be a string")
	}
	response := ResponseMsgJson{Msg: msg}
	ctx.Writter.WriteHeader(200)
	if err := json.NewEncoder(ctx.Writter).Encode(response); err != nil {
		return err
	}
	return nil
}
func (ResponseMsgJSONBase) NameBase() string {
	return "response-json"
}
func init() {
	bases.RegistryBase(ResponseMsgJSONBase{})
	distros.Register(Tapi{})
}
