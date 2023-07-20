package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"auth.mminkov.net/api"
	"auth.mminkov.net/internal/config"
	"auth.mminkov.net/internal/state"
	"github.com/gin-gonic/gin"
)

type TestApp struct {
	Config   *config.Config
	AppState *state.AppState
	Server   *http.Server
}

func (t *TestApp) GET_healthcheck() *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	resp, err := http.Get(addr + "/api/users/healthcheck")
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) GET_currentuser(cookie *http.Cookie) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodGet, addr+"/api/users/currentuser", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) POST_sign_up(data map[string]interface{}) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(addr+"/api/users/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) POST_sign_in(data map[string]interface{}, cookie *http.Cookie) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, addr+"/api/users/signin", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) POST_sign_out(cookie *http.Cookie) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodPost, addr+"/api/users/signout", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func SpawnApp() *TestApp {
	applicationConfig := config.BuildApplicationConfig()
	databaseConfig := config.BuildDatabaseConfig()
	databaseConfig.Database = "auth_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	config := &config.Config{
		ApplicationConfig: applicationConfig,
		DatabaseConfig:    databaseConfig,
	}

	appState := api.BuildAppState(config)

	gin.SetMode(gin.ReleaseMode)
	server := api.BuildServer(config, appState)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait a bit for the server to start
	time.Sleep(500 * time.Millisecond)

	return &TestApp{
		Config:   config,
		AppState: appState,
		Server:   server,
	}
}

func (app *TestApp) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer app.AppState.DBCleanup()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}

func ReadJSON(r *http.Response, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("body contains badly-formed JSON (at character " + strconv.Itoa(int(syntaxError.Offset)) + ")")

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return errors.New("body contains incorrect JSON type for field " + unmarshalTypeError.Field)
			}
			return errors.New("body contains incorrect JSON type (at character " + strconv.Itoa(int(unmarshalTypeError.Offset)) + ")")

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return errors.New("body contains unknown JSON error")

		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
