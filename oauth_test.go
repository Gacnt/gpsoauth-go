package gpsoauth

import (
	"log"
	"os"
	"testing"
)

const androidID string = "9774d56d682e549c"
const oauth_service string = "audience:server:client_id:848232511240-7so421jotr2609rmqakceuu1luuq0ptb.apps.googleusercontent.com"
const app string = "com.nianticlabs.pokemongo"
const client_sig = "321187995bc7cdc2b5fc91b11a96e2baa8602c62"

func TestOAuth(t *testing.T) {
	email := os.Getenv("GPSOAUTH_EMAIL")
	password := os.Getenv("GPSOAUTH_PASSWORD")
	androidID, masterToken, err := Login(email, password, androidID)
	if err != nil {
		log.Println(err)
	}

	body, err := OAuth(email, masterToken, androidID, oauth_service, app, client_sig)
	if err != nil {
		log.Println(err)
	}

	if _, ok := body["Auth"]; !ok {
		t.Errorf("Missing AUTH. Could be an incorrect Email or Password, or 2 step authentication failure. (This package does not support 2 step auth)")
	}

}
