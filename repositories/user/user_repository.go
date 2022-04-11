package user

import (
	entityDomain "tupulung/entities/domain"
	web "tupulung/entities/web"

	"gorm.io/gorm"
)
type UserRepository struct {
	db *gorm.DB
}

// Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

/*
 * Find User by ID
 * -------------------------------
 * Mencari user berdasarkan ID
 */
func (repo UserRepository) Find(id int) (entityDomain.User, error) {

	// Get user dari database
	user := entityDomain.User{}
	tx := repo.db.Find(&user, id)
	if tx.Error != nil {

		// Return error dengan code 500 
		return entityDomain.User{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {
		
		// Return error dengan code 400 jika tidak ditemukan
		return entityDomain.User{}, web.WebError{Code: 400, Message: "cannot get user data with specified id"}
	}
	return user, nil
}

/*
 * Find By Column
 * -------------------------------
 * Mencari user tunggal berdasarkan column dan value
 */
func (repo UserRepository) FindBy(field string, value string) (entityDomain.User, error) {

	// Get user dari database
	user := entityDomain.User{}
	tx := repo.db.Where(field + " = ?", value).Find(&user)
	if tx.Error != nil {

		// return kode 500 jika terjadi error
		return entityDomain.User{}, web.WebError{ Code: 500, Message: tx.Error.Error() }
	} else if tx.RowsAffected <= 0 {

		// return kode 400 jika tidak ditemukan
		return entityDomain.User{}, web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}
	return user, nil
}


/*
 * Store
 * -------------------------------
 * Menambahkan user tunggal kedalam database
 */
func (repo UserRepository) Store(user entityDomain.User) (entityDomain.User, error) {
	
	// insert user ke database
	tx := repo.db.Create(&user)
	if tx.Error != nil {

		// return kode 500 jika error
		return entityDomain.User{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return user, nil
}


/*
 * Update User
 * -------------------------------
 * Mengedit user tunggal berdasarkan ID
 */
func (repo UserRepository) Update(user entityDomain.User, id int) (entityDomain.User, error) {

	// Update database
	tx := repo.db.Save(&user)
	if tx.Error != nil {

		// return Kode 500 jika error
		return entityDomain.User{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return user, nil
}

/*
 * Delete
 * -------------------------------
 * Delete user tunggal berdasarkan ID
 */
func (repo UserRepository) Delete(id int) error {

	// Delete from database
	tx := repo.db.Delete(&entityDomain.User{}, id)
	if tx.Error != nil {

		// return kode 500 jika error
		return web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return nil
}
