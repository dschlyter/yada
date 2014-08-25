package models

type AddResponse struct {
    Error   *string
    Success bool
}

func ReturnSuccess() AddResponse {
    return AddResponse{Success: true}
}

func ReturnError(errMsg string) AddResponse {
    return AddResponse{Error: &errMsg}
}
