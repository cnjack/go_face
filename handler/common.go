package handler

type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewCommonRespWithError(code int, err error) *CommonResp {
	return &CommonResp{
		Code: code,
		Msg:  err.Error(),
	}
}

func NewCommonRespWithData(data interface{}) *CommonResp {
	return &CommonResp{
		Code: 0,
		Msg:  "",
		Data: data,
	}
}
