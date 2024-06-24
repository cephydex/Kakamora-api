package xutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type RespItem struct {
	Site string `json:"site"`
	StatusCode int `json:"status_code"`
	Description string `json:"description"`
}


var timeout = time.Duration(4 * time.Second)

func ReqTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func ProcessReq(reqURL string, client http.Client, logit ...bool) RespItem {
	var resp, err = client.Get(reqURL)
	// if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		// retry
		fmt.Println("RETRY", err)
		resp, err = client.Get(reqURL)
	}
	if err != nil {
		log.Fatal(err)
	}
	
	defer resp.Body.Close()

	if len(logit) > 0 {
		fmt.Printf("\nResp: %s :: %d -- %s", reqURL, resp.StatusCode, resp.Status)
	}
	
	return RespItem{
		Site: reqURL,
		StatusCode: resp.StatusCode,
		Description: resp.Status,
	}
}

func ProcessReqMV(reqURL string, client http.Client) (string, int, string) {
	resp, err := client.Get(reqURL)
	if err != nil {
		// handle error
		fmt.Println("Error", err)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Resp: %s :: %d -- %s \n", reqURL, resp.StatusCode, resp.Status)
	// if resp.StatusCode != http.StatusOK {
	// }
	return reqURL, resp.StatusCode, resp.Status
}

func JsonResponse(w http.ResponseWriter, resp any) {
	if err := json.NewEncoder(w).Encode(resp) ; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Errorf("error json %w", err)
	}
}

func Encode[T any](w http.ResponseWriter, r *http.Request, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(v) ; err != nil {
		return fmt.Errorf("error json %w", err)
	}
	return nil
}

// func decode[T any](r *http.Request) (T, error) {
// 	var v T

// 	if err := json.NewDecoder(r.Body).Decode(v) ; err != nil {
// 		return v, fmt.Errorf("error json %w", err)
// 	}
// 	return v, nil
// }
