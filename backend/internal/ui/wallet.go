package ui

import (
	"errors"
	"fmt"

	"nft-engine/internal/wallet"

	"github.com/AlecAivazis/survey/v2"
)

var (
	ErrFailConfirm = errors.New("confirmation declined")
)

// the questions to ask
var addWalletQuestions = []*survey.Question{
	{
		Name: "mode",
		Prompt: &survey.Select{
			Message: "Import or Create a Wallet:",
			Options: []string{"Import Wallet", "New Wallet"},
			Default: "Import Wallet",
		},
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
	{
		Name:   "age",
		Prompt: &survey.Input{Message: "How old are you?"},
	},
}

type AddWalletForm struct {
	Mode       string
	Seed       string
	Passphrase string
}

// AddWalletPrompt prompts the user to import or create a wallet.
func AddWalletPrompt() (form *AddWalletForm, err error) {
	form = &AddWalletForm{}

	if err = survey.AskOne(&survey.Select{
		Message: "Import or Create a Wallet:",
		Options: []string{"Import Wallet", "New Wallet"},
		Default: "Import Wallet",
	}, &form.Mode); err != nil {
		return nil, err
	}

	switch form.Mode {
	case "Import Wallet":
		if err = survey.AskOne(&survey.Input{Message: "Input seed phrase to import.\n"}, &form.Seed, survey.WithValidator(survey.MinLength(40))); err != nil {
			return nil, err
		}

	case "New Wallet":
		form.Seed = wallet.NewMnemonic()
		confirm := false

		if err = survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("%s\n\nSAVE YOUR SEED PHRASE SOMEWHERE SECURE. PLEASE CONFIRM DO NOT CONTINUE UNTIL YOU HAVE IT SAVE!", form.Seed)}, &confirm); err != nil {
			return nil, err
		}
		if !confirm {
			return nil, ErrFailConfirm
		}

		confirm = false

		if err = survey.AskOne(&survey.Confirm{Message: "CONFIRM AGAIN."}, &confirm); err != nil {
			return nil, err
		}
		if !confirm {
			return nil, ErrFailConfirm
		}

	default:
		return nil, errors.New("unknown choice")
	}

	if err = survey.AskOne(&survey.Password{Message: "Input a passphrase between 8 - 14 characters: "}, &form.Passphrase); err != nil {
		return nil, err
	}

	match := false
	attempts := 0
	confirmation := ""

	for !match && attempts < 3 {
		if err = survey.AskOne(&survey.Password{Message: "Confirm passphrase: "}, &confirmation); err != nil {
			return nil, err
		}

		if confirmation == form.Passphrase {
			return
		}

		attempts++
	}

	return nil, errors.New("failed passphrase confirmation")
}
