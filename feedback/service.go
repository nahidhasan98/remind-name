package feedback

import (
	"time"
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

	err := service.repo.saveFeedback(data)

	if err != nil {
		return nil, err
	}

	return &Response{
		Status:  true,
		Message: "Thank you for your feedback! We appreciate your input.",
	}, nil
}
