package dto

type NoteDto interface {
	GetDueDate() string
	SetDueDate(string)
}
