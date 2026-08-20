package httptransport

import (
	"github.com/gin-gonic/gin"
	"github.com/wyw14/cry052/internal/application"
	"github.com/wyw14/cry052/internal/domain/masking"
	"net/http"
	"strconv"
)

func NewRouter(s *application.Services) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/readyz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
	v := r.Group("/api/v1")
	v.GET("/sources", func(c *gin.Context) { c.JSON(200, gin.H{"items": s.Store.ListSources(c, limit(c.Query("limit")))}) })
	v.POST("/sources", func(c *gin.Context) {
		var b struct{ Name, Adapter, Owner, DSN string }
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.RegisterSource(c, b.Name, b.Adapter, b.Owner, b.DSN)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(http.StatusCreated, v)
	})
	v.POST("/sources/:id/tables", func(c *gin.Context) {
		var b struct{ Name string }
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		if err := s.RegisterTable(c, c.Param("id"), b.Name); err != nil {
			fail(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
	v.GET("/fields", func(c *gin.Context) { c.JSON(200, gin.H{"items": s.Store.ListFields(c, c.Query("source"))}) })
	v.POST("/fields", func(c *gin.Context) {
		var b struct {
			Source, Table, Name, Type string
			Sensitivity               masking.Sensitivity
		}
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.AddField(c, b.Source, b.Table, b.Name, b.Type, b.Sensitivity)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(201, v)
	})
	v.GET("/rules", func(c *gin.Context) { c.JSON(200, gin.H{"items": s.Store.ListRules(c)}) })
	v.POST("/rules", func(c *gin.Context) {
		var b struct {
			Name                        string
			Kind                        masking.RuleKind
			Pattern, Replacement, Owner string
		}
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.CreateRule(c, b.Name, b.Kind, b.Pattern, b.Replacement, b.Owner)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(201, v)
	})
	v.POST("/pipelines", func(c *gin.Context) {
		var b struct{ Name, Owner, Source string }
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.CreatePipeline(c, b.Name, b.Owner, b.Source)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(201, v)
	})
	v.POST("/pipelines/:id/approve", func(c *gin.Context) {
		var b struct{ Actor string }
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		if err := s.ApprovePipeline(c, c.Param("id"), b.Actor); err != nil {
			fail(c, err)
			return
		}
		c.Status(204)
	})
	v.POST("/pipelines/:id/preview", func(c *gin.Context) {
		var b struct {
			Source string
			Rows   int
		}
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.Preview(c, c.Param("id"), b.Source, b.Rows)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(201, v)
	})
	v.POST("/previews/:id/confirm", func(c *gin.Context) {
		var b struct{ Actor string }
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		if err := s.ConfirmPreview(c, c.Param("id"), b.Actor); err != nil {
			fail(c, err)
			return
		}
		c.Status(204)
	})
	v.POST("/batches", func(c *gin.Context) {
		var b struct {
			Pipeline, Preview, Owner, Target, Key string
			Total                                 int
		}
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.StartBatch(c, b.Pipeline, b.Preview, b.Owner, b.Target, b.Key, b.Total)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(201, v)
	})
	v.POST("/batches/:id/run", func(c *gin.Context) {
		var b struct {
			Snapshot string
			Chunk    int
		}
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		v, err := s.RunBatch(c, c.Param("id"), b.Snapshot, b.Chunk)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(200, v)
	})
	v.POST("/batches/:id/cancel", func(c *gin.Context) {
		if err := s.CancelBatch(c, c.Param("id")); err != nil {
			fail(c, err)
			return
		}
		c.Status(204)
	})
	v.POST("/rules/:id/apply", func(c *gin.Context) {
		var b struct{ Value string }
		if err := c.ShouldBindJSON(&b); err != nil {
			fail(c, err)
			return
		}
		value, err := s.ApplyRule(c, c.Param("id"), b.Value)
		if err != nil {
			fail(c, err)
			return
		}
		c.JSON(200, gin.H{"value": value})
	})
	return r
}
func limit(raw string) int {
	n, _ := strconv.Atoi(raw)
	if n < 1 || n > 100 {
		return 20
	}
	return n
}
func fail(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"code": "MASKING_REQUEST_INVALID", "message": err.Error(), "request_id": c.GetHeader("X-Request-ID")})
}
