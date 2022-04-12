package category

import "tupulung/entities"

type CategoryRepositoryInterface interface {
	/*
	 * Find All
	 * -------------------------------
	 * Mengambil data category berdasarkan filters dan sorts
	 */
	FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Category, error)

	/*
	 * Find
	 * -------------------------------
	 * Mencari category tunggal berdasarkan ID
	 */
	Find(id int) (entities.Category, error)

	/*
	 * CountAll
	 * -------------------------------
	 * Menghitung semua events (ini digunakan untuk pagination di service)
	 */
	CountAll(filters []map[string]string) (int64, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan data category kedalam database
	 */
	Store(category entities.Category) (entities.Category, error)

	/*
	 * Update
	 * -------------------------------
	 * Mengupdate data category berdasarkan ID
	 */
	Update(category entities.Category, id int) (entities.Category, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete category berdasarkan ID
	 */
	Delete(id int) error
}