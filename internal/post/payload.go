package post

type PostCreateRequest struct {
	Image       string `json:"image" validate:"required,url"`
	Description string `json:"description" validate:"required"`
}

type PostUpdateRequest struct {
	Image       string `json:"image" validate:"omitempty,url"`
	Description string `json:"description"`
}

type GetPostsResponse struct {
	Posts []Post `json:"posts"`
	Count int64  `json:"count"`
}
