package main

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/migratooor/tokenLists/generators/common/helpers"
)

// TokenListToken is the token struct used in the default token list
type TokenListToken struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	LogoURI  string `json:"logoURI"`
	ChainID  uint64 `json:"chainId"`
	Decimals int    `json:"decimals"`
}

// TokenListData is the token list struct used in the default token list
type TokenListData struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	Version   struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Patch int `json:"patch"`
	} `json:"version"`
	LogoURI           string                    `json:"logoURI"`
	Keywords          []string                  `json:"keywords"`
	Tokens            []TokenListToken          `json:"tokens"`
	PreviousTokensMap map[string]TokenListToken `json:"-"`
	NextTokensMap     map[string]TokenListToken `json:"-"`
	Extra             map[string]interface{}    `json:"extra,omitempty"`
}

// InitTokenList initializes the token list
func InitTokenList() TokenListData {
	newTokenList := TokenListData{
		Name:      ``,
		Timestamp: ``,
		Keywords:  []string{},
		Tokens:    []TokenListToken{},
	}
	newTokenList.Version.Major = 0
	newTokenList.Version.Minor = 0
	newTokenList.Version.Patch = 0
	newTokenList.PreviousTokensMap = make(map[string]TokenListToken)
	newTokenList.NextTokensMap = make(map[string]TokenListToken)

	return newTokenList
}

func SetToken(
	address common.Address,
	name string, symbol string, logoURI string,
	chainID uint64, decimals int,
) (TokenListToken, error) {
	token := TokenListToken{}
	if name == `` {
		return token, errors.New(`token name is empty`)
	}
	if symbol == `` {
		return token, errors.New(`token symbol is empty`)
	}
	if decimals == 0 {
		return token, errors.New(`token decimals is 0`)
	}
	if helpers.IsIgnoredToken(chainID, address) {
		return token, errors.New(`token is ignored`)
	}
	if chainID == 0 || helpers.IsChainIDIgnored(chainID) {
		return token, errors.New(`chainID is ignored`)
	}

	token.ChainID = chainID
	token.Decimals = decimals
	token.Address = address.Hex()
	token.Name = name
	token.Symbol = symbol
	token.LogoURI = logoURI
	return token, nil
}

// Assign assigns the original token list to the current token list
func (list TokenListData) Assign(originalTokenList []TokenListToken) TokenListData {
	for _, token := range originalTokenList {
		if helpers.IsChainIDIgnored(token.ChainID) {
			continue
		}
		if (token.Name == `` || token.Symbol == `` || token.Decimals == 0) || helpers.IsIgnoredToken(token.ChainID, common.HexToAddress(token.Address)) {
			continue
		}
		key := getKey(token.ChainID, common.HexToAddress(token.Address))
		list.NextTokensMap[key] = token
	}
	return list
}
