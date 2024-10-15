package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type gameDatabaseReturn []struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Thumb string `json:"thumb"`
	Url   string `json:"url"`
}

type gameCoverDatabaseReturn []struct {
	GameID int    `json:"game"`
	Url    string `json:"url"`
}

type gameWesbiteDatabaseReturn []struct {
	Url      string `json:"url"`
	Category int    `json:"category"`
}

type steamPlayerCountReturn struct {
	Response struct {
		PlayerCount int `json:"player_count"`
	} `json:"response"`
}

type twitchGameIDReturn struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type twitchStreamsReturn struct {
	Data []struct {
		ViewerCount int `json:"viewer_count"`
	} `json:"data"`
	CursorData struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type playerAnalysisToReturn struct {
	Twitch *twitchAnalysis `json:"twitch,omitempty"`
	Steam  *steamAnalysis  `json:"steam,omitempty"`
}

type twitchAnalysis struct {
	CurrentViewers          int     `json:"current_viewers"`
	CurrentStreams          int     `json:"current_streams"`
	AverageViewersPerStream float32 `json:"average_viewers"`
}

type steamAnalysis struct {
	CurrentPlayers           int     `json:"current_players"`
	PercentageOfTotalPlayers float32 `json:"percentage_of_total"`
}

type queryParameter struct {
	key   string
	value string
}

const gamesDatabaseAPIGames string = "https://api.igdb.com/v4/games"
const gamesDatabaseAPICovers string = "https://api.igdb.com/v4/covers"
const gameDatabaseAPIWebsites string = "https://api.igdb.com/v4/websites"
const ourClientID string = "n092twgjuyvozrvzrtu3r7fwd0etmh"
const ourAuthorization string = "Bearer irmkybj1r349mha5oqv1tn1xvtdm8u"
const searchStatement string = "fields name, url; where category = 0 & (status = null | status = 4); search "
const getGameCoversStatement string = "fields game, url; where "
const gameEqualsStatement string = "game = "
const limitStatement string = "limit 10;"
const getGameWebsitesStatement string = "fields url, category; where game = "
const websiteTypeStatement string = " & (category = 13 | category = 6)"
const steamWebAPI string = "https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid="
const twitchAPIgetGames string = "https://api.twitch.tv/helix/games"
const twitchAPIgetStreams string = "https://api.twitch.tv/helix/streams"
const gameDatabaseID string = "igdb_id"
const twitchDatabaseID string = "game_id"

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(CORSMiddleware())
	router.GET("/getGames/:search", getGames)
	router.GET("/getPlayerAnalysis/:id", getGamePlayerStatistics)

	fmt.Println("Backend 2.0 online!")
	router.Run("0.0.0.0:8080")
}

func getGames(c *gin.Context) {
	var searchTerm string = c.Param("search")
	var bodyToSend = searchStatement + "\"" + searchTerm + "\"; " + limitStatement
	var jsonifiedFetch gameDatabaseReturn

	json.Unmarshal(fetchFromTwitchOrGameDatabase("POST", gamesDatabaseAPIGames, bodyToSend, nil), &jsonifiedFetch)

	var gameThumbs = make(map[int]string)
	var coverhttpRequestBody = getGameCoversStatement
	var firstGameFilter = true

	for _, result := range jsonifiedFetch {
		gameThumbs[result.ID] = ""

		if firstGameFilter {
			firstGameFilter = false
		} else {
			coverhttpRequestBody += " | "
		}

		coverhttpRequestBody += gameEqualsStatement + " " + strconv.Itoa(result.ID)
	}

	coverhttpRequestBody += ";"

	var jsonifiedCoverFetch gameCoverDatabaseReturn

	json.Unmarshal(fetchFromTwitchOrGameDatabase("POST", gamesDatabaseAPICovers, coverhttpRequestBody, nil), &jsonifiedCoverFetch)

	for _, cover := range jsonifiedCoverFetch {
		gameThumbs[cover.GameID] = cover.Url
	}

	for i := 0; i < len(jsonifiedFetch); i++ {
		jsonifiedFetch[i].Thumb = gameThumbs[jsonifiedFetch[i].ID]
	}

	c.IndentedJSON(http.StatusOK, jsonifiedFetch)
}

func getGamePlayerStatistics(c *gin.Context) {
	gameID := c.Param("id")

	if _, err := strconv.Atoi(c.Param("id")); err != nil {
		return
	}

	var gameWebsites gameWesbiteDatabaseReturn
	json.Unmarshal(fetchFromTwitchOrGameDatabase("POST", gameDatabaseAPIWebsites, getGameWebsitesStatement+gameID+websiteTypeStatement+";", nil), &gameWebsites)

	var steamURL string
	var twitchURL string

	for _, website := range gameWebsites {
		switch website.Category {
		case 6:
			twitchURL = website.Url
		case 13:
			steamURL = website.Url
		}
	}

	var analysisToReturn playerAnalysisToReturn

	if steamURL != "" {
		var steamPlayerCount steamPlayerCountReturn
		var steamTotalPlayerCount steamPlayerCountReturn
		json.Unmarshal(fetchFromSteamAPI(getSteamGameID(steamURL)), &steamPlayerCount)
		json.Unmarshal(fetchFromSteamAPI("0"), &steamTotalPlayerCount)

		var steamAnalysis steamAnalysis
		steamAnalysis.CurrentPlayers = steamPlayerCount.Response.PlayerCount
		steamAnalysis.PercentageOfTotalPlayers = (float32(steamPlayerCount.Response.PlayerCount) / float32(steamTotalPlayerCount.Response.PlayerCount)) * 100

		analysisToReturn.Steam = &steamAnalysis
	}

	if twitchURL != "" {
		var twitchGameID twitchGameIDReturn
		json.Unmarshal(fetchFromTwitchOrGameDatabase("GET", twitchAPIgetGames, "", []queryParameter{{gameDatabaseID, gameID}}), &twitchGameID)

		var twitchAnalysis twitchAnalysis
		var twitchStreams twitchStreamsReturn
		continueFetchingStreams := true

		for continueFetchingStreams {
			json.Unmarshal(fetchFromTwitchOrGameDatabase("GET", twitchAPIgetStreams, "", []queryParameter{{twitchDatabaseID, twitchGameID.Data[0].ID}, {"after", twitchStreams.CursorData.Cursor}, {"first", "100"}}), &twitchStreams)

			streamCount := len(twitchStreams.Data)

			if streamCount == 0 {
				continueFetchingStreams = false
			} else {
				twitchAnalysis.CurrentStreams += streamCount

				for _, stream := range twitchStreams.Data {
					twitchAnalysis.CurrentViewers += stream.ViewerCount
				}
			}
		}

		if twitchAnalysis.CurrentStreams == 0 {
			twitchAnalysis.AverageViewersPerStream = 0
		} else {
			twitchAnalysis.AverageViewersPerStream = float32(twitchAnalysis.CurrentViewers) / float32(twitchAnalysis.CurrentStreams)
		}

		analysisToReturn.Twitch = &twitchAnalysis
	}

	c.IndentedJSON(http.StatusOK, analysisToReturn)
}

func getSteamGameID(url string) string {
	var idString = strings.Split(url, "/")[4]
	_, parseError := strconv.Atoi(idString)

	if parseError != nil {
		fmt.Print(parseError.Error())
		return ""
	}

	return idString
}

func fetchFromTwitchOrGameDatabase(method string, requestUrl string, body string, parameters []queryParameter) []byte {
	gameDatabaseRequest, requestError := http.NewRequest(method, requestUrl, bytes.NewBufferString(body))

	if requestError != nil {
		fmt.Print(requestError.Error())
		return nil
	}

	gameDatabaseRequest.Header.Add("Client-ID", ourClientID)
	gameDatabaseRequest.Header.Add("Authorization", ourAuthorization)

	if parameters != nil {
		addParamsToQuery := gameDatabaseRequest.URL.Query()

		for _, parameter := range parameters {
			addParamsToQuery.Add(parameter.key, parameter.value)
		}
		gameDatabaseRequest.URL.RawQuery = addParamsToQuery.Encode()
	}

	gameSearchFetch, fetchError := http.DefaultClient.Do(gameDatabaseRequest)

	if fetchError != nil {
		fmt.Print(fetchError.Error())
		return nil
	}

	gameFetchBody, bodyReadError := io.ReadAll(gameSearchFetch.Body)

	if bodyReadError != nil {
		fmt.Print(bodyReadError.Error())
		return nil
	}

	return gameFetchBody
}

func fetchFromSteamAPI(id string) []byte {
	steamAPIrequest, requestError := http.Get(steamWebAPI + id)

	if requestError != nil {
		return nil
	}

	requestBody, requestBodyReadError := io.ReadAll(steamAPIrequest.Body)

	if requestBodyReadError != nil {
		return nil
	}

	return requestBody
}
