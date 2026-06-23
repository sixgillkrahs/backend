package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditEvent represents a single audit log event streamed to the client.
type AuditEvent struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Timestamp string            `json:"timestamp"`
	Details   map[string]string `json:"details"`
}

var auditEventTypes = []string{"login", "changeRole", "createRole", "viewMatrix"}
var auditActorNames = []string{
	"Alice Nguyen", "Bob Chen", "Carol Smith", "David Park",
	"Eve Torres", "Frank Lee", "Grace Kim", "Henry Wu",
}

func randomAuditEvent() AuditEvent {
	//nolint:gosec
	eventType := auditEventTypes[rand.Intn(len(auditEventTypes))]
	//nolint:gosec
	actor := auditActorNames[rand.Intn(len(auditActorNames))]
	details := map[string]string{"name": actor}
	if eventType == "changeRole" {
		roles := []string{"Admin", "Moderator", "User"}
		//nolint:gosec
		details["role"] = roles[rand.Intn(len(roles))]
		//nolint:gosec
		details["user"] = auditActorNames[rand.Intn(len(auditActorNames))]
	} else if eventType == "createRole" {
		customRoles := []string{"Reviewer", "Auditor", "Support", "Developer"}
		//nolint:gosec
		details["role"] = customRoles[rand.Intn(len(customRoles))]
	}

	return AuditEvent{
		ID:        fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		Type:      eventType,
		Timestamp: time.Now().Format("15:04:05"),
		Details:   details,
	}
}

// AuditStreamHandler streams real-time audit events as Server-Sent Events (SSE).
// @Summary Audit Event Stream
// @Description Streams audit log events in real-time using Server-Sent Events (SSE). The client must send a Bearer token in the Authorization header.
// @Tags Audit
// @Produce text/event-stream
// @Security BearerAuth
// @Success 200 {string} string "SSE stream of audit events"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/audit/stream [get]
func AuditStreamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set SSE headers
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

		// Flush immediately with a comment to establish the connection
		c.Writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(c.Writer, ": connected\n\n")
		c.Writer.Flush()

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		clientGone := c.Request.Context().Done()

		for {
			select {
			case <-clientGone:
				// Client disconnected
				return
			case <-ticker.C:
				event := randomAuditEvent()
				data, err := json.Marshal(event)
				if err != nil {
					continue
				}
				fmt.Fprintf(c.Writer, "id: %s\n", event.ID)
				fmt.Fprintf(c.Writer, "event: audit\n")
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				c.Writer.Flush()
			}
		}
	}
}
