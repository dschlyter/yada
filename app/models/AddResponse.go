package models

type AddResponse struct {
    Error *string
}

func ReturnSuccess() AddResponse {
    return AddResponse{}
}

func ReturnError(errMsg string) AddResponse {
    return AddResponse{Error: &errMsg}
}
