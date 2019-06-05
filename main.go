package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//client, ctx := util.GetGithubClient()
	//events, response, err := client.Issues.ListEvents(ctx, &github.ListOptions{})
	//fmt.Println(events)
	//fmt.Println(response)
	//fmt.Println(err)
	http.HandleFunc("/", webhookHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func webhookHandler(writer http.ResponseWriter, request *http.Request) {
	bodyData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("[webhookHandler] Unable to read body: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}
	body := fmt.Sprint(string(bodyData))
	query := fmt.Sprint(request.URL.RawQuery)
	writer.WriteHeader(200)
	fmt.Fprint(writer, body + "\n\n" + query)
}
