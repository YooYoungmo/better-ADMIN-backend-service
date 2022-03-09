package webhook

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type webHookRepository struct {
}

func (webHookRepository) Create(ctx context.Context, entity *WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return errors.New(err)
	}

	return nil
}

func (webHookRepository) FindAll(ctx context.Context, pageable dtos.Pageable) ([]WebHookEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&WebHookEntity{})

	var entities = make([]WebHookEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, errors.New(err)
	}

	return entities, totalCount, nil
}

func (webHookRepository) FindById(ctx context.Context, id uint) (WebHookEntity, error) {
	var entity WebHookEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.New(err)
	}

	return entity, nil
}

func (webHookRepository) Delete(ctx context.Context, entity WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return errors.New(err)
	}

	if err := db.Delete(&entity).Error; err != nil {
		return errors.New(err)
	}

	return nil
}

func (webHookRepository) Save(ctx context.Context, entity WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Save(entity).Error; err != nil {
		return errors.New(err)
	}

	return nil
}

func (webHookRepository) FindLast(ctx context.Context) (WebHookEntity, error) {
	var entity WebHookEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.New(err)
	}

	return entity, nil
}
