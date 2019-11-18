package server

// Package server will manage all the HTTP request made to web-msg-handler.

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/web-msg-handler/api"
	"github.com/Miguel-Dorta/web-msg-handler/pkg"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
)

const statusUnknownError = 502

var (
	Log *logolang.Logger
	s *sender.Mail
)

// Run will start a HTTP server in the port provided using the config file path provided.
// It ends when a termination or interrupt signal is received.
// It can end the program execution prematurely.
func Run(configFile, port string) {
	// Load config
	var err error
	s, err = config.LoadConfig(configFile)
	if err != nil {
		Log.Criticalf("error loading config file from path \"%s\": %s", configFile, err)
		os.Exit(1)
	}

	http.HandleFunc("/", handle)
	srv := http.Server{Addr: ":" + port}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit // Block until quit signal is received

		Log.Info("Shutting down")

		if err := srv.Shutdown(context.Background()); err != nil {
			Log.Criticalf("error while shutting down: %s", err)
			os.Exit(1)
		}
	}()

	Log.Infof("Listening to port %s", port)
	Log.Info("Press CTRL + C to exit")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		Log.Criticalf("Unexpected error which closed the server: %s", err)
		os.Exit(1)
	}
}

// handle is the function executed for each HTTP request received by web-msg-handler.
//
// It will:
//
// - Check if the HTTP method used is POST
//
// - Check if the Content-Type header is the MIME JSON.
//
// - Check if the request body is valid.
//
// - Check if the email provided is valid.
//
// - Check if the request have passed the ReCaptcha verification.
//
// - Send the message
func handle(w http.ResponseWriter, r *http.Request) {
	// Request ID for logging purposes
	Log.Debug("Request received")

	if method := r.Method; method != http.MethodPost {
		Log.Errorf("Invalid method: %s", method)
		statusWriter(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("method %s not supported", method))
		return
	}

	if contentType := r.Header.Get(pkg.MimeContentType); contentType != pkg.MimeJSON {
		Log.Errorf("Invalid content type: %s", contentType)
		statusWriter(w, http.StatusBadRequest, false, fmt.Sprintf("content-type %s not supported", contentType))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Errorf("Error while reading body: %s", err)
		statusWriter(w, statusUnknownError, false, fmt.Sprintf("unknown error while reading request body: %s", err.Error()))
		return
	}

	var r2 api.Request
	if err = json.Unmarshal(body, &r2); err != nil {
		Log.Errorf("Malformed JSON: %s", err)
		statusWriter(w, http.StatusBadRequest, false, "malformed JSON")
		return
	}

	if !sanitation.IsValidMail(r2.Mail) {
		Log.Error("Invalid email")
		statusWriter(w, http.StatusBadRequest, false, "invalid email")
		return
	}

	if err = s.CheckRecaptcha(r2.Recaptcha); err != nil {
		Log.Errorf("Recaptcha verification failed: %s", err)
		statusWriter(w, http.StatusBadRequest, false, "recaptcha verification failed")
		return
	}

	if err = s.Send(sanitation.SanitizeName(r2.Name), r2.Mail, sanitation.SanitizeMsg(r2.Msg)); err != nil {
		Log.Errorf("Sender failed: %s", err)
		statusWriter(w, http.StatusServiceUnavailable, false, "error sending message")
		return
	}

	statusWriter(w, http.StatusOK, true, "")
	Log.Debug("Success")
}

// statusWriter will write a response to the http.ResponseWriter provided.
// That response will be sent with the status code provided,
// and its body will consists in a JSON represented by api.Response with the success status and error provided.
func statusWriter(w http.ResponseWriter, statusCode int, success bool, msg string) {
	w.Header().Set(pkg.MimeContentType, pkg.MimeJSON)
	w.WriteHeader(statusCode)

	data, _ := json.Marshal(api.Response{
		Success: success,
		Err:     msg,
	})

	if _, err := w.Write(data); err != nil {
		Log.Errorf("error writing response: %s", err)
	}
}
