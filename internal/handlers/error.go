package handlers

import (
	"context"
	"encoding/json"

	"github.com/go-chi/chi/middleware"
)

const (
	ErrorReadHTTPBody = iota + 1
	ErrorUnmarshalHTTPBody
	ErrorMarshalHTTPBody
	ErrorCreateHouse
	ErrorParseURL
	ErrorGetFlatsByHouseID
	ErrorNotAuthorized
	ErrorRegisterUser
	ErrorLoginUser
	ErrorDummyLogin
	ErrorCreateFlat
	ErrorUpdateFlat
	ErrorSubscribeOnHouse
	ErrorNoAuthorized
	ErrorExtractRoleFromToken
)

const (
	ReadHTTPBodyMsg              = "can't read request"
	UnmarshalHTTPBodyMsg         = "can't unmarshal request"
	ErrorCreateHouseMsg          = "can't create house"
	ErrorMarshalHTTPBodyMsg      = "can't marshal response"
	ErrorParseURLMsg             = "can't parse url"
	ErrorGetFlatsByHouseIDMsg    = "can't get flats by house id"
	ErrorNotAuthorizedMsg        = "not authorized"
	ErrorRegisterUserMsg         = "can't register user"
	ErrorLoginUserMsg            = "can't login user"
	ErrorDummyLoginMsg           = "can't simple login"
	ErrorCreateFlatMsg           = "can't create flat"
	ErrorUpdateFlatMsg           = "can't update flat"
	ErrorSubscribeOnHouseMsg     = "can't subscribe on house"
	ErrorNoAccessMsg             = "no enough access rights"
	ErrorExtractRoleFromTokenMsg = "can't extract role"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func CreateErrorResponse(ctx context.Context, errCode int, msg string) []byte {
	var errResponse ErrorResponse
	errResponse.Code = errCode
	errResponse.RequestID = middleware.GetReqID(ctx)
	errResponse.Message = msg

	response, err := json.Marshal(errResponse)
	if err != nil {
		return nil
	}

	return response
}
