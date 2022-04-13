package participant

import "tupulung/entities"

type ParticipantRepositoryInterface interface {

	/*
	 * Append
	 * -------------------------------
	 * Menambahkan user ke event
	 */
	Append(user entities.User, event entities.Event) error

	/*
	 * Delete
	 * -------------------------------
	 * Menghapus user dari event
	 */
	Delete(user entities.User, event entities.Event) error
}
