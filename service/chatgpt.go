package service

import (
	"context"
	"fmt"
	"strings"

	util "github.com/restuwahyu13/gin-rest-api/utils"
	openai "github.com/sashabaranov/go-openai"
)

func ReviewCode(diffs []string) (string, error) {
	prePrompt := "As a senior developer, review the following code changes and answer code review questions about them. The code changes are provided as git diff strings:"
	questions := "\n\nQuestions:\n1. Can you summarise the changes in a succinct bullet point list\n2. In the diff, are the added or changed code written in a clear and easy to understand way?\n3. Does the code use comments, or descriptive function and variables names that explain what they mean?\n4. based on the code complexity of the changes, could the code be simplified without breaking its functionality? if so can you give example snippets?\n5. Can you find any bugs, if so please explain and reference line numbers?\n6. Do you see any code that could induce security issues?\n Plaese answer me in Chinese.\n"

	openai_key := util.GodotEnv("OPENAI_API_KEY")

	client := openai.NewClient(openai_key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a senior developer reviewing code changes.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("%s\n\n%s%s", prePrompt, strings.Join(diffs, ""), questions),
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "Format the response so it renders nicely in GitLab, with nice and organized markdown (use code blocks if needed), and send just the response no comments on the request, when answering include a short version of the question, so we know what it is.",
				},
			},
			Temperature: 0.7,
			Stream:      false,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	answer := resp.Choices[0].Message.Content
	fmt.Println(answer)

	return answer, nil
}
