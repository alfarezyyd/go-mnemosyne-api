package sharedNote

type Handler struct {
	sharedNoteService Service
}

func NewHandler() *Handler {
	return &Handler{}
}
