package ui

import (
	"errors"
	"fmt"
	"nft-engine/internal/eclient"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var ContractAddressRegex = eclient.ContractAddressRegex

type TaskForm struct {
	ContractAddress *common.Address
	MintFunction    *abi.Method
	Args            []*ArgumentInput
	PayableAmount   string
}

type ArgumentInput struct {
	*abi.Argument
	Value string
}

// validateContractAddress uses a simple regex validator to validate the user's input.
func validateContractAddress(val interface{}) error {
	if str, ok := val.(string); !ok || !ContractAddressRegex.MatchString(str) {
		return errors.New("input is not a valid contract address")
	}

	return nil
}

// ContractAddressPrompt prompts the user for a contract address.
func ContractAddressPrompt(form *TaskForm) error {
	input := ""
	if err := survey.AskOne(&survey.Input{Message: "Contract Address: "}, &input, survey.WithValidator(validateContractAddress)); err != nil {
		return err
	}

	address := common.HexToAddress(input)
	form.ContractAddress = &address

	return nil
}

// SelectMintFunctionPrompt prompts the user to select a mint function from the given contract.
func SelectMintFunctionPrompt(contract *eclient.Contract, form *TaskForm) error {
	if contract == nil {
		return errors.New("cannot pass nil contract")
	}

	methods := contract.GetMethodNames()
	if len(methods) == 0 {
		return errors.New("contract has no methods")
	}

	var fnName string
	var method *abi.Method
	var inputsValues []*ArgumentInput

	if err := survey.AskOne(&survey.Select{
		Message: "Select mint function",
		Options: methods,
		Default: methods[0],
	}, &fnName); err != nil {
		return err
	}
	payable := ""
	method = contract.GetMethodByName(fnName)
	if method.IsPayable() {
		if err := survey.AskOne(&survey.Input{Message: "Payable Amount: "}, &payable, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	for _, input := range method.Inputs {
		val := ""
		if err := survey.AskOne(&survey.Input{Message: fmt.Sprintf("%s (%s): ", input.Name, input.Type.String())}, &val, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		inputsValues = append(inputsValues, &ArgumentInput{
			Argument: &input,
			Value:    val,
		})
	}

	form.ContractAddress = contract.Address
	form.MintFunction = method
	form.Args = inputsValues
	form.PayableAmount = payable

	return nil
}
