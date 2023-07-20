package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ginHands struct {
	SerName         string
	Path            string
	Latency         time.Duration
	Method          string
	StatusCode      int
	ClientIP        string
	MsgStr          string
	UserAgent       string
	ResponseSize    int
	RequestHeaders  map[string][]string
	ResponseHeaders map[string][]string
}

func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(typ).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

func Logger(serName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		reqHeaders := c.Request.Header
		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}

		userAgent := c.GetHeader("User-Agent")
		responseSize := c.Writer.Size()
		respHeaders := c.Writer.Header()

		cData := &ginHands{
			SerName:         serName,
			Path:            path,
			Latency:         time.Duration(time.Since(t).Microseconds()),
			Method:          c.Request.Method,
			StatusCode:      c.Writer.Status(),
			ClientIP:        c.ClientIP(),
			MsgStr:          msg,
			UserAgent:       userAgent,
			ResponseSize:    responseSize,
			RequestHeaders:  reqHeaders,
			ResponseHeaders: respHeaders,
		}

		logSwitch(cData)
	}
}

func logSwitch(data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			log.Warn().
				Str("ser_name", data.SerName).
				Str("method", data.Method).
				Str("path", data.Path).
				Dur("resp_time", data.Latency).
				Int("status", data.StatusCode).
				Str("client_ip", data.ClientIP).
				Str("user_agent", data.UserAgent).
				Int("response_size", data.ResponseSize).
				Interface("request_headers", data.RequestHeaders).
				Interface("response_headers", data.ResponseHeaders).
				Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			log.Error().
				Str("ser_name", data.SerName).
				Str("method", data.Method).
				Str("path", data.Path).
				Dur("resp_time", data.Latency).
				Int("status", data.StatusCode).
				Str("client_ip", data.ClientIP).
				Str("user_agent", data.UserAgent).
				Int("response_size", data.ResponseSize).
				Interface("request_headers", data.RequestHeaders).
				Interface("response_headers", data.ResponseHeaders).
				Msg(data.MsgStr)
		}
	default:
		log.Info().
			Str("ser_name", data.SerName).
			Str("method", data.Method).
			Str("path", data.Path).
			Dur("resp_time", data.Latency).
			Int("status", data.StatusCode).
			Str("client_ip", data.ClientIP).
			Str("user_agent", data.UserAgent).
			Int("response_size", data.ResponseSize).
			Interface("request_headers", data.RequestHeaders).
			Interface("response_headers", data.ResponseHeaders).
			Msg(data.MsgStr)
	}
}
