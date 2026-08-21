package integrations

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"green-compass-backend/pkg/httpx"
)

// Handler receives provider callbacks (SMS inbound/delivery, USSD sessions,
// push receipts, source callbacks) with signature validation and replay
// protection per the integration contract.
type Handler struct {
	webhookSecret string
	seenKeys      *replayCache
}

func NewHandler(webhookSecret string) *Handler {
	return &Handler{webhookSecret: webhookSecret, seenKeys: newReplayCache(5 * time.Minute)}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/v1/integrations")

	group.POST("/sms/inbound", h.SMSInbound)
	group.POST("/sms/delivery", h.SMSDelivery)
	group.POST("/ussd/session", h.USSDSession)
	group.POST("/push/receipt", h.PushReceipt)
	group.POST("/sources/:source_key/callback", h.SourceCallback)
}

// validateSignature checks the X-Signature header (HMAC-SHA256 of raw body).
func (h *Handler) validateSignature(c *gin.Context, body []byte) bool {
	if h.webhookSecret == "" {
		return true // unsigned in local dev
	}
	sig := c.GetHeader("X-Signature")
	if sig == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sig), []byte(expected))
}

func readBody(c *gin.Context) ([]byte, bool) {
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20)) // 1 MiB cap
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "unable to read request body"))
		return nil, false
	}
	return body, true
}

// SMSInbound receives inbound SMS messages from the SMS provider.
// POST /v1/integrations/sms/inbound
func (h *Handler) SMSInbound(c *gin.Context) {
	body, ok := readBody(c)
	if !ok {
		return
	}
	if !h.validateSignature(c, body) {
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_SIGNATURE", Message: "invalid webhook signature", Status: http.StatusUnauthorized})
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"received": true}, httpx.GetRequestID(c)))
}

// SMSDelivery receives delivery receipts from the SMS provider.
// POST /v1/integrations/sms/delivery
func (h *Handler) SMSDelivery(c *gin.Context) {
	body, ok := readBody(c)
	if !ok {
		return
	}
	if !h.validateSignature(c, body) {
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_SIGNATURE", Message: "invalid webhook signature", Status: http.StatusUnauthorized})
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"received": true}, httpx.GetRequestID(c)))
}

// USSDSession handles USSD session callbacks from the aggregator.
// POST /v1/integrations/ussd/session
func (h *Handler) USSDSession(c *gin.Context) {
	body, ok := readBody(c)
	if !ok {
		return
	}
	if !h.validateSignature(c, body) {
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_SIGNATURE", Message: "invalid webhook signature", Status: http.StatusUnauthorized})
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"received": true}, httpx.GetRequestID(c)))
}

// PushReceipt receives push delivery receipts (FCM/APNS).
// POST /v1/integrations/push/receipt
func (h *Handler) PushReceipt(c *gin.Context) {
	body, ok := readBody(c)
	if !ok {
		return
	}
	if !h.validateSignature(c, body) {
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_SIGNATURE", Message: "invalid webhook signature", Status: http.StatusUnauthorized})
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"received": true}, httpx.GetRequestID(c)))
}

// SourceCallback receives external source callbacks with replay protection.
// POST /v1/integrations/sources/{source_key}/callback
func (h *Handler) SourceCallback(c *gin.Context) {
	sourceKey := c.Param("source_key")
	if sourceKey == "" {
		httpx.HandleError(c, httpx.InvalidParam("source_key", "source key is required"))
		return
	}

	body, ok := readBody(c)
	if !ok {
		return
	}
	if !h.validateSignature(c, body) {
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_SIGNATURE", Message: "invalid webhook signature", Status: http.StatusUnauthorized})
		return
	}

	// Replay protection via idempotency key
	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey != "" {
		if h.seenKeys.seen(idemKey) {
			c.JSON(http.StatusOK, httpx.Success(gin.H{"duplicate": true}, httpx.GetRequestID(c)))
			return
		}
		h.seenKeys.mark(idemKey)
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"source":  sourceKey,
		"received": true,
	}, httpx.GetRequestID(c)))
}

// replayCache is a simple TTL-based dedup cache for webhook keys.
type replayCache struct {
	ttl   time.Duration
	items map[string]time.Time
}

func newReplayCache(ttl time.Duration) *replayCache {
	return &replayCache{ttl: ttl, items: make(map[string]time.Time)}
}

func (r *replayCache) seen(key string) bool {
	ts, ok := r.items[key]
	if !ok {
		return false
	}
	if time.Since(ts) > r.ttl {
		delete(r.items, key)
		return false
	}
	return true
}

func (r *replayCache) mark(key string) {
	r.items[key] = time.Now()
}