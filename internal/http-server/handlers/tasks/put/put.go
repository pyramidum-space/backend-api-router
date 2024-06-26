package put

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pyramidum-space/protos/gen/go/tasks"
)

type Response struct {
	ID int64 `json:"user_id"`
}

type UpdateTaskRequest struct {
	ID               []byte               `json:"id"`
	Header           string               `json:"header"`
	Text             string               `json:"text"`
	ExternalImages   []string             `json:"external_images"`
	Deadline         time.Time            `json:"deadline"`
	ProgressStatus   tasks.ProgressStatus `json:"progress_status"`
	IsUrgent         bool                 `json:"is_urgent"`
	IsImportant      bool                 `json:"is_important"`
	OwnerID          int32                `json:"owner_id"`
	ParentID         []byte               `json:"parent_id"`
	PossibleDeadline time.Time            `json:"possible_deadline"`
	Weight           int32                `json:"weight"`
}

type TaskUpdator interface {
	Update(
		id []byte,
		header string,
		text string,
		external_images []string,
		deadline time.Time,
		progress_status tasks.ProgressStatus,
		is_urgent bool,
		is_important bool,
		owner_id int32,
		parent_id []byte,
		possible_deadline time.Time,
		weight int32) ([]byte, error)
}

func MakeGetHandlerFunc(log *slog.Logger, updator TaskUpdator) gin.HandlerFunc {
	const op = "http-server.handlers.tasks.put.MakeGetHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(c *gin.Context) {

		var req UpdateTaskRequest

		if err := c.BindJSON(&req); err != nil {
			log.Error("err: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := updator.Update(req.ID,
			req.Header,
			req.Text,
			req.ExternalImages,
			req.Deadline,
			req.ProgressStatus,
			req.IsUrgent,
			req.IsImportant,
			req.OwnerID,
			req.ParentID,
			req.PossibleDeadline,
			req.Weight)

		if err != nil {
			log.Error("error while registration")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"task_id": id,
		})
	}
}
