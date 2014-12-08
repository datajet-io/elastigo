// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http:www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elastigo

import (
	"encoding/json"
	"fmt"
)

type JsonAliases struct {
	Actions []AliasAction `json:"actions"`
}

type AliasAction struct {
	Add    *JsonAlias `json:"add,omitempty"`
	Remove *JsonAlias `json:"remove,omitempty"`
}

type JsonAlias struct {
	Index string `json:"index"`
	Alias string `json:"alias"`
}

type Aliases map[string]interface{}

// The API allows you to create an index alias through an API.
func (c *Conn) AddAlias(index string, alias string) (BaseResponse, error) {
	var url string
	var retval BaseResponse

	if len(index) > 0 {
		url = "/_aliases"
	} else {
		return retval, fmt.Errorf("You must specify an index to create the alias on")
	}

	jsonAliases := JsonAliases{}
	jsonAliasAdd := AliasAction{}
	jsonAliasAdd.Add = &JsonAlias{
		Index: index,
		Alias: alias,
	}
	jsonAliases.Actions = append(jsonAliases.Actions, jsonAliasAdd)
	requestBody, err := json.Marshal(jsonAliases)

	if err != nil {
		return retval, err
	}

	body, err := c.DoCommand("POST", url, nil, requestBody)
	if err != nil {
		return retval, err
	}

	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}

	return retval, err
}

func (c *Conn) CheckAlias(alias string) (bool, error) {
	_, err := c.DoCommand("GET", "/_alias/"+alias, nil, nil)
	if err != nil {
		if err == RecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

type AliasInfo map[string]interface{}

func (c *Conn) GetAliasIndex(alias string) (string, error) {
	body, err := c.DoCommand("GET", "/_alias/"+alias, nil, nil)
	if err != nil {
		return "", err
	}

	var retval AliasInfo
	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return "", jsonErr
	}

	for key, _ := range retval {
		return key, nil
	}

	return "", fmt.Errorf("No Alias %s found!", alias)
}

func (c *Conn) PutAliases(oldIndex, newIndex, alias string) (BaseResponse, error) {
	var retval BaseResponse

	action := AliasAction{}
	if len(oldIndex) > 0 {
		action.Remove = &JsonAlias{
			Alias: alias,
			Index: oldIndex,
		}
	}

	if len(newIndex) > 0 {
		action.Add = &JsonAlias{
			Alias: alias,
			Index: newIndex,
		}
	}

	aliasOption := JsonAliases{
		Actions: []AliasAction{action},
	}

	body, err := c.DoCommand("POST", "/_aliases", nil, aliasOption)
	if err != nil {
		return retval, err
	}

	if jsonErr := json.Unmarshal(body, &retval); jsonErr != nil {
		return retval, jsonErr
	}

	return retval, nil
}
