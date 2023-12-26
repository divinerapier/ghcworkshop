package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartAPI(t *testing.T) {
	r := NewEngine()

	t.Run("POST /cart", func(t *testing.T) {
		cart := &Cart{ID: 1}
		body, _ := json.Marshal(cart)
		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)
		var responseCart Cart
		json.Unmarshal(resp.Body.Bytes(), &responseCart)
		assert.Equal(t, cart.ID, responseCart.ID)
	})

	t.Run("GET /cart/:id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/cart/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)
		var responseCart Cart
		json.Unmarshal(resp.Body.Bytes(), &responseCart)
		assert.Equal(t, 1, responseCart.ID)
	})
}

func TestCartOperations(t *testing.T) {
	r := NewEngine()

	t.Run("Add items to cart", func(t *testing.T) {
		for _, id := range []int{1, 2, 3} {
			cart := &Cart{ID: id}
			body, _ := json.Marshal(cart)
			req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, 200, resp.Code)
		}
	})

	t.Run("Check items in cart", func(t *testing.T) {
		for _, id := range []int{1, 2, 3} {
			req, _ := http.NewRequest(http.MethodGet, "/cart/"+strconv.Itoa(id), nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, 200, resp.Code)
			// 校验返回的 beer 是否存在 beer id 1 2 3
			var responseCart Cart
			json.Unmarshal(resp.Body.Bytes(), &responseCart)
			assert.Equal(t, id, responseCart.ID)
			for _, beer := range responseCart.Beers {
				assert.Contains(t, []int{1, 2, 3}, beer.ID)
			}
		}
	})

	t.Run("Remove item from cart", func(t *testing.T) {
		// 调用 Put 接口，传入 remove 2
		body := []byte(`{"beer_id":2,"op":"remove"}`)
		req, _ := http.NewRequest(http.MethodPut, "/cart/2", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Check items in cart after removal", func(t *testing.T) {
		for _, id := range []int{1, 3} {
			req, _ := http.NewRequest(http.MethodGet, "/cart/"+strconv.Itoa(id), nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, 200, resp.Code)
			// 校验返回的 beer 是否存在 beer id 1  3
			var responseCart Cart
			json.Unmarshal(resp.Body.Bytes(), &responseCart)
			assert.Equal(t, id, responseCart.ID)
			for _, beer := range responseCart.Beers {
				assert.Contains(t, []int{1, 3}, beer.ID)
			}
		}
	})

	t.Run("Add new item to cart", func(t *testing.T) {
		// 调用 Put 接口，传入 add 10
		body := []byte(`{"beer_id":10,"op":"add"}`)
		req, _ := http.NewRequest(http.MethodPut, "/cart/1", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Check items in cart after adding new item", func(t *testing.T) {
		for _, id := range []int{1} {
			req, _ := http.NewRequest(http.MethodGet, "/cart/"+strconv.Itoa(id), nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, 200, resp.Code)
			// 判断返回的 beer 中的id 应该是 1 3 10
			var responseCart Cart
			json.Unmarshal(resp.Body.Bytes(), &responseCart)
			assert.Equal(t, id, responseCart.ID)
			for _, beer := range responseCart.Beers {
				assert.Contains(t, []int{1, 3, 10}, beer.ID)
			}
		}
	})
}

// 测试 search 接口
func TestSearchAPI(t *testing.T) {
	r := NewEngine()

	t.Run("Search beers by name", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/beer/search?q=IPA", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)
		var beers []Beer
		json.Unmarshal(resp.Body.Bytes(), &beers)
		assert.Equal(t, 3, len(beers))
	})
}
