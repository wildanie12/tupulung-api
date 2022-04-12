package category

import (
	"tupulung/entities"
	web "tupulung/entities/web"
	categoryRepository "tupulung/repositories/category"

	"github.com/jinzhu/copier"
)

type CategoryService struct {
	categoryRepo categoryRepository.CategoryRepositoryInterface
}

func NewCategoryService(repository categoryRepository.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{
		categoryRepo: repository,
	}
}

/*
 * --------------------------
 * Get List of category 
 * --------------------------
 */
func (service CategoryService) FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CategoryResponse, error) {

	offset := (page - 1) * limit

	categorysRes := []entities.CategoryResponse{}
	categorys, err := service.categoryRepo.FindAll(limit, offset, filters, sorts)
	copier.Copy(&categorysRes, &categorys)
	return categorysRes, err
}


/*
 * --------------------------
 * Load pagination data 
 * --------------------------
 */
func (service CategoryService) GetPagination(page, limit int, filters []map[string]string) (web.Pagination, error) {
	totalRows, err := service.categoryRepo.CountAll(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	totalPages :=  totalRows / int64(limit)

	return web.Pagination{
		Page: page,
		Limit: limit,
		TotalPages: int(totalPages),
	}, nil
}

/*
 * --------------------------
 * Get single category data based on ID
 * --------------------------
 */
func (service CategoryService) Find(id int) (entities.CategoryResponse, error) {
	
	category, err := service.categoryRepo.Find(id)
	categoryRes  := entities.CategoryResponse{}
	copier.Copy(&categoryRes, &category)

	return categoryRes, err
}


/*
 * --------------------------
 * Create category resource
 * --------------------------
 */
func (service CategoryService) Create(categoryRequest entities.CategoryRequest) (entities.CategoryResponse, error) {
	
	// convert request to domain entity
	category := entities.Category{}
	copier.Copy(&category, &categoryRequest)

	// Repository action
	category, err := service.categoryRepo.Store(category)
	if err != nil {
		return entities.CategoryResponse{}, err
	}

	// process domain entity to response
	categoryRes := entities.CategoryResponse{}
	copier.Copy(&categoryRes, &category)

	return categoryRes, nil
}


/*
 * --------------------------
 * Update category resource
 * --------------------------
 */
func (service CategoryService) Update(categoryRequest entities.CategoryRequest, id int) (entities.CategoryResponse, error) {

	// Find category
	category, err := service.categoryRepo.Find(id)
	if err != nil {
		return entities.CategoryResponse{}, web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}

	// Merge updated data request to domain entity
	copier.CopyWithOption(&category, &categoryRequest, copier.Option{ IgnoreEmpty: true, DeepCopy: true })

	// repository action
	category, err = service.categoryRepo.Update(category, id)

	// Convert category domain to category response
	categoryRes := entities.CategoryResponse{}
	copier.Copy(&categoryRes, &category)

	return categoryRes, err
}

/*
 * --------------------------
 * Delete resource data 
 * --------------------------
 */
func (service CategoryService) Delete(id int) error {
	// Find category
	_, err := service.categoryRepo.Find(id)
	if err != nil {
		return web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}
	
	// Copy request to found category
	err = service.categoryRepo.Delete(id)
	return err
}