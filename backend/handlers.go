package main

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/topstories?page=1&limit=20
func GetTopStoriesHandler(c *gin.Context, client *HNClient) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	typeFilter := c.Query("type")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	ids, err := client.GetTopStoryIDs(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed fetch topstories", "detail": err.Error()})
		return
	}

	// concurrency limit semaphore
	sem := make(chan struct{}, 10) // batasi 10 concurrent fetch
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	items := make([]map[string]interface{}, 0, len(ids))

	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			key := "item:" + strconv.Itoa(id)
			if v, ok := CacheGet(key); ok {
				if itm, ok2 := v.(map[string]interface{}); ok2 {
					mu.Lock()
					items = append(items, itm)
					mu.Unlock()
					return
				}
			}

			it, err := client.GetItem(ctx, id)
			if err != nil {
				return
			}
			CacheSet(key, it, 10*time.Minute)
			mu.Lock()
			items = append(items, it)
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	// filter berdasarkan type jika diberikan
	filtered := items
	if typeFilter != "" {
		tmp := make([]map[string]interface{}, 0, len(items))
		for _, it := range items {
			if t, ok := it["type"].(string); ok {
				if t == typeFilter {
					tmp = append(tmp, it)
				}
			}
		}
		filtered = tmp
	}

	// urutkan semua item dari terbaru ke terlama
	sort.Slice(filtered, func(i, j int) bool {
		ti, _ := filtered[i]["time"].(float64)
		tj, _ := filtered[j]["time"].(float64)
		return ti > tj
	})

	// pagination setelah sorting/filter
	total := len(filtered)
	start := (page - 1) * limit
	if start >= total {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "total": total})
		return
	}
	end := start + limit
	if end > total {
		end = total
	}
	pagedItems := filtered[start:end]

	c.JSON(http.StatusOK, gin.H{"items": pagedItems, "total": len(items)})
}

// GET /api/item/:id
func GetItemHandler(c *gin.Context, client *HNClient) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	key := "item:" + idStr
	if v, ok := CacheGet(key); ok {
		c.JSON(http.StatusOK, v)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	it, err := client.GetItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	CacheSet(key, it, 10*time.Minute)
	c.JSON(http.StatusOK, it)
}

// GET /api/user/:id
func GetUserHandler(c *gin.Context, client *HNClient) {
	user := c.Param("id")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user"})
		return
	}
	key := "user:" + user
	if v, ok := CacheGet(key); ok {
		c.JSON(http.StatusOK, v)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	u, err := client.GetUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	CacheSet(key, u, 30*time.Minute)
	c.JSON(http.StatusOK, u)
}
