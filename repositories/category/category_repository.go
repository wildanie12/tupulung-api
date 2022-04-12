package category

import (
	"tupulung/entities"
	web "tupulung/entities/web"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data category berdasarkan filters dan sorts
 */
func (repo CategoryRepository) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Category, error) {

	categories := []entities.Category{}
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
	tx := builder.Find(&categories)
	if tx.Error != nil {
		return []entities.Category{}, web.WebError{Code: 500, Message: tx.Error.Error()} 
	}
	return categories, nil
}

/*
 * CountAll
 * -------------------------------
 * Menghitung semua kategori (ini digunakan untuk pagination di service)
 */
func (repo CategoryRepository) CountAll(filters []map[string]string) (int64, error) {
	var count int64
	builder := repo.db.Model(&entities.Category{})
	
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"] + " " + filter["operator"] + " ?", filter["value"])
	}
	tx := builder.Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, Message: tx.Error.Error()}
	}
	return count, nil
}


/*
 * Find
 * -------------------------------
 * Mencari category tunggal berdasarkan ID
 */
func (repo CategoryRepository) Find(id int) (entities.Category, error) {
	category := entities.Category{}
	tx := repo.db.Find(&category, id)
	if tx.Error != nil {
		return entities.Category{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {
		return entities.Category{}, web.WebError{Code: 400, Message: "cannot get category data with specified id"}
	}
	return category, nil
}


/*
 * Store
 * -------------------------------
 * Menambahkan data category kedalam database
 */
func (repo CategoryRepository) Store(category entities.Category) (entities.Category, error) {
	
	tx := repo.db.Create(&category)
	if tx.Error != nil {
		return entities.Category{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return category, nil
}

/*
 * Update
 * -------------------------------
 * Mengupdate data category berdasarkan ID
 */
func (repo CategoryRepository) Update(category entities.Category, id int) (entities.Category, error) {
	tx := repo.db.Save(&category)
	if tx.Error != nil {
		return entities.Category{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return category, nil
}

/*
 * Delete
 * -------------------------------
 * Delete category berdasarkan ID
 */
func (repo CategoryRepository) Delete(id int) error {
	tx := repo.db.Delete(&entities.Category{}, id)
	if tx.Error != nil {
		return web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return nil
}
