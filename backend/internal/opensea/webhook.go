package piratesea

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Bot) sendWebhook(embeds []Embed) error {
	webhook := Webhook{}
	webhook.Username = "OpenC"
	webhook.AvatarURL = "https://lh3.googleusercontent.com/5qyDayRBJP1YUQnlE7Z6O-R64CSNdQiUPalWtzvSlmX_KlrHtlcBwccaQdQovnt7Ymmes-vuHczK1rygoMM585mL2IIrCEqKy9sJNg=s0"
	webhook.Embeds = embeds

	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	res, err := c.request.Post("https://discord.com/api/webhooks/738893176194072636/JJXUz8SZnlJjlNsTke2YGsA3g-n-Fm8af4A-AKPXobT8M7h03x9lILFP99cevEvslPRt", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return errors.New("Webhook failed with statuscode: " + fmt.Sprint(res.StatusCode))
	}

	return nil
}

type Webhook struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}
type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}
type Fields struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty,omitempty"`
}
type Thumbnail struct {
	URL string `json:"url,omitempty"`
}
type Image struct {
	URL string `json:"url,omitempty"`
}
type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}
type Embed struct {
	Author      Author    `json:"author,omitempty"`
	Title       string    `json:"title,omitempty"`
	URL         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Color       int       `json:"color,omitempty"`
	Fields      []Fields  `json:"fields,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Image       Image     `json:"image,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
}
