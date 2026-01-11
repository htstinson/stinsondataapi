package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/google/uuid"
	searcher "github.com/htstinson/business_searcher"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

type key struct {
	Value string `json:"value"`
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Test")

	ctx := r.Context()

	apiKey, err := getSecret("Google_Custom_Search")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var k = key{}

	err = json.Unmarshal([]byte(apiKey), &k)
	if err != nil {
		fmt.Println(err.Error())
	}

	var search_definition model.SearchDefinition
	if err := json.NewDecoder(r.Body).Decode(&search_definition); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "invalid search_definition")
		return
	}
	defer r.Body.Close()

	subscriber, err := h.db.GetSubscriber(ctx, search_definition.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	search_definition, err = h.db.GetSearchDefinition(ctx, *subscriber, search_definition.Id, 10, 0)

	search_engine_list, err := h.db.SelectSearchDefinitionEnginesView(ctx, search_definition, 10, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	daterange := searcher.DateRangeConfig{
		Type:      search_definition.SearchType,
		StartDate: search_definition.StartDate.Format("2006-01-02"),
		EndDate:   search_definition.EndDate.Format("2006-01-02"),
	}

	fmt.Println("search engines count:", len(search_engine_list))

	//Load each search
	for _, v := range search_engine_list {

		searches := make([]searcher.SearchQuery, 0)

		var googleSearchConfig = searcher.GoogleSearchConfig{
			DefaultMaxResults: 10,
			DefaultSortByDate: true,
		}

		var config searcher.Config
		var client *searcher.SearchClient
		var output searcher.OutputResult
		var searchquery searcher.SearchQuery
		var count int
		var search_engines = make(map[string]string)

		search_engines[v.SearchEngineName] = v.SearchEngineId

		searchquery = searcher.SearchQuery{
			Name:       v.SearchEngineName,
			Query:      search_definition.Query,
			ExactMatch: search_definition.ExactMatch,
			CSEIDs:     []string{search_engines[v.SearchEngineName]},
			DateRange:  &daterange,
			MaxResults: search_definition.MaxResults,
			SortByDate: search_definition.SortByDate,
		}
		searches = append(searches, searchquery)

		config = searcher.Config{
			GoogleSearch:  googleSearchConfig,
			SearchEngines: search_engines,
			Searches:      searches,
		}

		// Create Google Search CLient
		client, err = searcher.NewSearchClient(k.Value, &config)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Build output structure
		output = searcher.OutputResult{
			Timestamp:     time.Now().Format(time.RFC3339),
			Configuration: client.BuildConfigurationOutput(),
			Searches:      client.ExecuteAllSearches(),
		}

		count = 0
		for _, w := range output.Searches {

			for _, n := range w.Results {
				for _, b := range n.Items {

					subscriberId, err := uuid.Parse(subscriber.Id)
					if err != nil {
						fmt.Println(err.Error())
					}
					search_definition_engine_id, err := uuid.Parse(v.Id)
					if err != nil {
						fmt.Println(err.Error())
					}

					search_time, err := time.Parse(time.RFC3339, output.Timestamp)
					if err != nil {
						fmt.Printf("error parsing time: %v\n", err)
					}

					published, err := extractdate(b.Snippet)
					if err != nil {
						fmt.Println(err.Error())
					}

					calbrate_search_result := model.CalibrateSearchResult{
						Link:                     &b.Link,
						Snippet:                  &b.Snippet,
						Title:                    &b.Title,
						SubscriberID:             subscriberId,
						SearchDefinitionEngineID: &search_definition_engine_id,
						SearchTime:               &search_time,
						Published:                &published,
					}
					_, err = h.db.CreateSearchResult(ctx, *subscriber, calbrate_search_result)
					if err != nil {
						fmt.Println(err.Error())
					}
					count++
					fmt.Println(calbrate_search_result.Published.Format(time.RFC3339))
				}
				fmt.Println("Total Results", count)
				fmt.Println("------------------------------------------------------------------------------------------", v.Id)
				fmt.Println()
			}
			count = 0
		}

	}

}

func extractdate(input string) (time.Time, error) {
	var t time.Time

	// Try to match absolute date pattern like "Sep 24, 2024"
	datePattern := regexp.MustCompile(`([A-Z][a-z]{2})\s+(\d{1,2}),\s+(\d{4})`)
	match := datePattern.FindString(input)

	if match != "" {
		// Parse the absolute date
		parsedDate, err := time.Parse("Jan 2, 2006", match)
		if err != nil {
			return t, fmt.Errorf("error parsing date: %v", err)
		}
		return parsedDate, nil
	}

	// Try to match relative date pattern like "3 days ago", "2 hours ago", etc.
	relativePattern := regexp.MustCompile(`(\d+)\s+(second|minute|hour|day|week|month|year)s?\s+ago`)
	relativeMatch := relativePattern.FindStringSubmatch(input)

	if len(relativeMatch) > 0 {
		// Parse the number
		amount, err := strconv.Atoi(relativeMatch[1])
		if err != nil {
			return t, fmt.Errorf("error parsing relative time amount: %v", err)
		}

		// Get the time unit
		unit := relativeMatch[2]
		now := time.Now()

		// Calculate the date based on the unit
		switch unit {
		case "second":
			return now.Add(-time.Duration(amount) * time.Second), nil
		case "minute":
			return now.Add(-time.Duration(amount) * time.Minute), nil
		case "hour":
			return now.Add(-time.Duration(amount) * time.Hour), nil
		case "day":
			return now.AddDate(0, 0, -amount), nil
		case "week":
			return now.AddDate(0, 0, -amount*7), nil
		case "month":
			return now.AddDate(0, -amount, 0), nil
		case "year":
			return now.AddDate(-amount, 0, 0), nil
		}
	}

	return t, errors.New("no date found")
}

func getSecret(secret_name string) (string, error) {
	secretName := secret_name
	region := "us-west-2"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString

	return secretString, err
}
