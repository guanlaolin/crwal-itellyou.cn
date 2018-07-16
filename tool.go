/**
 *   Copyright 2018 Guanlaolin
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func craw() {

}

/**
 * Send HTTP request and unmarshal response body
 *
 * @param method string
 *   http request method
 * @param url string
 *   http request url
 * @param header map[string][]string
 *	 http request header
 * @param body string
 * 	 http request body
 * @param v interface{}
 *   struct of unmarshal json
 *
 * @return []byte
 *	 http response body
 * @return error
 * 	 error message
 */
func unmarshalBody(method string, url string, header map[string][]string,
	body string, v interface{}) error {

	buf, err := fetchBody(method, url, header, body)
	if err != nil {
		LOG_DEBUG("fetch:", err)
		return err
	}

	return json.Unmarshal(buf, v)
}

/**
 * Send HTTP request and read body
 *
 * @param method string
 *   http request method
 * @param url string
 *   http request url
 * @param header map[string][]string
 *	 http request header
 * @param body string
 * 	 http request body
 *
 * @return []byte
 *	 http response body
 * @return error
 * 	 error message
 */
func fetchBody(method string, url string, header map[string][]string,
	body string) ([]byte, error) {

	resp, err := fetch(method, url, header, body)
	if err != nil {
		LOG_DEBUG("fetch:", err)
		return nil, err
	}

	return readRespBody(resp, true)
}

/**
 * Send HTTP request
 *
 * @param method string
 *   http request method
 * @param url string
 *   http request url
 * @param header map[string][]string
 *	 http request header
 * @param body string
 * 	 http request body
 *
 * @return *http.Response
 *	 http response
 * @return error
 * 	 error message
 */
func fetch(method string, url string, header map[string][]string,
	body string) (*http.Response, error) {

	client := http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		LOG_DEBUG("NewRequest:", err)
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		LOG_DEBUG("Do:", err)
		return nil, err
	}

	return resp, nil
}

/**
 * Read HTTP response body
 *
 * @param resp *http.Response
 *   http response
 * @param c bool
 *   true : close resp.Body
 *   false : do not close resp.Body
 *
 * @return []byte
 *	 http response body
 * @return error
 * 	 error message
 */
func readRespBody(resp *http.Response, c bool) ([]byte, error) {
	if c {
		defer resp.Body.Close()
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LOG_DEBUG("ReadAll:", err)
		return nil, err
	}

	return buf, nil
}

func anylazeAll(buf []byte, rule string) [][][]byte {
	reg := regexp.MustCompile(rule)
	res := reg.FindAllSubmatch(buf, -1)

	return res
}

func bytes2float(b []byte) float64 {
	f, err := strconv.ParseFloat(string(b), 8)
	if err != nil {
		log.Fatalln("convert bytes to float error:", err)
	}

	return f
}
