package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"sync"
	"time"
)

type jsonUser struct {
	Name string `json:"login"`
	Blog string `json:"blog"`
}

type jsonAuth struct {
	client_id     string `json:"login"`
	client_secret string `json:"blog"`
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	languages := map[string]int64{}
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: "4314e2c3c994dbcb39927a6758ee063de5825e1e"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	opt := &github.SearchOptions{
		Sort: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	t := time.Now()
	date := fmt.Sprintf("created:>%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	results, _, err := client.Search.Repositories(date, opt)
	if err != nil {
		log.Printf("[Error]: %s\n", err)
	} else {
		for i := range results.Repositories {
			go func(i int) {
				defer wg.Done()
				lgs, _, err := client.Repositories.ListLanguages(
					*results.Repositories[i].Owner.Login,
					*results.Repositories[i].Name)
				if err != nil {
					log.Printf("[Error] : %s\n", err)
				} else {
					for key, value := range lgs {
						languages[key] += int64(value)
					}
				}
			}(i)
			i += 1
		}
	}
	wg.Wait()
	for key, value := range languages {
		fmt.Printf("%s: %d lines\n", key, value)
	}
}
