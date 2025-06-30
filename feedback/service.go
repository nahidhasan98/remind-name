package feedback

import (
	"time"

	"github.com/nahidhasan98/remind-name/logger"
)

type FeedbackService struct {
	repo *repository
}

func NewFeedbackService() *FeedbackService {
	return &FeedbackService{
		repo: newRepository(),
	}
}

func (service *FeedbackService) SaveFeedback(data *Feedback) (*Response, error) {
	now := time.Now().Unix()
	data.UpdatedAt = now
	data.CreatedAt = now

	logger.Info("Attempting to save feedback: %+v", data)
	err := service.repo.saveFeedback(data)

	if err != nil {
		logger.Error("Failed to save feedback: %v", err)
		return nil, err
	}

	logger.Info("Feedback saved successfully for: %v", data.CreatedAt)

	return &Response{
		Status:  true,
		Message: "Thank you for your feedback! We appreciate your input.",
	}, nil
}
