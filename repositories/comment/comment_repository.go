package comment

import (
	"tupulung/entities"
	"tupulung/entities/web"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data comment berdasarkan filters dan sorts
 */
func (repo CommentRepository) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Comment, error) {
	comments := []entities.Comment{}
	builder := repo.db.Limit(limit).Offset(offset)

	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"] + " " + filter["operator"] + " ?", filter["value"])
	}

	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}

	// Operation
	tx := builder.Find(&comments)
	if tx.Error != nil {
		return []entities.Comment{}, web.WebError{Code: 500, Message: tx.Error.Error()} 
	}
	return comments, nil
}

/*
 * Find
 * -------------------------------
 * Mencari comment tunggal berdasarkan ID
 */
func (repo CommentRepository) Find(id int) (entities.Comment, error) {
	// Get user dari database
	user := entities.Comment{}
	tx := repo.db.Find(&user, id)
	if tx.Error != nil {

		// Return error dengan code 500 
		return entities.Comment{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {
		
		// Return error dengan code 400 jika tidak ditemukan
		return entities.Comment{}, web.WebError{Code: 400, Message: "cannot get user data with specified id"}
	}
	return user, nil
}

/*
 * Find User
 * -------------------------------
 * Mencari comment berdasarkan field tertentu
 */
func (repo CommentRepository) FindBy(field string, value string) (entities.Comment, error) {

	// Get user dari database
	user := entities.Comment{}
	tx := repo.db.Where(field + " = ?", value).Find(&user)
	if tx.Error != nil {

		// return kode 500 jika terjadi error
		return entities.Comment{}, web.WebError{ Code: 500, Message: tx.Error.Error() }
	} else if tx.RowsAffected <= 0 {

		// return kode 400 jika tidak ditemukan
		return entities.Comment{}, web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}
	return user, nil
}

/*
 * CountAll
 * -------------------------------
 * Menghitung semua comments (ini digunakan untuk pagination di service)
 */
func (repo CommentRepository) CountAll(filters []map[string]string) (int64, error) {
	
	var count int64
	builder := repo.db.Model(&entities.Event{})
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	tx := builder.Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, Message: tx.Error.Error()}
	}
	return count, nil
}

/*
 * Store
 * -------------------------------
 * Menambahkan data comment kedalam database
 */
func (repo CommentRepository) Store(comment entities.Comment) (entities.Comment, error) {
	tx := repo.db.Create(&comment)
	if tx.Error != nil {
		return entities.Comment{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return comment, nil
}

/*
 * Update
 * -------------------------------
 * Mengupdate data comment berdasarkan ID
 */
func (repo CommentRepository) Update(comment entities.Comment, id int) (entities.Comment, error) {
	tx := repo.db.Save(&comment)
	if tx.Error != nil {
		return entities.Comment{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return comment, nil
}

/*
 * Delete
 * -------------------------------
 * Delete comment berdasarkan ID
 */
func (repo CommentRepository) Delete(id int) error {
	tx := repo.db.Delete(&entities.Category{}, id)
	if tx.Error != nil {
		return web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return nil
}