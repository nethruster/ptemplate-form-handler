package api

// Request represents the content of the request that ptemplate-form-handler will accept.
// It must be in JSON.
type Request struct {
	Name      string `json:"name"`
	Mail      string `json:"mail"`
	Msg       string `json:"msg"`
	Recaptcha string `json:"g-recaptcha-response"`
}
