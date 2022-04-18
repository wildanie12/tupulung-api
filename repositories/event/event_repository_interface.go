package event

import "tupulung/entities"

type EventRepositoryInterface interface {
	/*
	 * Find All
	 * -------------------------------
	 * Mengambil data event berdasarkan filters dan sorts
	 */
	FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Event, error)

	/*
	 * Find
	 * -------------------------------
	 * Mencari event tunggal berdasarkan ID
	 */
	Find(id int) (entities.Event, error)

	/*
	 * Find User
	 * -------------------------------
	 * Mencari event berdasarkan field tertentu
	 */
	FindBy(field string, value string) (entities.Event, error)

	/*
	 * CountAll
	 * -------------------------------
	 * Menghitung semua events (ini digunakan untuk pagination di service)
	 */
	CountAll(filters []map[string]string) (int64, error)

	/*
	 * Store
	 * -------------------------------
	 * Menambahkan data event kedalam database
	 */
	Store(event entities.Event) (entities.Event, error)

	/*
	 * Update
	 * -------------------------------
	 * Mengupdate data event berdasarkan ID
	 */
	Update(event entities.Event, id int) (entities.Event, error)

	/*
	 * Delete
	 * -------------------------------
	 * Delete event berdasarkan ID
	 */
	Delete(id int) error

	/*
	 * Delete Batch
	 * -------------------------------
	 * Delete multiple event berdasarkan filter tertentu
	 */
	DeleteBatch(filters []map[string]string) error
}
