package elastigo

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Setting struct {
	NumberOfShards   int             `json:"number_of_shards"`
	NumberOfReplicas int             `json:"number_of_replicas"`
	Analysis         *AnalysisOption `json:"analysis,omitempty"`
}

type AnalysisOption struct {
	Analyzer  map[string]AnalyzerOption  `json:"analyzer"`
	Tokenizer map[string]TokenizerOption `json:"tokenizer"`
	Filter    map[string]FilterOption    `json:"filter"`
}

type AnalyzerOption struct {
	Type       string   `json:"type"`
	Tokenizer  string   `json:"tokenizer"`
	Filter     []string `json:"filter,omitempty"`
	CharFilter []string `json:"char_filter,omitempty"`
}

type TokenizerOption struct {
	Type           string   `json:"type"`
	MaxTokenLength int      `json:"max_token_length,omitempty"`
	MinGram        int      `json:"min_gram,omitempty"`
	MaxGram        int      `json:"max_gram,omitempty"`
	TokenChars     []string `json:"token_chars,omitempty"`
}

type FilterOption struct {
	Type       string   `json:"type"`
	Name       string   `json:"name,omitempty"`
	Stopwords  []string `json:"stopwords,omitempty"`
	Min        int      `json:"min,omitempty"`
	Max        int      `json:"max,omitempty"`
	MinGram    int      `json:"min_gram,omitempty"`
	MaxGram    int      `json:"max_gram,omitempty"`
	TokenChars []string `json:"token_chars,omitempty"`
}

func (c *Conn) PutSettings(index string, settings interface{}) (BaseResponse, error) {

	var url string
	var retval BaseResponse

	settingsType := reflect.TypeOf(settings)
	if settingsType.Kind() != reflect.Struct {
		return retval, fmt.Errorf("Settings kind was not struct")
	}

	if len(index) > 0 {
		url = fmt.Sprintf("/%s/_settings", index)
	} else {
		url = fmt.Sprintf("/_settings")
	}

	requestBody, err := json.Marshal(settings)

	if err != nil {
		return retval, err
	}

	body, errDo := c.DoCommand("PUT", url, nil, requestBody)
	if errDo != nil {
		return retval, errDo
	}

	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}

	return retval, err
}
