package dto

import (
	"github.com/google/uuid"
	"time"
)

type Proposition struct {
	Id              uuid.UUID `db:"proposition_id"`
	OriginalTextUrl string    `db:"proposition_original_text_url"`
	Title           string    `db:"proposition_title"`
	Content         string    `db:"proposition_content"`
	SubmittedAt     time.Time `db:"proposition_submitted_at"`
	ImageUrl        string    `db:"proposition_image_url"`
	CreatedAt       time.Time `db:"proposition_created_at"`
	UpdatedAt       time.Time `db:"proposition_updated_at"`
	*PropositionType
}
