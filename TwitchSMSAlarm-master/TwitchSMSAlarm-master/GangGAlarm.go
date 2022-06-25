package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/solapi/solapi-go"
)

func getTwitchAcessToken(clientId string, clientSecret string) string {
	url := "https://id.twitch.tv/oauth2/token?client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=client_credentials"
	reqBody := bytes.NewBufferString("Post")

	resp, err := http.Post(url, "", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(data)[17:47]
}

func GetStreamerLiveB(clientId string, twitchAccessTocken string, streamerId string) bool {
	url := "https://api.twitch.tv/helix/search/channels?query=" + streamerId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("client-id", clientId)
	req.Header.Add("Authorization", "Bearer "+twitchAccessTocken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)

	//fmt.Println(str)
	str = str[strings.Index(str, "\""+streamerId+"\""):]

	if strings.Contains(str, "is_live") {
		pos := strings.Index(str, "is_live") + 9
		str = str[pos : pos+5]

		if strings.Contains(str, "true") {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("[ DU Alarm System ]")

	client := solapi.NewClient()
	fmt.Println("> 솔라피 API 잘 받음.")
	twitchToken := getTwitchAcessToken("qya2dfi9rpv22trnv4w066rjd1k5f4", "2ulhczpqs73lix3nkrfd8qx5h0z90v")
	fmt.Println("> 트위치 토큰 잘 받음.")

	fmt.Println("> 감시 시작...")

	var bLive bool = false
	var bSwitch bool = false
	for {
		bLive = GetStreamerLiveB("qya2dfi9rpv22trnv4w066rjd1k5f4", twitchToken, "xhfl091")

		if bLive {
			if !bSwitch {
				bSwitch = true

				fmt.Println(" > 라디유 방송 시작!")

				message := make(map[string]interface{})
				message["to"] = "01074716322"
				message["from"] = "01074716322"
				message["text"] = "라디유 등장(테스트)"
				message["type"] = "SMS"

				params := make(map[string]interface{})
				params["message"] = message

				result, err := client.Messages.SendSimpleMessage(params)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf(" > solapi result: %+v\n", result)
			} else {
				fmt.Println(" > DU stream ongoing...")
			}
		} else {
			fmt.Println(" > DU stream off")
			if bSwitch {
				bSwitch = false
			}
		}

		time.Sleep(30 * time.Second)
	}
}
