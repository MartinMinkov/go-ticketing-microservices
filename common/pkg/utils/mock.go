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
func MockAuthenticatedUser1() *http.Cookie {
	return &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJpYXQiOjE2ODk4ODk2MTEsImlkIjoiNjRiOWFiNGJjZmVhNmI2MmMxZTIyMjlmIn0.5kE_BCrTi7TwVMAyAN8hL1t_SzMIcbyP6wYJdWd1sXo",
	}
}

/*
Generated from:
auth-service-jwt
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiaWF0IjoxNjg5ODg5NDQ5LCJpZCI6IjY0YjlhYWE5Y2ZlYTZiNjJjMWUyMjI5ZCJ9.iomMo7sgB-_EOjPAto7UhOTMjGFVO3iry8ND60zD75s
Path=/
HttpOnly
Expires=Fri, 21 Jul 2023 09:44:09 GMT

	{
	    "email": "test1@test.com",
	    "password": "password123"
	}
*/
func MockAuthenticatedUser2() *http.Cookie {
	return &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxQHRlc3QuY29tIiwiaWF0IjoxNjg5ODg5NDQ5LCJpZCI6IjY0YjlhYWE5Y2ZlYTZiNjJjMWUyMjI5ZCJ9.iomMo7sgB-_EOjPAto7UhOTMjGFVO3iry8ND60zD75s",
	}
}

/*
Generated From:
auth-service-jwt
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQHRlc3QuY29tIiwiaWF0IjoxNjg5ODg5NzAxLCJpZCI6IjY0YjlhYmE1Y2ZlYTZiNjJjMWUyMjJhMSJ9.IgByMJ6ybTt8nn81WHI2RI6SPvz_ivgbNJvialhKAIQ
Path=/
HttpOnly
Expires=Fri, 21 Jul 2023 09:48:21 GMT

	{
	    "email": "test2@test.com",
	    "password": "password123"
	}
*/
func MockAuthenticatedUser3() *http.Cookie {
	return &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQHRlc3QuY29tIiwiaWF0IjoxNjg5ODg5NzAxLCJpZCI6IjY0YjlhYmE1Y2ZlYTZiNjJjMWUyMjJhMSJ9.IgByMJ6ybTt8nn81WHI2RI6SPvz_ivgbNJvialhKAIQ",
	}
}
