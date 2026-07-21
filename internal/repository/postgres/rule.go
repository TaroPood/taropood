package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TaroPood/taropood/internal/domain"
	"github.com/TaroPood/taropood/internal/repository/postgres/model"
	"github.com/TaroPood/taropood/internal/repository/postgres/query"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type RuleRepository struct {
	db *gorm.DB
	q  *query.Query
}

func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{
		db: db,
		q:  query.Use(db),
	}
}

func (r *RuleRepository) Create(ctx context.Context, rule *domain.Rule) error {
	if rule.ID == "" {
		rule.ID = newID()
	}
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now

	for i := range rule.Actions {
		if rule.Actions[i].ID == "" {
			rule.Actions[i].ID = newID()
		}
		rule.Actions[i].RuleID = rule.ID
		rule.Actions[i].CreatedAt = now
	}

	m, err := model.RuleToModel(rule)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInvalidRule, err)
	}

	if err := r.q.RuleModel.WithContext(ctx).Create(m); err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %s", domain.ErrDuplicateName, rule.Name)
		}
		return fmt.Errorf("create rule: %w", err)
	}

	*rule = *m.ToDomain()
	return nil
}

func (r *RuleRepository) GetByID(ctx context.Context, id string) (*domain.Rule, error) {
	var m model.RuleModel
	err := r.db.WithContext(ctx).
		Preload("Actions", func(db *gorm.DB) *gorm.DB {
			return db.Order(`"order" ASC`)
		}).
		First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("rule %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get rule %s: %w", id, err)
	}

	return m.ToDomain(), nil
}

func (r *RuleRepository) List(ctx context.Context, filter domain.RuleFilter) ([]*domain.Rule, error) {
	q := r.q.RuleModel.WithContext(ctx).
		Order(r.q.RuleModel.Priority.Desc(), r.q.RuleModel.CreatedAt.Asc())

	if filter.Name != nil {
		q = q.Where(r.q.RuleModel.Name.Eq(*filter.Name))
	}
	if filter.Enabled != nil {
		q = q.Where(r.q.RuleModel.Enabled.Eq(*filter.Enabled))
	}
	if len(filter.Tags) > 0 {
		q = q.Where(field.NewUnsafeFieldRaw("tags ?| array[?]", filter.Tags))
	}
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit)
	} else {
		q = q.Limit(100)
	}
	if filter.Offset > 0 {
		q = q.Offset(filter.Offset)
	}

	models, err := q.Find()
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}

	rules := make([]*domain.Rule, 0, len(models))
	for i := range models {
		rules = append(rules, models[i].ToDomain())
	}
	return rules, nil
}

func (r *RuleRepository) Update(ctx context.Context, rule *domain.Rule) error {
	rule.UpdatedAt = time.Now().UTC()

	m, err := model.RuleToModel(rule)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInvalidRule, err)
	}

	err = r.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.ActionModel.WithContext(ctx).
			Where(tx.ActionModel.RuleID.Eq(rule.ID)).
			Delete(); err != nil {
			return err
		}
		for i := range m.Actions {
			m.Actions[i].RuleID = rule.ID
			if m.Actions[i].ID == "" {
				m.Actions[i].ID = newID()
			}
			m.Actions[i].CreatedAt = rule.CreatedAt
		}
		return tx.UnderlyingDB().Save(m).Error
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %s", domain.ErrDuplicateName, rule.Name)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("rule %s: %w", rule.ID, domain.ErrNotFound)
		}
		return fmt.Errorf("update rule %s: %w", rule.ID, err)
	}

	return nil
}

func (r *RuleRepository) Delete(ctx context.Context, id string) error {
	info, err := r.q.RuleModel.WithContext(ctx).
		Where(r.q.RuleModel.ID.Eq(id)).
		Delete()
	if err != nil {
		return fmt.Errorf("delete rule %s: %w", id, err)
	}
	if info.RowsAffected == 0 {
		return fmt.Errorf("rule %s: %w", id, domain.ErrNotFound)
	}
	return nil
}

func (r *RuleRepository) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.q.Transaction(func(tx *query.Query) error {
		return fn(context.WithValue(ctx, txKey{}, tx.UnderlyingDB()))
	})
}

type txKey struct{}

func (r *RuleRepository) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

func newID() string {
	var buf [16]byte
	_, _ = rand.Read(buf[:])
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
