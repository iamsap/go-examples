// summary: A simple library showing how to pull profiles from github using channels
// author: robbymillsap@gmail.com

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var userch = make(chan string)

func main() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Tell me who you'd like to lookup (type \"exit\" to quit)")
	reader := bufio.NewReader(os.Stdin)

	go printUserInfo(userch)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Ok I'm done")
			break
		}

		userch <- text
	}

}

func printUserInfo(ch <-chan string) {

	for username := range ch {

		client := http.Client{}
		t0 := time.Now()

		var profile Profile
		resp, err := client.Get("https://api.github.com/users/" + username)

		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("An error occurred loading", username)
			return
		}

		err = json.Unmarshal(data, &profile)

		t1 := time.Now()

		if profile.Name != "" {
			fmt.Printf("Looks like we found %s from %s\n", profile.Name, profile.Location)
		} else {
			fmt.Println("Sorry, couldn't find", username, "but hey we tried!")
		}

		fmt.Printf("Duration %v\n", t1.Sub(t0))

	}
}

type Profile struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Location  string `json:"location"`
}
