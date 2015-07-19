package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"sort"
	"sync"
	"time"
)

func print_sorted(languages map[string]int) {
	var values []int
	for k := range languages {
		values = append(values, languages[k])
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	for i := range values {
		for k, v := range languages {
			if v == values[i] {
				fmt.Printf("%s, %d\n", k, v)
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	languages := make(map[string]int)
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: "4314e2c3c994dbcb39927a6758ee063de5825e1e"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	t := time.Now()
	query := fmt.Sprintf("created:>%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	opt := &github.SearchOptions{
		Sort: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	results, _, err := client.Search.Repositories(query, opt)
	if err != nil {
		log.Printf("[Error]: %s\n", err)
	} else {
		for i := range results.Repositories {
			go func(i int) {
				defer wg.Done()
				langs, _, err := client.Repositories.ListLanguages(
					*results.Repositories[i].Owner.Login,
					*results.Repositories[i].Name)
				if err != nil {
					log.Printf("[Error] : %s\n", err)
				} else {
					for key, value := range langs {
						languages[key] += value
					}
				}
			}(i)
			i += 1
		}
	}
	wg.Wait()
	print_sorted(languages)
}
