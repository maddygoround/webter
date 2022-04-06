package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"

	"github.com/maddygoround/webter/pkg/sd"
)

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}

var (
	key   string
	nonce string
)

func encode(this js.Value, i []js.Value) interface{} {
	encoded, err := func() (string, string) {
		answerSd := sd.SessionDescription{
			Sdp:   i[0].String(),
			Key:   key,
			Nonce: nonce,
		}
		// Encrypt with the shared keys from the offer
		if key != "" {
			if err := answerSd.Encrypt(); err != nil {
				return "", err.Error()
			}
		}

		// Don't upload the keys, the host has them
		answerSd.Key = ""
		answerSd.Nonce = ""
		return sd.Encode(answerSd), ""
	}()
	i[1].Invoke(encoded, err)
	return nil
}

func decode(this js.Value, i []js.Value) interface{} {
	sdp, tkbsl, err := func() (string, string, string) {
		offer, err := sd.Decode(i[0].String())
		if err != nil {
			return "", "", err.Error()
		}
		if offer.Key != "" {
			key = offer.Key
			nonce = offer.Nonce
			if err := offer.Decrypt(); err != nil {
				return "", "", err.Error()
			}
		}
		return offer.Sdp, offer.TenKbSiteLoc, ""
	}()
	i[1].Invoke(sdp, tkbsl, err)
	return nil
}

func read10kbfile(this js.Value, i []js.Value) interface{} {
	status, body, err := func() (int, string, string) {
		resp, err := http.Get(fmt.Sprintf("%s%s", i[0].String(), i[1].String()))
		if err != nil {
			return 0, "", err.Error()
		}
		body, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusOK {
			return resp.StatusCode, "", fmt.Errorf(
				"Resp %d 10kb.site error: %s", resp.StatusCode, string(body)).Error()
		}
		if err != nil {
			return resp.StatusCode, "", err.Error()
		}
		return resp.StatusCode, string(body), ""
	}()
	i[1].Invoke(status, body, err)
	return nil
}

func registerCallbacks() {
	js.Global().Set("encode", js.FuncOf(encode))
	js.Global().Set("decode", js.FuncOf(decode))
	js.Global().Set("read10kbfile", js.FuncOf(read10kbfile))
}
