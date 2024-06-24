package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"svclookup/service"
	"svclookup/xutil"
	"time"

	zlog "github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

type RespGeneral struct {
	Message string `json:"message"`
	Success bool `json:"success"`
	// Data any `json:"data"`
	Data []xutil.RespItem `json:"data"`
}

var URLs = []string {
    "http://omnistrategies.net",
    "http://akofisgroup.com",
    "http://akofisengineering.com",

    "https://megafortunelottery.com",
    "https://api.megafortunelottery.com",
    "https://public-api.megafortunelottery.com/swagger/index.html",
    "https://worker.megafortunelottery.com/",
    "https://backend.megafortunelottery.com",
    "https://admin.megafortunelottery.com/",

    "https://backend.mypolicy.market/",
    "https://api.mypolicy.market",
    "https://admin.mypolicy.market",
    "https://temp.mypolicy.market",
    // "https://mypolicy.market",
    // "https://datacollection.mypolicy.market/",
}

func CommentPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Printf("server: %s %s /\n", r.Method, r.RequestURI)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, `{"pong": "Yaay we got %s"}`, id)
}

func DiscoverPage(w http.ResponseWriter, r *http.Request) {
	respData := service.LookupSvc(URLs[:])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var resp = RespGeneral{
		Message: "request completed successfully",
		Success: true,
		Data: respData,
	}

	// encode response to json
	xutil.Encode(w, r, resp)
}

func DiscoverPageTt(w http.ResponseWriter, r *http.Request) {	
	var respData = service.LookupSvc(URLs[:])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var resp = RespGeneral{
		Message: "request completed successfully",
		Success: true,
		Data: respData,
	}
	xutil.PostSlackBetaSignup(r, "test emssage from monitoring")

	// encode response to json
	xutil.JsonResponse(w, resp)
}

func AutoDiscover(currTime time.Time) RespGeneral {
	ft := currTime.Format("2006-01-02 15:04:05")
	// ftz := currTime.Format(time.RFC3339)

	var respData = service.LookupSvcAlt(URLs[:])

	minArr := CleanRespItems(respData) // remove empty result objects
	// send msg only if result is not empty
	if len(minArr) > 0 {
		
		// result json
		dataJson, jErr := json.Marshal(minArr)
		if jErr != nil {
			fmt.Println("ERR JSON", jErr)
		}
		dataStr := "\n"+ft + " | " + string(dataJson)
			
		// send email
		// code, body := SendEmail("s.alfa@akofisgroup.com", "Sumaila Alfa", dataStr)
		// zlog.Debug().Str("Status Code :: %s", fmt.Sprint(code))
		// zlog.Info().Str("Details :: %s", body)
	
		slack := initSlack() // slack init
		sErr := slack.SendSlackMessage( fmt.Sprintf("Lookup Result :: %s", dataStr) )
		if sErr != nil {
			zlog.Error().Msgf("Message could not be sent to %s channel", slack.ChannelID)
		} else {
			zlog.Info().
				// Str("message", message).
				Msgf("Message sent successfully to %s channel", slack.ChannelID)
		}	
			
		// fmt.Println("EXP", minArr)
		AppendToFile(dataStr) //log results to file
	}

	var response = RespGeneral{
		Message: "request completed successfully",
		Success: true,
		Data: minArr,
	}

	// fmt.Println("\nRES >> ", response)
	fmt.Println("RES >> ", response)
	b, bErr := json.Marshal(response)
	if bErr != nil {
        fmt.Printf("AD ERR: %s", bErr)
    }
	zlog.Debug().Str("Result", string(b))

	return response
}

func JsonCheckPage(w http.ResponseWriter, r *http.Request) {
	var respTmp = []xutil.RespItem {
		{Site: "megafortune", StatusCode: 200, Description: "Okay"},
		{Site: "mypolicy.market", StatusCode: 200, Description: "Okay"},
		{Site: "hasura", StatusCode: 500, Description: "Server error"},
		{Site: "public-api", StatusCode: 404, Description: "Not found"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var resp = RespGeneral{
		Message: "request completed successfully",
		Success: true,
		Data: respTmp,
	}
	
	// encode response to json
	xutil.JsonResponse(w, resp)
}

func initSlack() *xutil.Slack {
	config, err := xutil.LoadConfig("..")
    if err != nil {
        zlog.Fatal().Msgf("cannot load config (api):", err)
    }
	slackToken := config.SlackToken
	slackChannel := config.SlackChannel

	return &xutil.Slack{
		SlackClient: slack.New(slackToken),
		ChannelID:   slackChannel,
	}
}

func SendEmail(recipientEmail string, recipientName string, htmlContent string) (int, string) {
	config, err := xutil.LoadConfig("..")
    if err != nil {
        zlog.Fatal().Msgf("cannot load config (api):", err)
    }
	sendgridURL := config.SendGridUrl
	sendgridAPIKey := config.SendGridApiKey
	
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{} {
			{
				"to": []map[string]string{
					{
						"email": recipientEmail,
						"name":  recipientName,
					},
				},
				"subject": "System monitoring",
			},
		},
		"content": []map[string]string{
			{
				"type":  "text/html",
				// "type":  "text/plain",
				"value": htmlContent,
			},
		},
		"from": map[string]string{
			"email": config.FromEmail,
			"name":  config.FromName,
		},
		"reply_to": map[string]string{
			"email": config.FromEmail,
			"name":  config.FromName,
		},
	}
	
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		zlog.Fatal().Msgf("Cannot encode payload:", err)
	}
	
	var req, _ = http.NewRequest("POST", sendgridURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sendgridAPIKey))

	client := service.GetClient()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Email error", err)
	}
	defer resp.Body.Close()
	zlog.Info().
		Str("Code", fmt.Sprint(resp.StatusCode)).
		Str("Description", fmt.Sprint(resp.Status)).
		Msgf("Email payload", string(payloadBytes))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return resp.StatusCode, string(body)
}

// func initSlack(cfg *utils.Config) *xutil.Slack {
// 	return &xutil.Slack{
// 		SlackClient: slack.New(cfg.Slack.SlackToken),
// 		ChannelID:   cfg.Slack.SlackChannelID,
// 	}
// }
