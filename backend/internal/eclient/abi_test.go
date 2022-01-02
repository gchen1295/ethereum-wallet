package eclient

import (
	"io/ioutil"
	"testing"
)

var ABIPATH = `.\contracts\treasureMarket.json`

func TestABI(t *testing.T) {
	rawData, err := ioutil.ReadFile(ABIPATH)
	if err != nil {
		t.Fatal(err)
	}

	contractABI, err := ParseABI(rawData)
	if err != nil {
		t.Fatal(err)
	}

	if len(contractABI) != 33 {
		t.Fatalf("Error parsing abi!\nexpected length: %d\nactual length:%d\n", 33, len(contractABI))
	}

	contractMap := GetMethodMapping(contractABI)
	if len(contractMap) != 33 {
		t.Fatalf("Error mapping abi!\nexpected length: %d\nactual length:%d\n", 33, len(contractMap))
	}

	for _, method := range contractABI {
		t.Log(method.Name)
		_, found := contractMap[method.Name]
		if !found {
			t.Fatalf("Error mapping abi!\nmissing method: %s", method.Name)
		}
	}
}
