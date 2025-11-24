package inbox

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/unmei211/notifyme/internal/pkg/orm"
	"go.uber.org/zap"
)

type postgresRepository struct {
	log  *zap.SugaredLogger
	cfg  *orm.DatabaseConfig
	gorm *orm.Gorm
}

func (r *postgresRepository) ExistsByMessageId(messageId uuid.UUID) (bool, error) {
	var messageInbox MessageInbox

	var count int64
	r.gorm.DB.Model(&messageInbox).Where("message_id = ?", messageId).Count(&count)

	return count > 0, nil
}

func (r *postgresRepository) Add(msg *MessageInbox) (*MessageInbox, error) {

	if err := r.gorm.DB.Create(&msg).Error; err != nil {
		return nil, errors.Wrap(err, "Can't add inbox message")
	}

	return msg, nil
}

func InitRepository(
	log *zap.SugaredLogger,
	cfg *orm.DatabaseConfig,
	gorm *orm.Gorm,
) Repository {
	return &postgresRepository{
		log:  log,
		cfg:  cfg,
		gorm: gorm,
	}
}
