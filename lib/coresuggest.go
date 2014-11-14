// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
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

// Suggest performs a very basic suggest search on an index via the request URI API.
//
// params:
//   @index:  the elasticsearch index
//   @_type:  optional ("" if not used) search specific type in this index
//   @args:   a map of URL parameters. Allows all the URI-request parameters allowed by ElasticSearch.
//   @query:  this can be one of 3 types:
//              1)  string value that is valid elasticsearch
//              2)  io.Reader that can be set in body (also valid elasticsearch string syntax..)
//              3)  other type marshalable to json (also valid elasticsearch json)
//
//   out, err := Suggest(true, "github", map[string]interface{} {"from" : 10}, qryType)
//
// http://www.elasticsearch.org/guide/reference/api/search/uri-request.html
func (c *Conn) Suggest(index string, _type string, args map[string]interface{}, query interface{}) (SuggestResult, error) {
	var uriVal string
	var retval SuggestResult
	if len(_type) > 0 && _type != "*" {
		uriVal = fmt.Sprintf("/%s/%s/_suggest", index, _type)
	} else {
		uriVal = fmt.Sprintf("/%s/_suggest", index)
	}
	body, err := c.DoCommand("POST", uriVal, args, query)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal([]byte(body), &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	retval.RawJSON = body
	return retval, err
}

type SuggestResult struct {
	RawJSON     []byte
	ShardStatus Status `json:"_shards"`
}

type SuggesterQuery struct {
	Text       string           `json:"text"`
	Completion *CompletionQuery `json:"completion,omitempty"`
}

type CompletionQuery struct {
	Size  int           `json:"size,omitempty"`
	Field string        `json:"field,omitempty"`
	Fuzzy *FuzzyOptions `json:"fuzzy,omitempty"`
}

type FuzzyOptions struct {
	Fuzziness string `json:"fuzziness,omitempty"`
}
