package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexconrey/go-hs-api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func jsonBodyResponse(key string, value string) ([]byte, error) {
	resp := make(map[string]string)
	resp[key] = value
	return json.Marshal(resp)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Blizzard!")
}

func parseCardsListTemplate(cards []hs_api.Card) string {
	var page = `
	<html>
		<body>
			<table>
				<tr>
					<th>ID</th>
					<th>Card Image</th>
					<th>Name</th>
					<th>Type</th>
					<th>Rarity</th>
					<th>Set</th>
					<th>Class</th>
				</tr>
			{{ range $item := . }}
				<tr>
					<td> {{ $item.ID }} </td>
					<td><img width="66%" height="66%" src="{{ $item.Image }}" /></td>
					<td> {{ $item.Name }} </td>
					<td> {{ $item.Type.Name }} </td>
					<td> {{ $item.Rarity.Name }} </td>
					<td> {{ $item.Set.Name }} </td>
					<td> {{ $item.CardClass.Name }} </td>
				</tr>
			{{ end }}
			</table>
		</body>
	</html>
	`
	output := new(bytes.Buffer)
	tpl := template.New("cards_table")
	tpl.Parse(page)
	tpl.Execute(output, cards)
	return output.String()
}

func cardsListHandler(w http.ResponseWriter, r *http.Request, hs hs_api.HearthstoneAPIClient) {
	q := r.URL.Query()
	requiredParams := []string{
		"classes",
		"manaCost",
		"rarity",
	}

	for _, key := range requiredParams {
		if q.Get(key) == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp, err := jsonBodyResponse(
				"message",
				fmt.Sprintf("Missing parameter: %s", key),
			)
			if err != nil {
				log.Fatalf("JSON marshal error: %s", err)
			}
			w.Write(resp)
			return
		}
	}

	manaMin, err := strconv.Atoi(q.Get("manaCost"))
	if err != nil {
		fmt.Println(err.Error())
	}

	cardListRequest := CardListRequest{
		Classes:  strings.Split(q.Get("classes"), ","),
		ManaCost: manaMin,
		Rarity:   q.Get("rarity"),
	}

	err = cardListRequest.validate()
	if err != nil {
		resp, err := jsonBodyResponse(
			"message",
			fmt.Sprintf("Invalid request: %s", err.Error()),
		)
		if err != nil {
			log.Fatalf("JSON Marshal error: %s", err)
		}

		fmt.Fprint(w, string(resp))
	}

	var cards []hs_api.Card

	manaMaxEnv := os.Getenv("HS_MAX_MANA_COST")
	if manaMaxEnv == "" {
		log.Fatal("HS_MAX_MANA_COST is not set")
	}

	manaMax, err := strconv.Atoi(manaMaxEnv)
	if err != nil {
		fmt.Println(err.Error())
	}

	cards, err = hs.GetCardsWithClassesManaRaritySpec(cardListRequest.Classes, cardListRequest.ManaCost, manaMax, cardListRequest.Rarity)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Sort the ordering based on card ID
	sort.Slice(cards, func(i, j int) bool { return cards[i].ID < cards[j].ID })

	// If limit is defined, set a limit to the returned items
	if q.Get("limit") != "" {
		limit, err := strconv.Atoi(q.Get("limit"))
		if err != nil {
			fmt.Println(err.Error())
		}

		// Limit can only be applied if it is less than or equal to the length of cards
		if len(cards) >= limit {
			cards = cards[:limit]
		}
	}

	fmt.Fprint(w, parseCardsListTemplate(cards))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func main() {
	var (
		clientId     = os.Getenv("HS_API_CLIENT_ID")
		clientSecret = os.Getenv("HS_API_CLIENT_SECRET")
	)

	if clientId == "" {
		log.Fatal("HS_API_CLIENT_ID not set")
	}

	if clientSecret == "" {
		log.Fatal("HS_API_CLIENT_SECRET not set")
	}

	hearthstoneApi, err := hs_api.NewClient(
		clientId,
		clientSecret,
	)
	hearthstoneApi.EndpointURL = "https://us.api.blizzard.com/hearthstone"

	if err != nil {
		log.Fatal(err.Error())
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/cards/list", func(w http.ResponseWriter, r *http.Request) {
		cardsListHandler(w, r, hearthstoneApi)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
