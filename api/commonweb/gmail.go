package commonweb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getClient(oauthConfig *oauth2.Config, region string) *http.Client {
	tok, err := tokenFromSecret("gmail-token", region)
	if err != nil {
		tok = getTokenFromWeb(oauthConfig)
		saveTokenToSecret("gmail-token", tok)
	}

	tokenSource := oauthConfig.TokenSource(context.Background(), tok)

	// Save refreshed token back to Secrets Manager
	newTok, err := tokenSource.Token()
	if err == nil && newTok.AccessToken != tok.AccessToken {
		saveTokenToSecret("gmail-token", newTok)
	}

	return oauth2.NewClient(context.Background(), tokenSource)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	config.RedirectURL = "http://localhost:8080"

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Opening browser for authorization...")

	codeCh := make(chan string)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintln(w, "Authorization successful! You can close this tab.")
		codeCh <- code
	})

	go http.ListenAndServe(":8080", nil)

	// Open browser automatically
	fmt.Printf("Visit this URL if browser doesn't open:\n%v\n", authURL)

	tok, err := config.Exchange(context.TODO(), <-codeCh)
	if err != nil {
		log.Fatal("Unable to exchange token:", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	return tok, json.NewDecoder(f).Decode(tok)
}

func tokenFromSecret(secretName string, region string) (*oauth2.Token, error) {
	data, err := GetSecretString(secretName, region)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	return tok, json.Unmarshal(data, tok)
}

func saveTokenToSecret(secretName string, token *oauth2.Token) {
	data, err := json.Marshal(token)
	if err != nil {
		log.Fatal("Could not marshal token:", err)
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	client := secretsmanager.NewFromConfig(cfg)

	secretStr := string(data)
	_, err = client.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId:     &secretName,
		SecretString: &secretStr,
	})
	if err != nil {
		log.Fatal("Could not save token to Secrets Manager:", err)
	}
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("Unable to save token:", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sendEmail(service *gmail.Service, to, subject, body string) error {
	raw := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
	encoded := base64.URLEncoding.EncodeToString([]byte(raw))
	msg := &gmail.Message{Raw: encoded}
	_, err := service.Users.Messages.Send("me", msg).Do()
	return err
}

func SendMail(to string, subject string, body string, region string) {
	fmt.Println("commonweb SendMail")
	ctx := context.Background()

	// Download credentials.json from Google Cloud Console
	b, err := GetSecretString("gmail-credentials", region)
	if err != nil {
		fmt.Println("1 ", err.Error())
		log.Fatal("Cannot read gmail credentials:", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		fmt.Println("2 ", err.Error())
		log.Fatal("Cannot parse credentials:", err)
	}

	client := getClient(config, region)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("3 ", err.Error())
		log.Fatal("Cannot create Gmail service:", err)
	}

	err = sendEmail(srv, to, subject, body)
	if err != nil {
		fmt.Println("4 ", err.Error())
		log.Fatal("Send failed:", err)
	}
	fmt.Println("Email sent!")
}
