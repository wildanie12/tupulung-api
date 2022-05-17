package like

type LikeServiceInterface interface {
	Append(ID int, eventID int) error
	Delete(ID int, eventID int) error
}
