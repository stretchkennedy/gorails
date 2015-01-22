package main

import(
	"github.com/mattetti/goRailsYourself/crypto"
	"net/http"
	"net/url"
	"log"
	"fmt"
)

type AuthHandler struct {
	handler http.Handler
	secret string
	UserId int
}

func NewAuthHandler(secret string, handler http.Handler) *AuthHandler {
	return &AuthHandler{handler: handler, secret: secret}
}

func decrypt(railsSecret string, ciphertext string) (interface{}, error) {
	encryptedCookieSalt := []byte("encrypted cookie")
	encryptedSignedCookieSalt := []byte("signed encrypted cookie")

	kg := crypto.KeyGenerator{Secret: railsSecret}
	secret := kg.CacheGenerate(encryptedCookieSalt, 32)
	signSecret := kg.CacheGenerate(encryptedSignedCookieSalt, 64)
	e := crypto.MessageEncryptor{Key: secret, SignKey: signSecret}
	//e.Serializer = crypto.NullMsgSerializer{}

	var data interface{}
	err := e.DecryptAndVerify(ciphertext, &data)
	return data, err
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	showerr := func(err interface{}, code int) {
		log.Printf("AuthHandler: %d: %v", code, err)
		http.Error(w, fmt.Sprintf("%v", err), code)
	}


	// extract cookie
	cookie, err := r.Cookie("_BeCollective_session")
	if err != nil {
		showerr(err, 422)
		return
	}

	// extract ciphertext
	ciphertext, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		showerr(err, 422)
		return
	}

	// decrypt
	data, err := decrypt(h.secret, ciphertext)
	if err != nil {
		showerr(err, 422)
		return
	}

	// extract user id
	// {... "warden.user.user.key": [[1234], ...], ...}
	hash := data.(map[string]interface{})
	warden := hash["warden.user.user.key"]
	switch warden.(type) {
	case nil:
		showerr("Verification failed: Logged out", 401)
		return
	}
	num := warden.([]interface{})[0].([]interface{})[0].(float64)
	h.UserId = int(num)

	// serve
	h.handler.ServeHTTP(w, r)
}
