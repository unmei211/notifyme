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

var _ Repository = (*postgresRepository)(nil)

func (r *postgresRepository) Update(inbox *MessageInbox) error {
	res := r.gorm.DB.Save(inbox)
	return res.Error
}

func (r *postgresRepository) FindInboxesForWorker(workerId int, workersCount int, page int, pageSize int) ([]*MessageInbox, error) {

	sql := `
		SELECT * FROM inbox i
		WHERE
		    mod(abs(hashtext(i.message_key)::bigint), ?) = ?
			AND
		    i.processed_at IS NULL 
		ORDER BY i.received_at
		LIMIT ? OFFSET ?
    `

	var result []*MessageInbox

	tx := r.gorm.DB.Raw(sql, workersCount, workerId, pageSize, (page-1)*pageSize).Scan(&result)

	return result, tx.Error
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
