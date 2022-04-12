package user

import entity "tupulung/entities"

type UserRepositoryInterface interface {
	/*
	 * Find User
	 * -------------------------------
	 * Mencari user berdasarkan ID
	 */
	Find(id int) (entity.User, error)

	/*
	 * Find By Column
	 * -------------------------------
	 * Mencari user tunggal berdasarkan column dan value
	 */
	FindBy(field string, value string) (entity.User, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan user tunggal kedalam database
	 */
	Store(user entity.User) (entity.User, error)

	/*
	 * Update User
	 * -------------------------------
	 * Mengedit user tunggal berdasarkan ID
	 */
	Update(user entity.User, id int) (entity.User, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete user tunggal berdasarkan ID
	 */
	Delete(id int) error
}