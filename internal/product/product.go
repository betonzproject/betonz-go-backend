package product

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/doorman2137/betonz-go/internal/etg"
	"github.com/doorman2137/betonz-go/internal/utils/sliceutils"
	"github.com/redis/go-redis/v9"
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

var AllProductTypes = []ProductType{
	LiveCasino,
	Slots,
	Sports,
	CardAndBoard,
	Fishing,
	Cockfighting,
}

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
	BigGaming      Product = 5
	CreativeGaming Product = 127
	CQ9            Product = 39
	Dreamtech      Product = 120
	KingMaker      Product = 117
	IBCbet         Product = 86
	Jdb            Product = 41
	Jili           Product = 59
	Joker          Product = 22
	M8Bet          Product = 17
	MPoker         Product = 162
	PGSoft         Product = 102
	Live22         Product = 9
	PlayNGo        Product = 204
	PragmaticPlay  Product = 6
	SBObet         Product = 80
	SexyBaccarat   Product = 3
	Spadegaming    Product = 12
	SV388          Product = 32
	TFGaming       Product = 138
	WMCasino       Product = 70
	WSGaming       Product = 147
	YLfishing      Product = 52
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
	case BigGaming:
		return "Big Gaming"
	case CreativeGaming:
		return "Creative Gaming"
	case CQ9:
		return "CQ9"
	case Dreamtech:
		return "Dream Tech"
	case IBCbet:
		return "IBCbet"
	case Jdb:
		return "JDB"
	case Jili:
		return "JILI"
	case Joker:
		return "Joker"
	case KingMaker:
		return "King Maker"
	case M8Bet:
		return "M8Bet"
	case MPoker:
		return "M Poker"
	case PGSoft:
		return "PG Soft"
	case Live22:
		return "Live22"
	case PlayNGo:
		return "Play N Go"
	case PragmaticPlay:
		return "Pragmatic Play"
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
	case YLfishing:
		return "YL Fishing"
	default:
		return ""
	}
}

var productToUriComponentMap = map[Product]string{
	_1GPoker:       "1g-poker",
	_3Win8:         "3win8",
	Allbet:         "allbet",
	AsiaGaming:     "asia-gaming",
	BigGaming:      "big-gaming",
	CreativeGaming: "creative-gaming",
	CQ9:            "cq9",
	Dreamtech:      "dream-tech",
	IBCbet:         "ibcbet",
	Jdb:            "jdb",
	Jili:           "jili",
	Joker:          "joker",
	KingMaker:      "king-maker",
	M8Bet:          "m8bet",
	MPoker:         "m-poker",
	PGSoft:         "pgsoft",
	PlayNGo:        "playngo",
	PragmaticPlay:  "pragmatic-play",
	SBObet:         "sbobet",
	SexyBaccarat:   "sexy-baccarat",
	Spadegaming:    "spadegaming",
	SV388:          "sv388",
	TFGaming:       "tf-gaming",
	WMCasino:       "wm-casino",
	WSGaming:       "ws-gaming",
	Live22:         "live-22",
	YLfishing:      "yl-fishing",
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
	BigGaming,
	CreativeGaming,
	CQ9,
	Dreamtech,
	IBCbet,
	Jdb,
	Jili,
	Joker,
	KingMaker,
	M8Bet,
	MPoker,
	PGSoft,
	Live22,
	PlayNGo,
	PragmaticPlay,
	SBObet,
	SexyBaccarat,
	Spadegaming,
	SV388,
	TFGaming,
	WMCasino,
	WSGaming,
	YLfishing,
}
var LiveCasinoProducts = []Product{
	PragmaticPlay,
	SexyBaccarat,
	Allbet,
	AsiaGaming,
	BigGaming,
	WMCasino,
}
var SlotsProducts = []Product{
	PragmaticPlay,
	Jdb,
	Joker,
	Jili,
	PGSoft,
	Live22,
	PlayNGo,
	_3Win8,
	AsiaGaming,
	CreativeGaming,
	CQ9,
	Dreamtech,
	Spadegaming,
}
var SportsProducts = []Product{SBObet, IBCbet, M8Bet, TFGaming, WSGaming}
var CardAndBoardProducts = []Product{_1GPoker, MPoker, KingMaker}
var FishingProducts = []Product{Jdb, Jili, YLfishing}
var CockfightingProducts = []Product{SV388}

var WalletGroup1 = []Product{AsiaGaming}
var WalletGroup2 = []Product{Jdb, Jili, SexyBaccarat, SV388}

func HasGameList(productType ProductType, product Product) bool {
	if productType == Slots {
		return product == CreativeGaming || product == Joker || product == PragmaticPlay || product == Spadegaming || product == PlayNGo || product == CQ9 || product == Dreamtech
	} else if productType == CardAndBoard {
		return product == _1GPoker
	}
	return false
}

func SharesSameWallet(p1 Product, p2 Product) bool {
	if p1 == p2 {
		return true
	}
	return (slices.Contains(WalletGroup1, p1) && slices.Contains(WalletGroup1, p2) ||
		slices.Contains(WalletGroup2, p1) && slices.Contains(WalletGroup2, p2))
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
		return "", fmt.Errorf("Can't get balance of %s (%d) for %s: %s", product, product, etgUsername, err)
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
		return "", err
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

func GetGameList(redis *redis.Client, ctx context.Context, productType ProductType, product Product) ([]GameInfo, error) {
	// Get from redis cache
	key := "games:" + productType.UriComponent() + "_" + product.UriComponent()
	cached, err := redis.Get(ctx, key).Result()
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
		return nil, err
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
	redis.Set(ctx, key, buf.String(), time.Hour)

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
		return "", err
	}

	if gameResponse.Err != etg.Success {
		return "", fmt.Errorf("%d: %s", gameResponse.Err, gameResponse.Desc)
	}

	return gameResponse.Url, nil
}

var PragmaticPlayTopGame = []string{
	"Gates of Olympus",
	"Sweet Bonanza",
	"Sugar Rush 1000",
	"Gates of Olympus 1000",
	"Zeus vs Hades - Gods of War",
	"The Dog House",
	"Starlight Princess",
	"Starlight Princess 1000",
	"Gates of Olympus DICE",
	"Fire Portals",
	"The Dog House - Dog or Alive",
	"Big Bass - Secrets of the Golden Lake",
	"Fruit Party",
	"Sugar Rush",
	"Cleocatra",
}
