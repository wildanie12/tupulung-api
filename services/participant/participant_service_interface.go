package participant

type ParticipantServiceInterface interface {
	Append(userID, eventID int) error
	Delete(userID, eventID int) error
}
