package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/steinfletcher/apitest"
	json "github.com/steinfletcher/apitest-jsonpath"
)

func TestAPIUptime(t *testing.T) {
	apitest.New().
		HandlerFunc(FiberToHandlerFunc(newApp())).
		Get("/").
		Expect(t).
		Assert(
			json.Chain().
				Equal("error", false).
				Equal("message", "API Uptime").
				End(),
		).
		Status(http.StatusOK).
		End()
}

func TestAPINotFound(t *testing.T) {
	apitest.New().
		HandlerFunc(FiberToHandlerFunc(newApp())).
		Get("/login").
		Expect(t).
		Assert(
			json.Chain().
				Equal("error", "Not found").
				End(),
		).
		Status(http.StatusNotFound).
		End()
}

func TestLogin(t *testing.T) {
	t.Run("Missing username/password", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/login").
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Missing username/password").
					End(),
			).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Invalid username", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/login").
			JSON(`{"username": "admin","password":"admin"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Invalid username").
					End(),
			).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Invalid password", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/login").
			JSON(`{"username": "user1","password":"admin"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Invalid password").
					End(),
			).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Success", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/login").
			JSON(`{"username": "user1","password":"user1"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("user.name", "User 1").
					Equal("user.username", "user1").
					Present("token").
					End(),
			).
			Status(http.StatusOK).
			End()
	})
}

func FiberToHandlerFunc(app *fiber.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := app.Test(r)
		if err != nil {
			panic(err)
		}

		// copy headers
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)

		if _, err := io.Copy(w, resp.Body); err != nil {
			panic(err)
		}
	}
}
