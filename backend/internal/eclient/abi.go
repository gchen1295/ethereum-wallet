package eclient

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type IOTypes struct {
	Indexed      bool   `json:"indexed,omitempty"`
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

type ABIMethods struct {
	Inputs          []*IOTypes `json:"inputs"`
	StateMutability string     `json:"stateMutability,omitempty"`
	Type            string     `json:"type"`
	Anonymous       bool       `json:"anonymous,omitempty"`
	Name            string     `json:"name,omitempty"`
	Outputs         []*IOTypes `json:"outputs,omitempty"`
}

type ABIJson []*ABIMethods

type ABI map[string]*ABIMethods

// GetMethodMapping returns a mapping of the povided ABI for easier lookup.
func GetMethodMapping(abiArr ABIJson) ABI {
	mapping := ABI{}

	for _, method := range abiArr {
		mapping[method.Name] = method
	}

	return mapping
}

// ParseABI parses a raw ABI array into an ABIJSON struct.
func ParseABI(raw []byte) (ABIJson, error) {
	parsedAbi := ABIJson{}

	if err := json.Unmarshal(raw, &parsedAbi); err != nil {
		return nil, err
	}

	return parsedAbi, nil
}

// GetMethodByName retrieves a method by name or nil if none found.
func (contract *Contract) GetMethodByName(methodName string) *abi.Method {
	target, ok := contract.ABI.Methods[methodName]
	if !ok {
		return nil
	}

	return &target
}

// GetMethodNames retrieves the method names of the contract.
func (contract *Contract) GetMethodNames() []string {
	names := []string{}
	for i, _ := range contract.ABI.Methods {
		names = append(names, i)
	}

	return names
}
