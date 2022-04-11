package user

import entityDomain "tupulung/entities/domain"

type UserRepositoryInterface interface {
	/*
	 * Find User
	 * -------------------------------
	 * Mencari user berdasarkan ID
	 */
	Find(id int) (entityDomain.User, error)

	/*
	 * Find By Column
	 * -------------------------------
	 * Mencari user tunggal berdasarkan column dan value
	 */
	FindBy(field string, value string) (entityDomain.User, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan user tunggal kedalam database
	 */
	Store(user entityDomain.User) (entityDomain.User, error)

	/*
	 * Update User
	 * -------------------------------
	 * Mengedit user tunggal berdasarkan ID
	 */
	Update(user entityDomain.User, id int) (entityDomain.User, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete user tunggal berdasarkan ID
	 */
	Delete(id int) error
}