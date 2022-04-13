package comment

import "tupulung/entities"

type CommentRepositoryInterface interface{
	/*
	 * Find All
	 * -------------------------------
	 * Mengambil data comment berdasarkan filters dan sorts
	 */
	FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Comment, error)

	/*
	 * Find
	 * -------------------------------
	 * Mencari comment tunggal berdasarkan ID
	 */
	Find(id int) (entities.Comment, error)

	/*
	 * Find User
	 * -------------------------------
	 * Mencari comment berdasarkan field tertentu
	 */
	FindBy(field string, value string) (entities.Comment, error)

	/*
	 * CountAll
	 * -------------------------------
	 * Menghitung semua comments (ini digunakan untuk pagination di service)
	 */
	CountAll(filters []map[string]string) (int64, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan data comment kedalam database
	 */
	Store(comment entities.Comment) (entities.Comment, error)

	/*
	 * Update
	 * -------------------------------
	 * Mengupdate data comment berdasarkan ID
	 */
	Update(comment entities.Comment, id int) (entities.Comment, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete comment berdasarkan ID
	 */
	Delete(id int) error
}