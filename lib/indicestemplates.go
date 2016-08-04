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
	"reflect"
)

func (c *Conn) CheckTemplate(name string) (bool, error) {
	_, err := c.DoCommand("GET", "/_template/"+name, nil, nil)
	if err != nil {
		if IsRecordNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CreateTemplateWithConfig creates an template with config
func (c *Conn) CreateTemplateWithConfig(name string, config interface{}) (BaseResponse, error) {
	var url string
	var retval BaseResponse

	if len(name) > 0 {
		url = fmt.Sprintf("/_template/%s", name)
	} else {
		return retval, fmt.Errorf("You must specify a template name")
	}

	configType := reflect.TypeOf(config)
	if configType.Kind() != reflect.Struct && configType.Kind() != reflect.Map {
		return retval, fmt.Errorf("Config kind was not struct or map")
	}

	requestBody, err := json.Marshal(config)
	if err != nil {
		return retval, err
	}

	body, err := c.DoCommand("PUT", url, nil, requestBody)
	if err != nil {
		return retval, err
	}

	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}

	return retval, err
}
