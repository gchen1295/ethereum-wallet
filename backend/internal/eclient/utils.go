package eclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"nft-engine/pkg/utils"

	"github.com/ethereum/go-ethereum/core/types"
)

type (
	Webhook   = utils.Webhook
	Embed     = utils.Embed
	Author    = utils.Author
	Fields    = utils.Fields
	Thumbnail = utils.Thumbnail
	Image     = utils.Image
	Footer    = utils.Footer
)

func (c *Client) SendWebhook(webhookUrl string, block *types.Header) error {
	webhook := Webhook{}
	webhook.Username = "Paws"
	webhook.AvatarURL = "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Ethereum-icon-purple.svg/1200px-Ethereum-icon-purple.svg.png"
	webhook.Embeds = []Embed{}
	webhook.Embeds = append(webhook.Embeds, Embed{
		Author: Author{
			Name: fmt.Sprintf("New Block - %s", block.Number.String()),
			URL:  fmt.Sprintf("https://etherscan.io/block/%s", block.Number.String()),
		},
		Color: 0x00ff00,
		Footer: Footer{
			Text: "Made with ❤️ by @undefined#6969",
		},
		Fields: []Fields{
			{
				Name:   "Gas Limit",
				Value:  fmt.Sprintf("%d", block.GasLimit),
				Inline: true,
			},
			{
				Name:   "Base Fee",
				Value:  block.BaseFee.String(),
				Inline: true,
			},
			{
				Name:   "Gas Used",
				Value:  fmt.Sprintf("%d", block.GasUsed),
				Inline: true,
			},
			{
				Name:   "Parent Hash",
				Value:  block.ParentHash.String(),
				Inline: false,
			},
			{
				Name:   "Txn Hash",
				Value:  block.TxHash.String(),
				Inline: false,
			},
			{
				Name:   "Time",
				Value:  fmt.Sprintf("<t:%d>", block.Time),
				Inline: true,
			},
			{
				Name:   "Timestamp",
				Value:  fmt.Sprintf("%d", block.Time),
				Inline: true,
			},
		},
	})

	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	res, err := c.request.Post(webhookUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return errors.New("Webhook failed with statuscode: " + fmt.Sprint(res.StatusCode))
	}

	return nil
}
