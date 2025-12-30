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
)

type key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("h Test")

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
	fmt.Println(apiKey)
	fmt.Println("name", k.Name)
	fmt.Println("value", k.Value)

	var googleSearchConfig = searcher.GoogleSearchConfig{
		DefaultMaxResults: 10,
		DefaultSortByDate: true,
		DefaultCSEID:      "1031fbeefdfa24158",
	}

	search_engines := make(map[string]string)
	search_engines["auto_carfax"] = "60d7580b159ae40c3"
	search_engines["auto_cargurus"] = "10d608a7a7c314d6b"
	search_engines["auto_carsforsale"] = "e16b5ec3c749f4e8a"
	search_engines["auto_cylex"] = "d51c57acf13c140d9"
	search_engines["auto_powersports"] = "a2300c57c664540f9"
	search_engines["bbb"] = "4493419a4560045c9"
	search_engines["general_web"] = "1031fbeefdfa24158"
	search_engines["linkedin"] = "30739577b50e043fa"
	search_engines["kvdailyexpress"] = "0160629e137244237"
	search_engines["missouri_times"] = "b29bcb26023d24b15"
	search_engines["schuylercountytimes"] = "01ec1d78a3c654bfc"
	search_engines["kmov"] = "66af520420a4d4f01"
	search_engines["ktvo"] = "20ab9c7bcadd44ea2"
	search_engines["kttn"] = "44b91671818b741fc"
	search_engines["yelp"] = "44a819b1403034ef6"

	search_engines["facebook"] = "134a52f3313b84b07"

	searches := make([]searcher.SearchQuery, 1)
	searchquery := searcher.SearchQuery{
		Name:       "",
		Query:      "",
		ExactMatch: false,
		CSEID:      "",
	}
	searches = append(searches, searchquery)

	var config = searcher.Config{
		GoogleSearch:  googleSearchConfig,
		SearchEngines: search_engines,
		Searches:      searches,
	}

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

	fmt.Println(output)
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
