package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	searcher "github.com/htstinson/business_searcher"
	common "github.com/htstinson/stinsondataapi/api/commonweb"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

type key struct {
	Value string `json:"value"`
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("SearchDefinition", search_definition)

	search_engines := make(map[string]string)

	search_engine_list, err := h.db.SelectSearchDefinitionEnginesView(ctx, search_definition, 10, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("start date", search_definition.StartDate)
	fmt.Println("end date", search_definition.EndDate)
	fmt.Println("query", search_definition.Query)
	fmt.Println("search type", search_definition.SearchType)

	searches := make([]searcher.SearchQuery, 0)

	daterange := searcher.DateRangeConfig{
		Type:      search_definition.SearchType,
		StartDate: search_definition.StartDate.Format("2006-01-02"),
		EndDate:   search_definition.EndDate.Format("2006-01-02"),
	}

	var googleSearchConfig = searcher.GoogleSearchConfig{
		DefaultMaxResults: 10,
		DefaultSortByDate: true,
	}

	//Load each search
	for _, v := range search_engine_list {
		search_engines[v.SearchEngineName] = v.SearchEngineId
		fmt.Println("search engine name", v.SearchEngineName)
		fmt.Println("search engine id", v.SearchEngineId)

		searchquery := searcher.SearchQuery{
			Name:       v.SearchEngineName,
			Query:      search_definition.Query,
			ExactMatch: search_definition.ExactMatch,
			CSEIDs:     []string{search_engines[v.SearchEngineName]},
			DateRange:  &daterange,
			MaxResults: search_definition.MaxResults,
			SortByDate: search_definition.SortByDate,
		}
		searches = append(searches, searchquery)

	}

	var config = searcher.Config{
		GoogleSearch:  googleSearchConfig,
		SearchEngines: search_engines,
		Searches:      searches,
	}

	// Create Google Search CLient
	client, err := searcher.NewSearchClient(k.Value, &config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Build output structure
	output := searcher.OutputResult{
		Timestamp:     time.Now().Format(time.RFC3339),
		Configuration: client.BuildConfigurationOutput(),
		Searches:      client.ExecuteAllSearches(),
	}

	fmt.Println("Timestamp", output.Timestamp)
	fmt.Println("Total Searches", output.Configuration.TotalSearches)

	for k, v := range output.Searches {
		for m, n := range v.Results {
			for a, b := range n.Items {
				fmt.Println("---------------------------")
				fmt.Println(k, m, a, b.Link)
				fmt.Println(k, m, a, b.Position)
				fmt.Println(k, m, a, b.Snippet)
				fmt.Println(k, m, a, b.Title)
			}
		}
	}

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
