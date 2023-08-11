package service

import (
	"fmt"
	"net/http"

	util "github.com/restuwahyu13/gin-rest-api/utils"

	"github.com/gin-gonic/gin"
)

type CodeReviewRequest struct {
	ObjectKind        string           `json:"object_kind"`
	EventName         string           `json:"event_name"`
	Before            string           `json:"before"`
	After             string           `json:"after"`
	Ref               string           `json:"ref"`
	CheckoutSHA       string           `json:"checkout_sha"`
	Message           string           `json:"message"`
	UserID            int              `json:"user_id"`
	UserName          string           `json:"user_name"`
	UserEmail         string           `json:"user_email"`
	UserAvatar        string           `json:"user_avatar"`
	ProjectID         int              `json:"project_id"`
	Project           ProjectInfo      `json:"project"`
	ObjectAttributes  ObjectAttributes `json:"object_attributes"`
	Commits           []CommitInfo     `json:"commits"`
	TotalCommitsCount int              `json:"total_commits_count"`
	PushOptions       PushOptions      `json:"push_options"`
}

type ProjectInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
}

type CommitInfo struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	URL       string     `json:"url"`
	Author    AuthorInfo `json:"author"`
}

type AuthorInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PushOptions struct {
	CI CIOptions `json:"ci"`
}

type CIOptions struct {
	Skip bool `json:"skip"`
}

type ObjectAttributes struct {
	Action string `json:"action"`
	Iid    int    `json:"iid"`
}

// webhook handler
func Webhook(c *gin.Context) {

	requestBody := CodeReviewRequest{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Request.Header.Get("X-Gitlab-Token") != util.GodotEnv("WEBHOOK_VERIFY_TOKEN") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
		})
		return
	}

	if requestBody.ObjectKind == "merge_request" {
		if requestBody.ObjectAttributes.Action != "open" {
			c.JSON(http.StatusOK, gin.H{
				"status": "Not a  PR open",
			})
		}

		fmt.Println(requestBody)
		projectId := requestBody.Project.ID
		//projectCommitId := requestBody.Commits
		mrId := requestBody.ObjectAttributes.Iid

		diffs, err := GetChanges(projectId, mrId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"get gitlab changes error": err.Error()})
			return
		}

		comments, err := ReviewCode(diffs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"chatgpt review code error": err.Error()})
			return
		}

		if err := WriteComments(projectId, mrId, comments); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"write gitlab comments error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "review complete!",
	})
}

// func test(){
// 	openai_key := util.GodotEnv("OPENAI_API_KEY")

// 	client := openai.NewClient(openai_key)
// 	resp, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model: openai.GPT3Dot5Turbo16K,
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    openai.ChatMessageRoleUser,
// 					Content: "Hello!",
// 				},
// 			},
// 		},
// 	)

// 	if err != nil {
// 		fmt.Printf("ChatCompletion error: %v\n", err)
// 		return
// 	}

// 	fmt.Println(resp.Choices[0].Message.Content)

// 	c.JSON(200, gin.H{
// 		"message": "review complete!",
// 		"result":  resp.Choices[0].Message.Content,
// 	})
// }
