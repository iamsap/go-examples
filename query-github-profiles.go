// summary: A simple library showing how to pull profiles from github using channels
// author: robbymillsap@gmail.com

package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strings"
	"time"
)

var userch = make(chan string)

func main() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Tell me who you'd like to query (type \"exit\" to quit)")
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
		t0 := time.Now()

		session, err := mgo.Dial("localhost:27017")
		if err != nil {
			fmt.Println("Cannot connect", err)
			panic(err)
		}
		defer session.Close()

		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)

		c := session.DB("local").C("profiles")

		result := Profile{}
		err = c.Find(bson.M{"login": username}).One(&result)

		fmt.Println("Queried: ", result)

		if err != nil || result.Login == "" {
			fmt.Println("Cannot find user", err)
			continue
		}

		fmt.Println("Found:", result.Name)

		t1 := time.Now()
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
