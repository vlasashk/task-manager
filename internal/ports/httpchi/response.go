package httpchi

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResp struct {
	Param string `json:"param,omitempty"`
	Value string `json:"value,omitempty"`
	Error string `json:"error"`
}

type MsgResp struct {
	Msg string `json:"message"`
}

func NewErr(param, val, err string) ErrResp {
	return ErrResp{
		Param: param,
		Value: val,
		Error: err,
	}
}

func NewMsg(msg string) MsgResp {
	return MsgResp{
		Msg: msg,
	}
}

func (resp ErrResp) Send(w http.ResponseWriter, r *http.Request, status int) {
	render.Status(r, status)
	render.JSON(w, r, resp)
}
func (resp MsgResp) Send(w http.ResponseWriter, r *http.Request, status int) {
	render.Status(r, status)
	render.JSON(w, r, resp)
}
