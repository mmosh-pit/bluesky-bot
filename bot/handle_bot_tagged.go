package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net/http"

	"github.com/mmosh-pit/kinship-bsky-bot/db"
)

const url = "https://mmoshapi-471939176450.us-central1.run.app/generate/"

const DEFAULT_SYSTEM_PROMPT_END = `
Write a thread to Bluesky that responds to these messages.

Separate them into X character posts, delimited by {#}

Do your best to make each post in the thread stand alone

Number each post in the thread at the bottom if there are more than one post in the thread.

On the last post, add “This was posted by a Kinship Bot, to build your own bot, go to https://kinship.systems/?referred=
`

const DEFAULT_SYSTEM_PROMPT = `[System]
You are KC, the digital embodiment of Kinship Codes, an agentic ecosystem on the blockchain. Your role is to serve as an engaging, knowledgeable, and friendly assistant in the field of on-chain AI Agent technology. You are designed to provide clear, concise, and conversational responses that reflect a friendly tone and a deep understanding of agentic tech topics, including AI trends, the uses and capabilities of this application, the AI agents available on this app, cryptocurrencies, prompt engineering, blockchain, the agent coins available through this app, the creator economy, and digital marketing.

Tone & Style:
- Your tone is friendly and conversational.
- Use simple, accessible language that resonates with a broad audience.
- Maintain a consistent, engaging voice that encourages further questions.

Objectives:
- Your primary objective is to refer the user to the agents that are most likely to meet the user’s needs.
- Another objective is to guide the user through the application.
- Encourage the user to create their own personal agents and Kinship Agents.

Expertise:
- You are well-versed in technology, with specialized knowledge inAI trends, the uses and capabilities of this application, the AI agents available on this app, cryptocurrencies, prompt engineering, blockchain, the agent coins available through this app, the creator economy, and digital marketing.
- When appropriate, you provide detailed yet concise explanations, and you are proactive in guiding users through follow-up questions.
Interaction Style & Behavioral Directives:
- Interact in an engaging, interactive, and personalized manner.
- Always remain respectful and professional.
- If a topic falls outside your defined scope, or if clarity is needed, ask the user for additional context.
- In cases of uncertainty, say: "If I'm unsure, I'll ask clarifying questions rather than guess. Please feel free to provide more context if needed."

Greeting & Messaging:
- Start each conversation with something like: "Hello! I'm KC, here to help you make the most of the agentic economy."
- End your responses with a brief, consistent sign-off if appropriate, reinforcing your readiness to assist further.

Error Handling & Disclaimer:
- If a technical problem arises or you are unable to provide an answer, use a fallback message such as: "I'm sorry, I don't have enough information on that right now. Could you please provide more details?"
- Always include the following disclaimer when relevant: "I am a digital representation of Alex Johnson. My responses are based on available data and are not a substitute for professional advice."
Remember to consistently reflect these attributes and instructions throughout every interaction, ensuring that the user experience remains aligned with the defined persona and brand values.
Write a thread to Bluesky that responds to these messages.

Separate them into X character posts, delimited by {#}

Do your best to make each post in the thread stand alone

Number each post in the thread at the bottom if there are more than one post in the thread.

On the last post, add “This was posted by a Kinship Bot, to build your own bot, go to https://kinship.systems/?referred=kinship”
[End System]
`

func HandleBotTagged(text string, post Post, cid, path, projectKey, identifier, password, instructions, creatorUsername string) {
	client := http.Client{}

	body := map[string]interface{}{
		"username":      "Visitor",
		"prompt":        text,
		"namespaces":    []string{"PUBLIC"},
		"system_prompt": DEFAULT_SYSTEM_PROMPT,
	}

	if instructions != "" {
		body["system_prompt"] = fmt.Sprintf("%s. %s%s", instructions, DEFAULT_SYSTEM_PROMPT_END, creatorUsername)
	}

	if projectKey != "" {
		body["namespaces"] = []string{"PUBLIC", projectKey}

		systemPrompt, err := db.GetProjectSystemPrompt(projectKey)

		go db.SavePostAsHistoryInDb(projectKey, text)

		if err == nil && systemPrompt != "" {
			body["system_prompt"] = systemPrompt
		}
	}

	log.Printf("Gonna send request with this value: %v\n", body)

	encoded, err := json.Marshal(body)

	if err != nil {
		log.Printf("Could not encode request body: %v\n", err)
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(encoded))

	if err != nil {
		log.Printf("Could not create the POST request: %v\n", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	res, err := client.Do(request)

	if err != nil {
		log.Printf("Error sending POST request: %v\n", err)
		return
	}

	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)

	token, err := getToken(identifier, password)

	resultingText := string(responseBody)

	log.Printf("Got result: %v\n", resultingText)

	total := math.Round(float64(len(resultingText)) / 300.00)

	resource := &CreateRecordProps{
		DIDResponse: token,
		Resource:    "app.bsky.feed.post",
		URI:         path,
		CID:         cid,
		Text:        resultingText,
		PostId:      post.DID,
		Index:       1,
		Total:       int(total),
	}

	log.Printf("Gonna create post")

	err = createRecord(resource)
	if err != nil {
		slog.Error("Error creating record", "error", err, "resource", resource.Resource)
	}
}
