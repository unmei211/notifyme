package inbox

import (
	"github.com/unmei211/notifyme/internal/pkg/orm"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresRepository struct {
	log  *zap.SugaredLogger
	cfg  *orm.DatabaseConfig
	gorm *gorm.DB
}

func (r *postgresRepository) ExistsByMessageId(messageId string) (bool, error) {
	var messageInbox MessageInbox

	var count int64
	r.gorm.Model(&messageInbox).Where("message_id = ?", messageId).Count(&count)

	return count > 0, nil
}

func (r *postgresRepository)

func InitRepository(
	log *zap.SugaredLogger,
	cfg *orm.DatabaseConfig,
	gorm *gorm.DB,
) Repository {
	return &postgresRepository{
		log:  log,
		cfg:  cfg,
		gorm: gorm,
	}
}
