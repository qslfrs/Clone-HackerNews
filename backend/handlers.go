package main

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/topstories?page=1&limit=20
func GetTopStoriesHandler(c *gin.Context, client *HNClient) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	ids, err := client.GetTopStoryIDs(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed fetch topstories", "detail": err.Error()})
		return
	}

	// calculate slice for pagination
	start := (page - 1) * limit
	if start >= len(ids) {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "total": len(ids)})
		return
	}
	end := start + limit
	if end > len(ids) {
		end = len(ids)
	}
	pageIDs := ids[start:end]

	// concurrency limit semaphore
	sem := make(chan struct{}, 10) // max 10 concurrent item fetches
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	items := make([]map[string]interface{}, 0, len(pageIDs))

	for _, id := range pageIDs {
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
				// skip failed item but continue others
				return
			}
			CacheSet(key, it, 5*time.Minute)
			mu.Lock()
			items = append(items, it)
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	// keep original order according to pageIDs
	// build map for quick lookup
	orderMap := make(map[int]map[string]interface{}, len(items))
	for _, it := range items {
		if idf, ok := it["id"].(float64); ok {
			orderMap[int(idf)] = it
		}
	}
	ordered := make([]map[string]interface{}, 0, len(pageIDs))
	for _, id := range pageIDs {
		if it, ok := orderMap[id]; ok {
			ordered = append(ordered, it)
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": ordered, "total": len(ids)})
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
