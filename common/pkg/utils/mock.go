package utils

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
)

type UserMock struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	JWT       string `json:"jwt"`
}

var UserMock1 = UserMock{
	ID:        "64b9af66100af5c9c84eeb1e",
	Email:     "test1@test.com",
	CreatedAt: "2023-07-20T22:04:22.606443547Z",
	UpdatedAt: "2023-07-20T22:04:22.606443587Z",
	JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiaWF0IjoxNjg5ODkwNjYyLCJpZCI6IjY0YjlhZjY2MTAwYWY1YzljODRlZWIxZSJ9.cFifyvhanOWxFxnmi14xKCfcPjF3cDnG_uOyL8bOUVo",
}

var UserMock2 = UserMock{
	ID:        "64b9b00a100af5c9c84eeb20",
	Email:     "test2@test.com",
	CreatedAt: "2023-07-20T22:07:06.698012567Z",
	UpdatedAt: "2023-07-20T22:07:06.698012617Z",
	JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQHRlc3QuY29tIiwiaWF0IjoxNjg5ODkwODI2LCJpZCI6IjY0YjliMDBhMTAwYWY1YzljODRlZWIyMCJ9.xy5fWppuDwv6vphdEx4FX904MIYAEI0fJl7PGlHSxs8",
}

var UserMock3 = UserMock{
	ID:        "64b9b042100af5c9c84eeb22",
	Email:     "test3@test.com",
	CreatedAt: "2023-07-20T22:08:02.703566481Z",
	UpdatedAt: "2023-07-20T22:08:02.703566521Z",
	JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QzQHRlc3QuY29tIiwiaWF0IjoxNjg5ODkwODgyLCJpZCI6IjY0YjliMDQyMTAwYWY1YzljODRlZWIyMiJ9.g_vVeR7FZKQG_GmMsGI604qOJ7Ii_wYHLJXQmVzBjoo",
}

func SetUserCookieOnRequest(req *http.Request, user *UserMock) {
	cookie := &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: user.JWT,
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
}
