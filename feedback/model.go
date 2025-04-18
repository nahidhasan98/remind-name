package feedback

type Feedback struct {
	Name      string `bson:"name" form:"name" binding:"required"`
	Email     string `bson:"email" form:"email" binding:"required"`
	Feedback  string `bson:"feedback" form:"feedbackText" binding:"required"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

type Response struct {
	Status  bool
	Message string
}
