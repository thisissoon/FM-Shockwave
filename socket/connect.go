// Connect to Perceptor for WS Events

package socket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type PerceptorService struct {
	channel chan []byte
	addr    *string
	secret  *string
}

// Connect to Perceptor WS Service and Consume the Messages from the Service
func (p *PerceptorService) Run() {
	// Create Dialer
	d := &websocket.Dialer{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	// Attempt to connect to WS Service
	for {
		conn, _, err := d.Dial(fmt.Sprintf("ws://%s", *p.addr), p.headers())
		if err != nil {
			log.Errorf("WS Connection Failure: %s", err)
			time.Sleep(time.Second)
			continue
		}
	ReadLoop:
		for {
			t, m, e := conn.ReadMessage()
			// On Error close the connection and break the loop
			if e != nil {
				log.Errorf("WS Connection Lost: %s", e)
				conn.Close()
				break ReadLoop
			}

			if t == websocket.TextMessage {
				// Put the message on the event channel
				log.Debugf("Recived Message: %s", m)
				p.channel <- m
			}
		}
	}
}

// Generates HTTP Headers for connecting the Perceptor WS Service
func (p *PerceptorService) headers() http.Header {
	data := []byte("") // No request data is sent

	// Generate HMAC Signature
	mac := hmac.New(sha256.New, []byte(*p.secret))
	mac.Write(data)
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Create Headers
	headers := http.Header{}
	headers.Add("Signature", fmt.Sprintf("%s:%s", "shockwave", sig))

	return headers
}

// Creates a new PerceptorService
func NewPerceptorService(addr *string, secret *string, channel chan []byte) *PerceptorService {
	return &PerceptorService{
		channel: channel,
		addr:    addr,
		secret:  secret,
	}
}
