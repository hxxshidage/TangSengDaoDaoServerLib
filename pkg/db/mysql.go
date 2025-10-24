package db

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/gocraft/dbr/v2"
	migrate "github.com/rubenv/sql-migrate"
)

type sqlLogger struct{}

func (s *sqlLogger) Event(eventName string) {
	log.Printf("[Event]:%s", eventName)
}

func (s *sqlLogger) EventKv(eventName string, kvs map[string]string) {
	log.Printf("[EventKv]:%s, kvs:%v", eventName, kvs)
}

func (s *sqlLogger) EventErr(eventName string, err error) error {
	log.Printf("[EventErr]:%s, error:%v", eventName, err)
	return err
}

func (s *sqlLogger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	log.Printf("[EventErrKv]:%s, error:%v, kvs: %v", eventName, err, kvs)
	return err
}

func (s *sqlLogger) Timing(eventName string, nanoseconds int64) {
	log.Printf("[Timing]:%s, duration:%d ns", eventName, nanoseconds)
}

func (s *sqlLogger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	log.Printf("[TimingKv]:%s, kvs:%v, duration: %d ns,", eventName, kvs, nanoseconds)
}

// NewMySQL 创建一个MySQL db，[path]db存储路径 [sqlDir]sql脚本目录
func NewMySQL(addr string, maxOpenConns int, maxIdleConns int, connMaxLifetime time.Duration, outputLog bool) *dbr.Session {
	var sqlLog dbr.EventReceiver
	if outputLog {
		sqlLog = &sqlLogger{}
	}

	conn, err := dbr.Open("mysql", addr, sqlLog)
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetConnMaxLifetime(connMaxLifetime) //mysql 默认超时时间为 60*60*8=28800 SetConnMaxLifetime设置为小于数据库超时时间即可

	session := conn.NewSession(nil)

	return session
}

func Migration(sqlDir string, session *dbr.Session) error {
	migrations := &FileDirMigrationSource{
		Dir: sqlDir,
	}
	_, err := migrate.Exec(session.DB, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}
	return nil
}

type byID []*migrate.Migration

func (b byID) Len() int           { return len(b) }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byID) Less(i, j int) bool { return b[i].Less(b[j]) }

// FileDirMigrationSource 文件目录源 遇到目录进行递归获取
type FileDirMigrationSource struct {
	Dir string
}

// FindMigrations FindMigrations
func (f FileDirMigrationSource) FindMigrations() ([]*migrate.Migration, error) {
	filesystem := http.Dir(f.Dir)
	migrations := make([]*migrate.Migration, 0, 100)
	err := f.findMigrations(filesystem, &migrations)
	if err != nil {
		return nil, err
	}
	// Make sure migrations are sorted
	sort.Sort(byID(migrations))

	return migrations, nil
}

func (f FileDirMigrationSource) findMigrations(dir http.FileSystem, migrations *[]*migrate.Migration) error {

	file, err := dir.Open("/")
	if err != nil {
		return err
	}

	files, err := file.Readdir(0)
	if err != nil {
		return err
	}

	for _, info := range files {

		if strings.HasSuffix(info.Name(), ".sql") {
			file, err := dir.Open(info.Name())
			if err != nil {
				return fmt.Errorf("Error while opening %s: %s", info.Name(), err)
			}

			migration, err := migrate.ParseMigration(info.Name(), file)
			if err != nil {
				return fmt.Errorf("Error while parsing %s: %s", info.Name(), err)
			}
			*migrations = append(*migrations, migration)

		} else if info.IsDir() {
			err = f.findMigrations(http.Dir(fmt.Sprintf("%s/%s", f.Dir, info.Name())), migrations)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
