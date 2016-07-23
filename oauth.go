package gpsoauth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	USER_AGENT string = "Dalvik/2.1.0 (Linux; U; Android 5.1.1; Andromax I56D2G Build/LMY47V"
	AUTH_URL   string = "https://android.clients.google.com/auth"
)

type GoogleOAuthOptions struct {
	URL         string
	Method      string
	ContentType string
	Data        *strings.Reader
	Token       string
}

func parseBody(body string) map[string]string {
	holder := make(map[string]string)
	for _, line := range strings.Split(body, "\n") {
		if strings.Contains(line, "=") {
			temp := strings.Split(line, "=")
			holder[temp[0]] = temp[1]
		}
	}
	return holder
}

func request(options *GoogleOAuthOptions) (map[string]string, error) {
	req, err := http.NewRequest(options.Method, options.URL, options.Data)
	if err != nil {
		log.Printf("%v \nFailed to create new request url with Method: %v, URL: %v", err, options.Method, options.URL)
		return nil, err
	}

	if options.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("GoogleLogin auth=%s", options.Token))
	}

	req.Header.Add("User-Agent", USER_AGENT)

	if options.ContentType == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", options.ContentType)
	}

	newClient := http.Client{Timeout: 10 * time.Second}
	resp, err := newClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	jsonBody := parseBody(string(body))

	return jsonBody, nil

}

func OAuth(email, master_token, android_id, service, app, client_sig string) (map[string]string, error) {
	form := url.Values{}
	form.Add("accountType", "HOSTED_OR_GOOGLE")
	form.Add("Email", email)
	form.Add("EncryptedPasswd", master_token)
	form.Add("has_permission", "1")
	form.Add("service", service)
	form.Add("source", "android")
	form.Add("androidId", android_id)
	form.Add("app", app)
	form.Add("client_sig", client_sig)
	form.Add("device_country", "us")
	form.Add("operatorCountry", "us")
	form.Add("lang", "en")
	form.Add("sdk_version", "17")

	formBody := strings.NewReader(form.Encode())

	body, err := request(&GoogleOAuthOptions{
		URL:         AUTH_URL,
		Method:      "POST",
		ContentType: "application/x-www-form-urlencoded",
		Data:        formBody,
	})

	return body, err
}

func Login(email, password, android_id string) (string, string, error) {
	form := url.Values{}
	form.Add("accountType", "HOSTED_OR_GOOGLE")
	form.Add("Email", strings.TrimSpace(email))
	form.Add("has_permission", "1")
	form.Add("add_account", "1")
	form.Add("Passwd", password)
	form.Add("service", "ac2dm")
	form.Add("source", "android")
	form.Add("androidId", android_id)
	form.Add("device_country", "us")
	form.Add("operatorCountry", "us")
	form.Add("lang", "en")
	form.Add("sdk_version", "17")

	formBody := strings.NewReader(form.Encode())

	body, err := request(&GoogleOAuthOptions{
		Method:      "POST",
		URL:         AUTH_URL,
		ContentType: "application/x-www-form-urlencoded",
		Data:        formBody,
	})

	return android_id, body["Token"], err
}
