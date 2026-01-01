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

	var payload model.Subscriber_Item
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println(1, err.Error())
		common.RespondError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	defer r.Body.Close()

	subscriber_item, err := h.db.GetSubscriberItem(ctx, payload.Id)
	if err != nil {
		fmt.Println("unable to find subscriber item")
		return
	}

	subscriber, err := h.db.GetSubscriber(ctx, subscriber_item.Subscriber_Id)
	if err != nil {
		fmt.Println("unable to find subscriber")
		return
	}

	search_engines := make(map[string]string)

	search_engine_list, err := h.db.SelectSearchEngines(ctx, *subscriber, 10, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range search_engine_list {
		search_engines[v.Name] = v.SearchEngineId
	}

	var googleSearchConfig = searcher.GoogleSearchConfig{
		DefaultMaxResults: 10,
		DefaultSortByDate: true,
	}

	searches := make([]searcher.SearchQuery, 0)

	//daterange := searcher.DateRangeConfig{
	//	//Type:  "y",  // "d", "w", "m", "y"
	//	//Value: 1000, //int
	//	StartDate: "20250101", //YYYYMMDD format
	//	EndDate:   "20260101", //YYYYMMDD format
	//}

	searchquery2 := searcher.SearchQuery{
		Name:       "Political",
		Query:      `"Jason Soseman"`,
		ExactMatch: false,
		CSEIDs:     []string{search_engines["facebook"]},
		//DateRange:  &daterange,
		MaxResults: 15,
		SortByDate: true,
	}
	searches = append(searches, searchquery2)

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

	for k, v := range output.Searches {
		fmt.Println(k, v)
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
