package repository

//type ProviderRepo interface {
//	InsertOne(ctx context.Context, provider *model.Provider) error
//	DeleteOne(ctx context.Context, id uint64) error
//	FindOne(ctx context.Context, id uint64) (*model.Provider, error)
//	FindAll(ctx context.Context) ([]*model.Provider, error)
//	UpdateOne(ctx context.Context, provider *model.Provider) error
//	ExistsHashId(ctx context.Context, hashId string) (bool, error)
//}
//
//func NewProviderRepo(r *Repository) ProviderRepo {
//	return &providerRepo{r}
//}
//
//type providerRepo struct {
//	*Repository
//}
//
//func (r *providerRepo) InsertOne(ctx context.Context, provider *model.Provider) error {
//	return r.DB(ctx).Create(provider).Error
//}
//
//func (r *providerRepo) DeleteOne(ctx context.Context, id uint64) error {
//	return r.DB(ctx).Model(&model.Provider{}).Where("id = ?", id).Delete(&model.Provider{}).Error
//}
//
//func (r *providerRepo) FindOne(ctx context.Context, id uint64) (*model.Provider, error) {
//	var provider model.Provider
//	err := r.DB(ctx).First(&provider, id).Error
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, fmt.Errorf("provider with ID %d not found", id)
//		}
//		return nil, fmt.Errorf("error fetching provider: %w", err)
//	}
//	return &provider, nil
//}
//
//func (r *providerRepo) FindAll(ctx context.Context) ([]*model.Provider, error) {
//	var providers []*model.Provider
//	err := r.DB(ctx).Find(&providers).Error
//	if err != nil {
//		return nil, fmt.Errorf("error fetching providers: %w", err)
//	}
//	return providers, nil
//}
//
//func (r *providerRepo) UpdateOne(ctx context.Context, provider *model.Provider) error {
//	err := r.DB(ctx).Updates(provider).Error
//	if err != nil {
//		return fmt.Errorf("error updating provider: %w", err)
//	}
//	return nil
//}
//
//func (r *providerRepo) ExistsHashId(ctx context.Context, hashId string) (bool, error) {
//	var provider model.Provider
//	err := r.DB(ctx).Where("hash_id = ?", hashId).First(&provider).Error
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return false, nil
//		}
//		return false, fmt.Errorf("error fetching provider: %w", err)
//	}
//	return true, nil
//}
