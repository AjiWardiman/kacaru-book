package main

import (
  "context"
  "fmt"
  "os"
  openrouter "github.com/OpenRouterTeam/go-sdk"
  "github.com/OpenRouterTeam/go-sdk/models/components"
  "log"
)

func main() {
  ctx := context.Background()

  s := openrouter.New(
    openrouter.WithSecurity(os.Getenv("OPENROUTER_API_KEY")),
  )

  res, err := s.Chat.Send(ctx, components.ChatRequest{
    Messages: []components.ChatMessages{
      components.CreateChatMessagesUser(
        components.ChatUserMessage{
          Content: components.CreateChatUserMessageContentStr(
            "What is the meaning of life?",
          ),
          Role: components.ChatUserMessageRoleUser,
        },
      ),
    },
    Model: openrouter.Pointer("openai/gpt-4o-mini"),
  })
  if err != nil {
    log.Fatal(err)
  }
  if res.ChatResult != nil {
    fmt.Println(res.ChatResult.Choices[0].Message.Content)
  }
}