package product

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/etg"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
)

type ProductType int

const (
	LiveCasino   ProductType = 1
	Slots        ProductType = 2
	Sports       ProductType = 3
	CardAndBoard ProductType = 5
	Fishing      ProductType = 6
	Cockfighting ProductType = 9
)

func (p ProductType) String() string {
	switch p {
	case LiveCasino:
		return "Live Casino"
	case Slots:
		return "Slots"
	case Sports:
		return "Sports"
	case CardAndBoard:
		return "Card & Board"
	case Fishing:
		return "Fishing"
	case Cockfighting:
		return "Cockfighting"
	default:
		return ""
	}
}

var productTypeToUriComponentMap = map[ProductType]string{
	LiveCasino:   "live-casino",
	Slots:        "slots",
	Sports:       "sports",
	CardAndBoard: "card-and-board",
	Fishing:      "fishing",
	Cockfighting: "cockfighting",
}

func (p ProductType) UriComponent() string {
	return productTypeToUriComponentMap[p]
}

func UriComponentToProductType(uri string) ProductType {
	for productType, uriComponent := range productTypeToUriComponentMap {
		if uri == uriComponent {
			return productType
		}
	}
	return 0
}

type Product int

const (
	MainWallet     Product = -1
	_1GPoker       Product = 146
	_3Win8         Product = 110
	Allbet         Product = 7
	AsiaGaming     Product = 1002
	CreativeGaming Product = 127
	GamingSoft     Product = 2
	IBCbet         Product = 86
	Jdb            Product = 41
	Jili           Product = 59
	Joker          Product = 22
	M8Bet          Product = 17
	MPoker         Product = 162
	PGSoft         Product = 102
	PlayNGo        Product = 204
	PragmaticPlay  Product = 6
	RedTiger       Product = 116
	SAGaming       Product = 1001
	SBObet         Product = 80
	SexyBaccarat   Product = 3
	Spadegaming    Product = 12
	SV388          Product = 32
	TFGaming       Product = 138
	WMCasino       Product = 70
	WSGaming       Product = 147
)

func (p Product) String() string {
	switch p {
	case MainWallet:
		return "Main Wallet"
	case _1GPoker:
		return "1G Poker"
	case _3Win8:
		return "3Win8"
	case Allbet:
		return "Allbet"
	case AsiaGaming:
		return "Asia Gaming"
	case CreativeGaming:
		return "Creative Gaming"
	case GamingSoft:
		return "Gaming Soft"
	case IBCbet:
		return "IBCbet"
	case Jdb:
		return "JDB"
	case Jili:
		return "JILI"
	case Joker:
		return "Joker"
	case M8Bet:
		return "M8Bet"
	case MPoker:
		return "M Poker"
	case PGSoft:
		return "PG Soft"
	case PlayNGo:
		return "Play N Go"
	case PragmaticPlay:
		return "Pragmatic Play"
	case RedTiger:
		return "Red Tiger"
	case SAGaming:
		return "SA Gaming"
	case SBObet:
		return "SBObet"
	case SexyBaccarat:
		return "Sexy Baccarat"
	case Spadegaming:
		return "Spade Gaming"
	case SV388:
		return "SV388"
	case TFGaming:
		return "TF Gaming"
	case WMCasino:
		return "WM Casino"
	case WSGaming:
		return "WS Gaming"
	default:
		return ""
	}
}

var productToUriComponentMap = map[Product]string{
	_1GPoker:       "1g-poker",
	_3Win8:         "3win8",
	Allbet:         "allbet",
	AsiaGaming:     "asia-gaming",
	CreativeGaming: "creative-gaming",
	GamingSoft:     "gaming-soft",
	IBCbet:         "ibcbet",
	Jdb:            "jdb",
	Jili:           "jili",
	Joker:          "joker",
	M8Bet:          "m8bet",
	MPoker:         "m-poker",
	PGSoft:         "pgsoft",
	PlayNGo:        "playngo",
	PragmaticPlay:  "pragmatic-play",
	RedTiger:       "red-tiger",
	SAGaming:       "sa-gaming",
	SBObet:         "sbobet",
	SexyBaccarat:   "sexy-baccarat",
	Spadegaming:    "spadegaming",
	SV388:          "sv388",
	TFGaming:       "tf-gaming",
	WMCasino:       "wm-casino",
	WSGaming:       "ws-gaming",
}

func (p Product) UriComponent() string {
	return productToUriComponentMap[p]
}

func UriComponentToProduct(uri string) Product {
	for product, uriComponent := range productToUriComponentMap {
		if uri == uriComponent {
			return product
		}
	}
	return 0
}

var AllProducts = []Product{
	_1GPoker,
	_3Win8,
	Allbet,
	AsiaGaming,
	CreativeGaming,
	GamingSoft,
	IBCbet,
	Jdb,
	Jili,
	Joker,
	M8Bet,
	MPoker,
	PGSoft,
	PlayNGo,
	PragmaticPlay,
	RedTiger,
	SAGaming,
	SBObet,
	SexyBaccarat,
	Spadegaming,
	SV388,
	TFGaming,
	WMCasino,
	WSGaming,
}
var LiveCasinoProducts = []Product{
	PragmaticPlay,
	SexyBaccarat,
	PlayNGo,
	Allbet,
	AsiaGaming,
	SAGaming,
	WMCasino,
}
var SlotsProducts = []Product{
	PragmaticPlay,
	RedTiger,
	Jdb,
	Joker,
	Jili,
	PGSoft,
	PlayNGo,
	_3Win8,
	AsiaGaming,
	CreativeGaming,
	GamingSoft,
	Spadegaming,
}
var SportsProducts = []Product{SBObet, IBCbet, M8Bet, TFGaming, WSGaming}
var CardAndBoardProducts = []Product{_1GPoker, MPoker}
var FishingProducts = []Product{Jdb, Jili}
var CockfightingProducts = []Product{SV388}

var WalletGroup1 = []Product{AsiaGaming, SAGaming}
var WalletGroup2 = []Product{Jdb, Jili, SexyBaccarat, SV388}

func HasGameList(productType ProductType, product Product) bool {
	if productType == Slots {
		return product == CreativeGaming || product == Joker || product == PragmaticPlay || product == Spadegaming || product == PGSoft || product == PlayNGo
	} else if productType == CardAndBoard {
		return product == _1GPoker
	}
	return false
}

type LaunchGameListRequest struct {
	Op   string      `json:"op"`
	Prod Product     `json:"prod"`
	Type ProductType `json:"type"`
	Mem  string      `json:"mem"`
	Pass string      `json:"pass"`
}

type LaunchGameListResponse struct {
	Url  string `json:"url"`
	Err  int    `json:"err"`
	Desc string `json:"desc"`
}

func LaunchGameList(etgUsername string, productType ProductType, product Product) (string, error) {
	// We get the player's balance to ensure the player's account is successfully
	// created in ETG
	_, err := GetUserBalance(etgUsername, product)
	if err != nil {
		log.Panicf("Can't get balance of %s (%d) for %s: %s\n", product, product, etgUsername, err)
	}

	payload := LaunchGameListRequest{
		Op:   os.Getenv("ETG_OPERATOR_CODE"),
		Prod: product,
		Type: productType,
		Mem:  etgUsername,
		Pass: "00000000",
	}
	var launchGameListResponse LaunchGameListResponse
	err = etg.Post("/game", payload, &launchGameListResponse)
	if err != nil {
		log.Panicf("Can't launch %s (%d) game list for %s: %s\nEndpoint: %s\nPayload: %+v\n", product, product, etgUsername, err, "/game", payload)
	}

	if launchGameListResponse.Err != etg.Success {
		return "", fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %s", launchGameListResponse.Err, launchGameListResponse.Desc, "/game", payload)
	}

	return launchGameListResponse.Url, nil
}

type GameListRequest struct {
	Op   string      `json:"op"`
	Prod Product     `json:"prod"`
	Type ProductType `json:"type"`
}

type GameListResponse struct {
	Data []DataItem `json:"data"`
	Err  int        `json:"err"`
	Desc string     `json:"desc"`
}

type DataItem struct {
	CName    string `json:"c_name"`
	CWebCode string `json:"c_web_code"`
	CImage   string `json:"c_image"`
}

type GameInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ImageLink string `json:"imageLink"`
}

func GetGameList(app *app.App, ctx context.Context, productType ProductType, product Product) ([]GameInfo, error) {
	// Get from redis cache
	key := "games:" + productType.UriComponent() + "_" + product.UriComponent()
	cached, err := app.Redis.Get(ctx, key).Result()
	if err == nil {
		buf := bytes.NewBufferString(cached)
		var games []GameInfo
		dec := gob.NewDecoder(buf)
		err := dec.Decode(&games)
		if err == nil {
			return games, nil
		}
	}

	payload := GameListRequest{
		Op:   os.Getenv("ETG_OPERATOR_CODE"),
		Prod: product,
		Type: productType,
	}
	var gameListResponse GameListResponse
	err = etg.Post("/getgamelist", payload, &gameListResponse)
	if err != nil {
		log.Panicf("Can't get %s (%d) game list: %s\nEndpoint: %s\nPayload: %+v\n", product, product, err, "/getgamelist", payload)
	}

	if gameListResponse.Err != etg.Success {
		return nil, fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", gameListResponse.Err, gameListResponse.Desc, "/getgamelist", payload)
	}

	games := sliceutils.Map(gameListResponse.Data, func(data DataItem) GameInfo {
		return GameInfo{
			Id:        data.CWebCode,
			Name:      data.CName,
			ImageLink: data.CImage,
		}
	})

	// Cache in redis
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(games)
	app.Redis.Set(ctx, key, buf.String(), time.Hour)

	return games, nil
}

type GameRequest struct {
	Op     string      `json:"op"`
	Prod   Product     `json:"prod"`
	Type   ProductType `json:"type"`
	GameId string      `json:"gameid"`
	H5     int         `json:"h5"`
	Mem    string      `json:"mem"`
	Pass   string      `json:"pass"`
}

type GameResponse struct {
	Url  string `json:"url"`
	Err  int    `json:"err"`
	Desc string `json:"desc"`
}

func LaunchGame(etgUsername string, productType ProductType, product Product, gameId string) (string, error) {
	payload := GameRequest{
		Op:     os.Getenv("ETG_OPERATOR_CODE"),
		Prod:   product,
		Type:   productType,
		GameId: gameId,
		H5:     1,
		Mem:    etgUsername,
		Pass:   "00000000",
	}
	var gameResponse GameResponse
	err := etg.Post("/game", payload, &gameResponse)
	if err != nil {
		log.Panicln("Can't launch game: " + err.Error())
	}

	if gameResponse.Err != etg.Success {
		return "", fmt.Errorf("%d: %s", gameResponse.Err, gameResponse.Desc)
	}

	return gameResponse.Url, nil
}
