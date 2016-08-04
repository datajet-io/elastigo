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
	"sort"
)

// TermVector returns information and statistics on terms in the fields of a particular document.
func (c *Conn) TermVector(index, _type, id string, termVecotrReq TermVectorRequest, args map[string]interface{}) (TermVecotrResponse, error) {
	var url string
	var retval TermVecotrResponse
	if len(index) <= 0 {
		return retval, fmt.Errorf("index name is empty")
	}
	if len(_type) <= 0 {
		return retval, fmt.Errorf("index type is empty")
	}
	if len(id) <= 0 {
		return retval, fmt.Errorf("doc id is empty")
	}

	url = fmt.Sprintf("/%s/%s/%s/_termvector", index, _type, id)

	body, err := c.DoCommand("GET", url, args, termVecotrReq)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	return retval, err
}

// TermVectorRequest contains request for _termvector
type TermVectorRequest struct {
	Fields          []string `json:"fields"`
	Offsets         bool     `json:"offsets"`
	Payloads        bool     `json:"payloads"`
	Positions       bool     `json:"positions"`
	TermStatistics  bool     `json:"term_statistics"`
	FieldStatistics bool     `json:"field_statistics"`
}

// TermVecotrResponse contains response from _termvector
type TermVecotrResponse struct {
	Index       string             `json:"_index"`
	Type        string             `json:"_type"`
	ID          string             `json:"_id"`
	Found       bool               `json:"found"`
	TermVectors map[string]TermMap `json:"term_vectors"`
}

// TermMap wraps value of "_terms" field
type TermMap struct {
	Terms map[string]TermValue
}

// TermValue defines structure for specific term
type TermValue struct {
	TermFreq int `json:"term_freq"`
}

// SortByFreq is used for sorting term by frequency
type SortByFreq struct {
	Term string
	Freq int
}

// SortByFreqList is the list of TermFreq
type SortByFreqList []SortByFreq

func (t SortByFreqList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t SortByFreqList) Len() int           { return len(t) }
func (t SortByFreqList) Less(i, j int) bool { return t[i].Freq > t[j].Freq }

// SortMapByTermFreq sort by term's frequency
func SortMapByTermFreq(m map[string]TermValue) SortByFreqList {
	p := make(SortByFreqList, len(m))
	i := 0
	for k, v := range m {
		p[i] = SortByFreq{k, v.TermFreq}
		i++
	}

	sort.Sort(p)
	return p
}
