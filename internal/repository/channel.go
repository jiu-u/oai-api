package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiu-u/oai-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryOption func(*gorm.DB) *gorm.DB

func WithChannelId(id uint64) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func WithChannelName(name string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func WithChannelStatus(status string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func WithChannelIds(ids []uint64) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}

func WithChannelStatuses(statuses []string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status in ?", statuses)
	}
}

func WithPage(page, pageSize int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
}

func WithHashId(hashId string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("hash_id = ?", hashId)
	}
}

type ChannelRepository interface {
	CreateChannel(ctx context.Context, channel *model.Channel) error
	CreateChannelBatch(ctx context.Context, channels []*model.Channel) error
	CreateChannelIfNotExists(ctx context.Context, channel *model.Channel) error
	FindChannelById(ctx context.Context, id uint64) (*model.Channel, error)
	FindChannelByIdForUpdate(ctx context.Context, id uint64) (*model.Channel, error)
	FindChannelByIdForShare(ctx context.Context, id uint64) (*model.Channel, error)
	FindAllChannels(ctx context.Context) ([]*model.Channel, error)
	FindAllChannelsByCondition(ctx context.Context, options ...QueryOption) ([]*model.Channel, error)
	ExistsChannel(ctx context.Context, options ...QueryOption) (bool, error)
	UpdateChannel(ctx context.Context, channel *model.Channel) error
	DeleteChannelByID(ctx context.Context, id uint64) error
	//FindByCondition(ctx context.Context, options ...QueryOption) (*model.Channel, error)
	//Count(ctx context.Context, condition map[string]interface{}) (int64, error)
	//UpdateByCondition(ctx context.Context, condition map[string]interface{}, channel *model.Channel) error
	//UpdateSelective(ctx context.Context, condition map[string]interface{}, channel *model.Channel) error
	//LoadBatch(ctx context.Context, channels []*model.Channel) error
	//DeleteByCondition(ctx context.Context, condition map[string]interface{}) error
	//DeleteAll(ctx context.Context) error
}

func NewChannelRepository(repo *Repository) ChannelRepository {
	return &channelRepository{repo}
}

type channelRepository struct {
	*Repository
}

func (r *channelRepository) CreateChannel(ctx context.Context, channel *model.Channel) error {
	return r.DB(ctx).Create(channel).Error
}

func (r *channelRepository) CreateChannelBatch(ctx context.Context, channels []*model.Channel) error {
	return r.DB(ctx).CreateInBatches(channels, 100).Error
}

func (r *channelRepository) CreateChannelIfNotExists(ctx context.Context, channel *model.Channel) error {
	var temp model.Channel
	row := r.DB(ctx).Where("hash_id = ?", channel.HashId).First(&temp)
	if row.Error != nil {
		if errors.Is(row.Error, gorm.ErrRecordNotFound) {
			err := r.DB(ctx).Create(channel).Error
			return err
		}
		return row.Error
	}
	channel.Id = temp.Id
	return nil
}

func (r *channelRepository) FindChannelById(ctx context.Context, id uint64) (*model.Channel, error) {
	var channel model.Channel
	err := r.DB(ctx).First(&channel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}
	return &channel, nil
}

func (r *channelRepository) FindChannelByIdForUpdate(ctx context.Context, id uint64) (*model.Channel, error) {
	var channel model.Channel
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&channel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}
	return &channel, nil
}

func (r *channelRepository) FindChannelByIdForShare(ctx context.Context, id uint64) (*model.Channel, error) {
	var channel model.Channel
	err := r.DB(ctx).Clauses(clause.Locking{Strength: "SHARE"}).First(&channel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("channel with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}
	return &channel, nil
}

func (r *channelRepository) FindAllChannels(ctx context.Context) ([]*model.Channel, error) {
	var channels []*model.Channel
	err := r.DB(ctx).Find(&channels).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching channels: %w", err)
	}
	return channels, nil
}

func (r *channelRepository) FindAllChannelsByCondition(ctx context.Context, options ...QueryOption) ([]*model.Channel, error) {
	var channels []*model.Channel
	query := r.DB(ctx).Model(&model.Channel{})
	for _, option := range options {
		query = option(query)
	}
	err := r.DB(ctx).Find(&channels).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching channels: %w", err)
	}
	return channels, nil
}

func (r *channelRepository) ExistsChannel(ctx context.Context, options ...QueryOption) (bool, error) {
	var count int64
	query := r.DB(ctx).Model(&model.Channel{})
	for _, option := range options {
		query = option(query)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error counting channels: %w", err)
	}
	return count > 0, nil
}

func (r *channelRepository) UpdateChannel(ctx context.Context, channel *model.Channel) error {
	return r.DB(ctx).Updates(channel).Error
}

//func (r *channelRepository) UpdateByCondition(ctx context.Context, condition map[string]interface{}, channel *model.Channel) error {
//	return r.DB(ctx).Where(condition).Updates(channel).Error
//}
//
//func (r *channelRepository) UpdateSelective(ctx context.Context, condition map[string]interface{}, channel *model.Channel) error {
//	return r.DB(ctx).Where(condition).Select("*").Updates(channel).Error
//}

//// LoadBatch `ON DUPLICATE KEY UPDATE` 存在就更新，不存在就插入 !!!更新所有字段
//func (r *channelRepository) LoadBatch(ctx context.Context, channels []*model.Channel) error {
//	return r.DB(ctx).Save(channels).Error
//}

func (r *channelRepository) DeleteChannelByID(ctx context.Context, id uint64) error {
	return r.DB(ctx).Delete(&model.Channel{}, id).Error
}

//func (r *channelRepository) DeleteByCondition(ctx context.Context, condition map[string]interface{}) error {
//	return r.DB(ctx).Where(condition).Delete(&model.Channel{}).Error
//}
//
//func (r *channelRepository) DeleteAll(ctx context.Context) error {
//	return r.DB(ctx).Delete(&model.Channel{}).Error
//}
