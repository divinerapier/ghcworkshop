package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Beer struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Tagline          string  `json:"tagline"`
	FirstBrewed      string  `json:"first_brewed"`
	Description      string  `json:"description"`
	ImageURL         string  `json:"image_url"`
	Price            float64 `json:"price"`
	Abv              float64 `json:"abv"`
	Ibu              float64 `json:"ibu"`
	TargetFg         float64 `json:"target_fg"`
	TargetOg         float64 `json:"target_og"`
	Ebc              float64 `json:"ebc"`
	Srm              float64 `json:"srm"`
	Ph               float64 `json:"ph"`
	AttenuationLevel float64 `json:"attenuation_level"`
	Volume           struct {
		Value int    `json:"value"`
		Unit  string `json:"unit"`
	} `json:"volume"`
	BoilVolume struct {
		Value int    `json:"value"`
		Unit  string `json:"unit"`
	} `json:"boil_volume"`
	FoodPairing   []string `json:"food_pairing"`
	BrewersTips   string   `json:"brewers_tips"`
	ContributedBy string   `json:"contributed_by"`
}

type Cart struct {
	ID    int    `json:"id"`
	Beers []Beer `json:"beers"`
}

func main() {
	r := NewEngine()
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}

func NewEngine() *gin.Engine {
	r := gin.New()

	beers, err := loadBeersFromFile("./beers.json")
	if err != nil {
		panic(err)
	}

	carts := NewCarts()

	{

		// curl -X POST http://localhost:8080/cart
		// Create cart
		r.POST("/cart", func(ctx *gin.Context) {
			cart := NewCart()
			carts.Add(cart)
			ctx.JSON(200, cart)
		})

		// curl -X GET http://localhost:8080/cart/1
		// Get cart by id
		r.GET("/cart/:id", func(ctx *gin.Context) {
			id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
			if err != nil {
				ctx.JSON(400, "无效的购物车 ID")
				return
			}
			cart := carts.Get(int(id))
			ctx.JSON(200, cart)
		})

		// curl -X PUT http://localhost:8080/cart/1 -d '{"beer_id": 1, "op": "add"}'
		// curl -X PUT http://localhost:8080/cart/1 -d '{"beer_id": 1, "op": "remove"}'
		// Update(add / remove) cart by id.
		r.PUT("/cart/:id", func(ctx *gin.Context) {
			type Args struct {
				BeerID int    `json:"beer_id" form:"beer_id" binding:"required"`
				Op     string `json:"op" form:"op" binding:"required"`
			}
			id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
			if err != nil {
				ctx.JSON(400, "无效的购物车 ID")
				return
			}
			var args Args
			if err := ctx.BindJSON(&args); err != nil {
				ctx.JSON(400, "无效的啤酒 ID")
				return
			}
			if args.Op != "add" && args.Op != "remove" {
				ctx.JSON(400, "无效的操作")
				return
			}
			beer, exists := beers.GetByID(int(id))
			if !exists {
				ctx.String(404, "无效的啤酒 ID")
				return
			}
			cart := carts.Get(int(id))
			switch args.Op {
			case "add":
				cart.AddBeer(*beer)
			case "remove":
				cart.RemoveBeer(*beer)
			default:
			}
			ctx.JSON(200, cart)
		})

		r.DELETE("/cart/:id", func(ctx *gin.Context) {
			id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
			if err != nil {
				ctx.JSON(400, "无效的购物车 ID")
				return
			}
			carts.Remove(int(id))
			ctx.JSON(200, id)
		})
	}

	beerGroup := r.Group("/beer")
	{
		beerGroup.GET("/", func(ctx *gin.Context) {
			remains := beers.GetRemains()
			ctx.JSON(200, remains)
		})

		beerGroup.GET("/search", func(ctx *gin.Context) {
			type Search struct {
				Q string `form:"q" binding:"required"`
			}
			// search from all bears if name / desc / tagline / food_pairing contains x
			var search Search
			if err := ctx.BindQuery(&search); err != nil {
				ctx.JSON(400, "无效的搜索参数")
				return
			}
			remains := beers.GetRemains()
			result := make([]Beer, 0)
			x := search.Q
			for _, beer := range remains {
				if strings.Contains(beer.Name, x) || strings.Contains(beer.Description, x) || strings.Contains(beer.Tagline, x) {
					result = append(result, beer)
					continue
				}
				for _, food := range beer.FoodPairing {
					if strings.Contains(food, x) {
						result = append(result, beer)
						break
					}
				}
			}
			ctx.JSON(200, result)
		})
	}
	return r
}
