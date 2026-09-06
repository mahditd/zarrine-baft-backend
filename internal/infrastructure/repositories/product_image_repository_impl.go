package repositories

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	domainRepositories "github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"

	"gorm.io/gorm"
)

type ProductImageRepositoryImpl struct {
	db *gorm.DB
}

func NewProductImageRepository(
	db *gorm.DB,
) domainRepositories.ProductImageRepository {

	return &ProductImageRepositoryImpl{
		db: db,
	}
}

func (r *ProductImageRepositoryImpl) Create(
	image *models.ProductImage,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		var count int64

		err := tx.Model(&models.ProductImage{}).
			Where("product_id = ?", image.ProductID).
			Count(&count).
			Error

		if err != nil {
			return err
		}

		// First image becomes cover; later uploads append at the end
		// so they never steal the cover (SRS 5.3: admin controls order).
		if count == 0 {
			image.DisplayOrder = 1
			image.IsCover = true
		} else {
			var maxOrder int
			err = tx.Model(&models.ProductImage{}).
				Where("product_id = ?", image.ProductID).
				Select("COALESCE(MAX(display_order), 0)").
				Scan(&maxOrder).
				Error

			if err != nil {
				return err
			}

			image.DisplayOrder = maxOrder + 1
			image.IsCover = false
		}

		return tx.Select(
			"ProductID",
			"ImageURL",
			"FilePath",
			"DisplayOrder",
			"IsCover",
		).Create(image).Error
	})
}

func (r *ProductImageRepositoryImpl) FindByID(
	id uint,
) (*models.ProductImage, error) {
	var image models.ProductImage
	err := r.db.First(&image, id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *ProductImageRepositoryImpl) FindByProductID(
	productID uint,
) ([]models.ProductImage, error) {

	var images []models.ProductImage

	err := r.db.
		Where("product_id = ?", productID).
		Order("display_order ASC, id ASC").
		Find(&images).
		Error

	return images, err
}

func (r *ProductImageRepositoryImpl) CountByProductID(
	productID uint,
) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

func (r *ProductImageRepositoryImpl) Update(
	image *models.ProductImage,
) error {
	return r.db.Save(image).Error
}

func (r *ProductImageRepositoryImpl) Reorder(
	productID uint,
	imageIDs []uint,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range imageIDs {
			isCover := (i == 0)
			err := tx.Model(&models.ProductImage{}).
				Where("id = ? AND product_id = ?", id, productID).
				Updates(map[string]interface{}{
					"display_order": i + 1,
					"is_cover":      isCover,
				}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProductImageRepositoryImpl) Delete(
	image *models.ProductImage,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		productID := image.ProductID
		wasCover := image.IsCover

		// Permanent deletion
		if err := tx.Unscoped().Delete(image).Error; err != nil {
			return err
		}

		// If the deleted image was cover, reassign cover to the new first image
		if wasCover {
			var first models.ProductImage
			err := tx.Where("product_id = ?", productID).
				Order("display_order ASC, id ASC").
				First(&first).
				Error
			if err == nil {

				if err := tx.Model(&models.ProductImage{}).
					Where("product_id = ?", productID).
					Update("is_cover", false).
					Error; err != nil {
					return err
				}

				if err := tx.Model(&first).
					Updates(map[string]interface{}{
						"is_cover":      true,
						"display_order": 1,
					}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
