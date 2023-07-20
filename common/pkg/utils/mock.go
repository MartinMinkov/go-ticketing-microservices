package utils

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
)

/*
Generated from:
auth-service-jwt
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJpYXQiOjE2ODk4ODk2MTEsImlkIjoiNjRiOWFiNGJjZmVhNmI2MmMxZTIyMjlmIn0.5kE_BCrTi7TwVMAyAN8hL1t_SzMIcbyP6wYJdWd1sXo
Path=/
HttpOnly
Expires=Fri, 21 Jul 2023 09:46:51 GMT

	{
			"email": "test@test.com",
			"password": "password123"
	}
*/
const User1JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJpYXQiOjE2ODk4ODk2MTEsImlkIjoiNjRiOWFiNGJjZmVhNmI2MmMxZTIyMjlmIn0.5kE_BCrTi7TwVMAyAN8hL1t_SzMIcbyP6wYJdWd1sXo"
const User2JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiaWF0IjoxNjg5ODg5NDQ5LCJpZCI6IjY0YjlhYWE5Y2ZlYTZiNjJjMWUyMjI5ZCJ9.iomMo7sgB-_EOjPAto7UhOTMjGFVO3iry8ND60zD75s"
const User3JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQHRlc3QuY29tIiwiaWF0IjoxNjg5ODg5NzAxLCJpZCI6IjY0YjlhYmE1Y2ZlYTZiNjJjMWUyMjJhMSJ9.IgByMJ6ybTt8nn81WHI2RI6SPvz_ivgbNJvialhKAIQ"

func MockAuthenticatedCookie(jwt string) *http.Cookie {
	return &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: jwt,
	}
}

func MockUnauthenticatedCookie() *http.Cookie {
	return &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: "",
	}
}
