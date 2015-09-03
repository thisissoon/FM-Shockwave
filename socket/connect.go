// Connect to Perceptor for WS Events

package socket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type PerceptorService struct {
	addr   *string
	secret *string
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
			fmt.Println(fmt.Sprintf("WS Connection Failure: %s", err))
			time.Sleep(time.Second)
			continue
		}
	ReadLoop:
		for {
			t, m, e := conn.ReadMessage()
			// On Error close the connection and break the loop
			if e != nil {
				fmt.Println("Error")
				conn.Close()
				break ReadLoop
			}

			// TODO: Push to a message channel
			fmt.Println(t, string(m[:]), e)
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
func NewPerceptorService(addr *string, secret *string) *PerceptorService {
	return &PerceptorService{
		addr:   addr,
		secret: secret,
	}
}
