package like

import "tupulung/entities"

type LikeRepositoryInterface interface {

	/*
	 * Count like
	 * -------------------------------
	 * Menghitung like berdasarkan event
	 */
	CountLikeByEvent(eventId int) (int64, error)

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
