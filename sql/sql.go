package sql

import (
	"context"
	"database/sql"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/uuumate/inf/logging"
)

const (
	MysqlDriver = "mysql"
)

const (
	_defaultMaxIdle        = 100
	_defaultMaxActive      = 100
	_defaultMaxLifetimeSec = 1800
)

var (
	SupportDBDriver = map[string]bool{
		MysqlDriver: true,
	}
)

type sqlGroup struct {
	name   string
	master *gorm.DB
	slaves []*gorm.DB
	next   uint64
}

type Config struct {
	Name   string   `json:"name" toml:"driver"`
	Driver string   `json:"driver" toml:"driver"`
	Master string   `json:"master" toml:"master"`
	Slaves []string `json:"slaves" toml:"slaves"`
}

func InitSqlGroup(cfg *Config) (*sqlGroup, error) {
	if !SupportDBDriver[cfg.Driver] {
		return nil, NotSupportSqlDriver
	}

	masterDB, err := OpenDB(cfg.Driver, cfg.Master)
	if err != nil {
		return nil, err
	}

	slaveDBs := make([]*gorm.DB, len(cfg.Slaves))
	for i := 0; i < len(cfg.Slaves); i++ {
		slaveDBs[i], err = OpenDB(cfg.Driver, cfg.Slaves[i])
		if err != nil {
			return nil, err
		}
	}

	return &sqlGroup{
		name:   cfg.Name,
		master: masterDB,
		slaves: slaveDBs,
		next:   0,
	}, nil
}

func OpenDB(driver, dsn string) (*gorm.DB, error) {
	addr, maxActive, maxIdle, maxLifetimeSec, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, addr)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(driver, db)
	if err != nil {
		return nil, err
	}

	gormDB = gormDB.Debug()

	gormDB.DB().SetMaxIdleConns(maxIdle)
	gormDB.DB().SetMaxOpenConns(maxActive)
	gormDB.DB().SetConnMaxLifetime(time.Second * time.Duration(maxLifetimeSec))
	gormDB.SetLogger(&logger{})

	return gormDB, nil
}

func parseDSN(dsn string) (string, int, int, int, error) {
	c, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", 0, 0, 0, err
	}

	maxActive, _ := strconv.ParseInt(c.Params["max_active"], 10, 64)
	maxIdle, _ := strconv.ParseInt(c.Params["max_idle"], 10, 64)
	maxLifetimeSec, _ := strconv.ParseInt(c.Params["max_lifetime_sec"], 10, 64)

	delete(c.Params, "max_active")
	delete(c.Params, "max_idle")
	delete(c.Params, "max_lifetime_sec")

	if maxActive <= 0 {
		maxActive = _defaultMaxActive
	}

	if maxIdle <= 0 {
		maxIdle = _defaultMaxIdle
	}

	if maxLifetimeSec <= 0 {
		maxLifetimeSec = _defaultMaxLifetimeSec
	}

	return c.FormatDSN(), int(maxActive), int(maxIdle), int(maxLifetimeSec), nil
}

func (sg *sqlGroup) Master(ctx context.Context) *gorm.DB {
	return sg.master
}

func (sg *sqlGroup) Slave(ctx context.Context) *gorm.DB {
	if len(sg.slaves) == 0 {
		return sg.master
	}

	atomic.AddUint64(&sg.next, 1)
	return sg.slaves[int(sg.next)%len(sg.slaves)]
}

type logger struct {
}

func (l *logger) Print(v ...interface{}) {
	logging.Info(sqlFormat(v...))
}
