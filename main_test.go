package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/steinfletcher/apitest"
	json "github.com/steinfletcher/apitest-jsonpath"
	"rakamin.com/project/config"
	"rakamin.com/project/controllers"
	"rakamin.com/project/models"
)

var (
	jwt_token_user1 string
	jwt_token_user2 string
	msg_response    models.MessagesWithConversation
	reply_response  models.MessagesWithConversation
)

func getToken(id int) string {
	var user models.Users
	if id == 1 {
		user = models.Users{
			ID:       1,
			Name:     "User 1",
			UserName: "user1",
		}
	} else if id == 2 {
		user = models.Users{
			ID:       2,
			Name:     "User 2",
			UserName: "user2",
		}
	} else if id == 3 {
		user = models.Users{
			ID:       3,
			Name:     "User 3",
			UserName: "user3",
		}
	} else {
		return ""
	}
	return controllers.GetToken(&user)
}

func beforeEach() {
	// Setup Database
	config.Initialization()
	models.SetupModels()

	tkn1 := fmt.Sprintf("Bearer %s", getToken(1))
	tkn2 := fmt.Sprintf("Bearer %s", getToken(2))
	jwt_token_user1 = tkn1
	jwt_token_user2 = tkn2
}

func TestAPIUptime(t *testing.T) {
	beforeEach()
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
	beforeEach()
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
	beforeEach()

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

func TestNewMessages(t *testing.T) {
	beforeEach()

	t.Run("Not authenticated", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation").
			JSON(`{"user_id": 1,"message":"messages"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Unauthorized").
					End(),
			).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Missing user_id", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation").
			Header("Authorization", jwt_token_user1).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Missing user_id").
					End(),
			).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Cannot chat yourself", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation").
			Header("Authorization", jwt_token_user1).
			JSON(`{"user_id": 1,"message":"test from testing"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Invalid user_id. Cannot chat yourself").
					End(),
			).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("User not found", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation").
			Header("Authorization", jwt_token_user1).
			JSON(`{"user_id": 5,"message":"test from testing"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Invalid user_id. User not found").
					End(),
			).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Empty messages", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation").
			Header("Authorization", jwt_token_user1).
			JSON(`{"user_id": 2,"message":""}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Message cannot be empty").
					End(),
			).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Success", func(t *testing.T) {
		response := apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation").
			Header("Authorization", jwt_token_user1).
			JSON(`{"user_id": 2,"message":"test from testing"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("message", "test from testing").
					Equal("sender.username", "user1").
					Present("conversation").
					Equal("read_status", false).
					Present("timestamp").
					End(),
			).
			Status(http.StatusOK).
			End()
		response.JSON(&msg_response)
	})
}

func TestListConversation(t *testing.T) {
	beforeEach()

	t.Run("User 1", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation").
			Header("Authorization", jwt_token_user1).
			Expect(t).
			Assert(
				json.Chain().
					Equal("data[0].unread", float64(0)).
					Present("data[0].users").
					Equal("data[0].message.message", msg_response.Messages.Messages).
					Present("data[0].message.sender").
					Equal("data[0].message.sender.username", "user1").
					End(),
			).
			Status(http.StatusOK).
			End()
	})

	t.Run("User 2", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation").
			Header("Authorization", jwt_token_user2).
			Expect(t).
			Assert(
				json.Chain().
					Equal("data[0].unread", float64(1)).
					Present("data[0].users").
					Equal("data[0].message.message", msg_response.Messages.Messages).
					Present("data[0].message.sender").
					Equal("data[0].message.sender.username", "user1").
					End(),
			).
			Status(http.StatusOK).
			End()
	})
}

func TestListMessages(t *testing.T) {
	beforeEach()

	t.Run("Forbidden conversation", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation/22").
			Header("Authorization", jwt_token_user2).
			Expect(t).
			Assert(
				json.Len("data", 0),
			).
			Status(http.StatusOK).
			End()
	})

	t.Run("User 1", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation/11").
			Header("Authorization", jwt_token_user1).
			Expect(t).
			Assert(
				json.Chain().
					Equal("data[0].read_status", false).
					Equal("data[0].message", msg_response.Messages.Messages).
					Present("data[0].sender").
					Equal("data[0].sender.username", "user1").
					End(),
			).
			Status(http.StatusOK).
			End()
	})

	t.Run("User 2", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation/11").
			Header("Authorization", jwt_token_user2).
			Expect(t).
			Assert(
				json.Chain().
					Equal("data[0].read_status", true).
					Equal("data[0].message", msg_response.Messages.Messages).
					Present("data[0].sender").
					Equal("data[0].sender.username", "user1").
					End(),
			).
			Status(http.StatusOK).
			End()
	})
}

func TestReplyMessages(t *testing.T) {
	beforeEach()

	t.Run("Not authenticated", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation/11").
			JSON(`{"user_id": 1,"message":"test reply from testing"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Unauthorized").
					End(),
			).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Conversation not found", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation/22").
			Header("Authorization", jwt_token_user2).
			JSON(`{"message":"test reply from testing"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Conversation not found").
					End(),
			).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Empty messages", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation/11").
			Header("Authorization", jwt_token_user2).
			JSON(`{"message":""}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("error", "Message cannot be empty").
					End(),
			).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Success", func(t *testing.T) {
		response := apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Post("/conversation/11").
			Header("Authorization", jwt_token_user2).
			JSON(`{"user_id": 2,"message":"test reply from testing"}`).
			Expect(t).
			Assert(
				json.Chain().
					Equal("message", "test reply from testing").
					Equal("sender.username", "user2").
					Present("conversation").
					Equal("read_status", false).
					Present("timestamp").
					End(),
			).
			Status(http.StatusOK).
			End()
		response.JSON(&reply_response)
	})
}

func TestListConversationAfterReplied(t *testing.T) {
	beforeEach()

	t.Run("User 1", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation").
			Header("Authorization", jwt_token_user1).
			Expect(t).
			Assert(
				json.Chain().
					Equal("data[0].unread", float64(1)).
					Present("data[0].users").
					Equal("data[0].message.message", reply_response.Messages.Messages).
					Present("data[0].message.sender").
					Equal("data[0].message.sender.username", "user2").
					End(),
			).
			Status(http.StatusOK).
			End()
	})

	t.Run("User 2", func(t *testing.T) {
		apitest.New().
			HandlerFunc(FiberToHandlerFunc(newApp())).
			Get("/conversation").
			Header("Authorization", jwt_token_user2).
			Expect(t).
			Assert(
				json.Chain().
					Equal("data[0].unread", float64(0)).
					Present("data[0].users").
					Equal("data[0].message.message", reply_response.Messages.Messages).
					Present("data[0].message.sender").
					Equal("data[0].message.sender.username", "user2").
					End(),
			).
			Status(http.StatusOK).
			End()
	})
}

func TestDeleteMessage(t *testing.T) {
	beforeEach()

	if msg_response.ID != 0 {
		t.Run("(Dev) Delete New Messages", func(t *testing.T) {
			apitest.New().
				HandlerFunc(FiberToHandlerFunc(newApp())).
				Delete(fmt.Sprintf("/conversation/%d/%d", msg_response.ConversationID, msg_response.ID)).
				Header("Authorization", jwt_token_user1).
				Expect(t).
				Assert(
					json.Chain().
						Equal("success", true).
						End(),
				).
				Status(http.StatusOK).
				End()
		})
	}

	if reply_response.ID != 0 {
		t.Run("(Dev) Delete Reply Messages", func(t *testing.T) {
			apitest.New().
				HandlerFunc(FiberToHandlerFunc(newApp())).
				Delete(fmt.Sprintf("/conversation/%d/%d", reply_response.ConversationID, reply_response.ID)).
				Header("Authorization", jwt_token_user2).
				Expect(t).
				Assert(
					json.Chain().
						Equal("success", true).
						End(),
				).
				Status(http.StatusOK).
				End()
		})
	}
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
