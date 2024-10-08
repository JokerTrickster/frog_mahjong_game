// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"main/utils/db/mysql"
)

func newRoomUsers(db *gorm.DB, opts ...gen.DOOption) roomUsers {
	_roomUsers := roomUsers{}

	_roomUsers.roomUsersDo.UseDB(db, opts...)
	_roomUsers.roomUsersDo.UseModel(&mysql.RoomUsers{})

	tableName := _roomUsers.roomUsersDo.TableName()
	_roomUsers.ALL = field.NewAsterisk(tableName)
	_roomUsers.ID = field.NewUint(tableName, "id")
	_roomUsers.CreatedAt = field.NewTime(tableName, "created_at")
	_roomUsers.UpdatedAt = field.NewTime(tableName, "updated_at")
	_roomUsers.DeletedAt = field.NewField(tableName, "deleted_at")
	_roomUsers.UserID = field.NewInt(tableName, "user_id")
	_roomUsers.RoomID = field.NewInt(tableName, "room_id")
	_roomUsers.Score = field.NewInt(tableName, "score")
	_roomUsers.OwnedCardCount = field.NewInt(tableName, "owned_card_count")
	_roomUsers.PlayerState = field.NewString(tableName, "player_state")
	_roomUsers.TurnNumber = field.NewInt(tableName, "turn_number")

	_roomUsers.fillFieldMap()

	return _roomUsers
}

type roomUsers struct {
	roomUsersDo

	ALL            field.Asterisk
	ID             field.Uint
	CreatedAt      field.Time
	UpdatedAt      field.Time
	DeletedAt      field.Field
	UserID         field.Int
	RoomID         field.Int
	Score          field.Int
	OwnedCardCount field.Int
	PlayerState    field.String
	TurnNumber     field.Int

	fieldMap map[string]field.Expr
}

func (r roomUsers) Table(newTableName string) *roomUsers {
	r.roomUsersDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r roomUsers) As(alias string) *roomUsers {
	r.roomUsersDo.DO = *(r.roomUsersDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *roomUsers) updateTableName(table string) *roomUsers {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewUint(table, "id")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")
	r.DeletedAt = field.NewField(table, "deleted_at")
	r.UserID = field.NewInt(table, "user_id")
	r.RoomID = field.NewInt(table, "room_id")
	r.Score = field.NewInt(table, "score")
	r.OwnedCardCount = field.NewInt(table, "owned_card_count")
	r.PlayerState = field.NewString(table, "player_state")
	r.TurnNumber = field.NewInt(table, "turn_number")

	r.fillFieldMap()

	return r
}

func (r *roomUsers) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *roomUsers) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 10)
	r.fieldMap["id"] = r.ID
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
	r.fieldMap["deleted_at"] = r.DeletedAt
	r.fieldMap["user_id"] = r.UserID
	r.fieldMap["room_id"] = r.RoomID
	r.fieldMap["score"] = r.Score
	r.fieldMap["owned_card_count"] = r.OwnedCardCount
	r.fieldMap["player_state"] = r.PlayerState
	r.fieldMap["turn_number"] = r.TurnNumber
}

func (r roomUsers) clone(db *gorm.DB) roomUsers {
	r.roomUsersDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r roomUsers) replaceDB(db *gorm.DB) roomUsers {
	r.roomUsersDo.ReplaceDB(db)
	return r
}

type roomUsersDo struct{ gen.DO }

type IRoomUsersDo interface {
	gen.SubQuery
	Debug() IRoomUsersDo
	WithContext(ctx context.Context) IRoomUsersDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRoomUsersDo
	WriteDB() IRoomUsersDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRoomUsersDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRoomUsersDo
	Not(conds ...gen.Condition) IRoomUsersDo
	Or(conds ...gen.Condition) IRoomUsersDo
	Select(conds ...field.Expr) IRoomUsersDo
	Where(conds ...gen.Condition) IRoomUsersDo
	Order(conds ...field.Expr) IRoomUsersDo
	Distinct(cols ...field.Expr) IRoomUsersDo
	Omit(cols ...field.Expr) IRoomUsersDo
	Join(table schema.Tabler, on ...field.Expr) IRoomUsersDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRoomUsersDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRoomUsersDo
	Group(cols ...field.Expr) IRoomUsersDo
	Having(conds ...gen.Condition) IRoomUsersDo
	Limit(limit int) IRoomUsersDo
	Offset(offset int) IRoomUsersDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRoomUsersDo
	Unscoped() IRoomUsersDo
	Create(values ...*mysql.RoomUsers) error
	CreateInBatches(values []*mysql.RoomUsers, batchSize int) error
	Save(values ...*mysql.RoomUsers) error
	First() (*mysql.RoomUsers, error)
	Take() (*mysql.RoomUsers, error)
	Last() (*mysql.RoomUsers, error)
	Find() ([]*mysql.RoomUsers, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*mysql.RoomUsers, err error)
	FindInBatches(result *[]*mysql.RoomUsers, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*mysql.RoomUsers) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRoomUsersDo
	Assign(attrs ...field.AssignExpr) IRoomUsersDo
	Joins(fields ...field.RelationField) IRoomUsersDo
	Preload(fields ...field.RelationField) IRoomUsersDo
	FirstOrInit() (*mysql.RoomUsers, error)
	FirstOrCreate() (*mysql.RoomUsers, error)
	FindByPage(offset int, limit int) (result []*mysql.RoomUsers, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRoomUsersDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	GetByID(id int) (result mysql.RoomUsers, err error)
	GetByRoles(rolesName []string) (result []*mysql.RoomUsers, err error)
	InsertValue(name string, age int) (err error)
}

// SELECT * FROM @@table WHERE id=@id
func (r roomUsersDo) GetByID(id int) (result mysql.RoomUsers, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM room_users WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// GetByRoles query data by roles and return it as *slice of pointer*
//
//	(The below blank line is required to comment for the generated method)
//
// SELECT * FROM @@table WHERE role IN @rolesName
func (r roomUsersDo) GetByRoles(rolesName []string) (result []*mysql.RoomUsers, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, rolesName)
	generateSQL.WriteString("SELECT * FROM room_users WHERE role IN ? ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// InsertValue insert value
//
// INSERT INTO @@table (name, age) VALUES (@name, @age)
func (r roomUsersDo) InsertValue(name string, age int) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO room_users (name, age) VALUES (?, ?) ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (r roomUsersDo) Debug() IRoomUsersDo {
	return r.withDO(r.DO.Debug())
}

func (r roomUsersDo) WithContext(ctx context.Context) IRoomUsersDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r roomUsersDo) ReadDB() IRoomUsersDo {
	return r.Clauses(dbresolver.Read)
}

func (r roomUsersDo) WriteDB() IRoomUsersDo {
	return r.Clauses(dbresolver.Write)
}

func (r roomUsersDo) Session(config *gorm.Session) IRoomUsersDo {
	return r.withDO(r.DO.Session(config))
}

func (r roomUsersDo) Clauses(conds ...clause.Expression) IRoomUsersDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r roomUsersDo) Returning(value interface{}, columns ...string) IRoomUsersDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r roomUsersDo) Not(conds ...gen.Condition) IRoomUsersDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r roomUsersDo) Or(conds ...gen.Condition) IRoomUsersDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r roomUsersDo) Select(conds ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r roomUsersDo) Where(conds ...gen.Condition) IRoomUsersDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r roomUsersDo) Order(conds ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r roomUsersDo) Distinct(cols ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r roomUsersDo) Omit(cols ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r roomUsersDo) Join(table schema.Tabler, on ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r roomUsersDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r roomUsersDo) RightJoin(table schema.Tabler, on ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r roomUsersDo) Group(cols ...field.Expr) IRoomUsersDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r roomUsersDo) Having(conds ...gen.Condition) IRoomUsersDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r roomUsersDo) Limit(limit int) IRoomUsersDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r roomUsersDo) Offset(offset int) IRoomUsersDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r roomUsersDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRoomUsersDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r roomUsersDo) Unscoped() IRoomUsersDo {
	return r.withDO(r.DO.Unscoped())
}

func (r roomUsersDo) Create(values ...*mysql.RoomUsers) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r roomUsersDo) CreateInBatches(values []*mysql.RoomUsers, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r roomUsersDo) Save(values ...*mysql.RoomUsers) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r roomUsersDo) First() (*mysql.RoomUsers, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*mysql.RoomUsers), nil
	}
}

func (r roomUsersDo) Take() (*mysql.RoomUsers, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*mysql.RoomUsers), nil
	}
}

func (r roomUsersDo) Last() (*mysql.RoomUsers, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*mysql.RoomUsers), nil
	}
}

func (r roomUsersDo) Find() ([]*mysql.RoomUsers, error) {
	result, err := r.DO.Find()
	return result.([]*mysql.RoomUsers), err
}

func (r roomUsersDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*mysql.RoomUsers, err error) {
	buf := make([]*mysql.RoomUsers, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r roomUsersDo) FindInBatches(result *[]*mysql.RoomUsers, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r roomUsersDo) Attrs(attrs ...field.AssignExpr) IRoomUsersDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r roomUsersDo) Assign(attrs ...field.AssignExpr) IRoomUsersDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r roomUsersDo) Joins(fields ...field.RelationField) IRoomUsersDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r roomUsersDo) Preload(fields ...field.RelationField) IRoomUsersDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r roomUsersDo) FirstOrInit() (*mysql.RoomUsers, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*mysql.RoomUsers), nil
	}
}

func (r roomUsersDo) FirstOrCreate() (*mysql.RoomUsers, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*mysql.RoomUsers), nil
	}
}

func (r roomUsersDo) FindByPage(offset int, limit int) (result []*mysql.RoomUsers, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r roomUsersDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r roomUsersDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r roomUsersDo) Delete(models ...*mysql.RoomUsers) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *roomUsersDo) withDO(do gen.Dao) *roomUsersDo {
	r.DO = *do.(*gen.DO)
	return r
}
