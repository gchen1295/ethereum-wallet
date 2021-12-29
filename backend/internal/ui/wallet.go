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
	mode       string
	seed       []byte
	passphrase []byte
}

func AddWalletPrompt() (form *AddWalletForm, err error) {
	form = &AddWalletForm{}
	modePrompt := &survey.Select{
		Message: "Import or Create a Wallet:",
		Options: []string{"Import Wallet", "New Wallet"},
		Default: "Import Wallet",
	}

	if err = survey.AskOne(modePrompt, &form.mode); err != nil {
		return nil, err
	}

	switch form.mode {
	case "Import Wallet":
		seed := ""
		if err = survey.AskOne(&survey.Input{Message: "Input seed phrase to import."}, &seed, survey.WithValidator(survey.MinLength(40))); err != nil {
			return nil, err
		}

	case "New Wallet":
		seed := wallet.NewMnemonic()
		confirm := false

		if err = survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("%s\n\nSAVE YOUR SEED PHRASE SOMEWHERE SECURE. PLEASE CONFIRM DO NOT CONTINUE UNTIL YOU HAVE IT SAVE!", seed)}, confirm); err != nil {
			return nil, err
		}
		if !confirm {
			return nil, ErrFailConfirm
		}

		confirm = false

		if err = survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("%s\n\nCONFIRM AGAIN.", seed)}, confirm); err != nil {
			return nil, err
		}
		if !confirm {
			return nil, ErrFailConfirm
		}

		form.seed = []byte(seed)

	default:
		return nil, errors.New("unknown choice")
	}

	if err = survey.AskOne(&survey.Input{Message: "Input a passphrase between 8 - 14 characters"}, &form.mode); err != nil {
		return nil, err
	}

	return
}
