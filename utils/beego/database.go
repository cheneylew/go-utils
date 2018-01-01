package beego

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/orm"
	"github.com/cheneylew/goutil/utils"
	"database/sql"

)

const (
	DataTypeInt = "INT(11)"
	DataTypeBigInt = "BIGINT"
	DataTypeVarChar = "VARCHAR(255)"
	DataTypeVarCharSize = "VARCHAR(%d)"
	DataTypeDateTime = "DATETIME"
	DataTypeDate = "DATE"

	DataTypeUnsigned = "UNSIGNED"
	DataTypeNotNull = "NOT NULL"

)

func DBUrl(user, password, host, port, dbName string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8`, user, password, host, port, dbName)
}

func InitRegistDB(user,pwd,host,port,dbname string) *BaseDataBase {
	url := DBUrl(user,pwd,host,port,dbname)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", url)

	if err != nil {
		utils.JJKPrintln("========database can't connect! error:" + err.Error()+"========")
	} else {
		utils.JJKPrintln("========database connected success！========")
	}

	//创建模型表结构
	orm.RunSyncdb("default",false,true)

	return &BaseDataBase {
		Orm:orm.NewOrm(),
	}
}

type BaseDataBase struct {
	Orm orm.Ormer
}

func (db *BaseDataBase)DBBaseTableCount(tablename string) int64 {
	a, err := db.Orm.QueryTable(tablename).Count()
	if err != nil {
		return 0
	}

	return a
}

func (db *BaseDataBase)DBBaseAnyTableSelect(tablename string, start int, count int) ([]orm.Params, error) {
	var list []orm.Params
	sql := fmt.Sprintf("SELECT * FROM %s limit %d,%d", tablename, start, count)
	num, err := db.Orm.Raw(sql).Values(&list)
	if err == nil && num > 0 {
		return list, nil
	}

	return nil, err
}

func (db *BaseDataBase)DBBaseAnyTableSelectOneRow(tablename string, id int) (orm.Params, error) {
	var list []orm.Params
	sql := fmt.Sprintf("SELECT * FROM %s where %s_id=%d", tablename, tablename, id)
	num, err := db.Orm.Raw(sql).Values(&list)
	if err == nil && num > 0 {
		return list[0], nil
	}

	return nil, err
}

func (db *BaseDataBase)DBBaseAnyTableSelectOneRowWithContentID(tablename string, content_id int64) (orm.Params, error) {
	var list []orm.Params
	sql := fmt.Sprintf("SELECT * FROM %s where %s=%d", tablename, "content_id", content_id)
	num, err := db.Orm.Raw(sql).Values(&list)
	if err == nil && num > 0 {
		return list[0], nil
	}

	return nil, err
}

func (db *BaseDataBase)DBBaseCreateTable(tableName string) error {
	sql := fmt.Sprintf("CREATE TABLE `%s` (`%s_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,PRIMARY KEY (`%s_id`));", tableName, tableName, tableName)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseCreateTableWithContentID(tableName string) error {
	sql := fmt.Sprintf("CREATE TABLE `%s` (`%s_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,`content_id` INT UNSIGNED NOT NULL,PRIMARY KEY (`%s_id`));", tableName, tableName, tableName)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseDropTable(tableName string) error {
	sql := fmt.Sprintf("DROP TABLE `%s`;", tableName)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseUpdateColumn(tableName, oldColumnName, newColumnName, dataType string) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `%s` `%s` %s NOT NULL DEFAULT `` ;", tableName, oldColumnName, newColumnName, dataType)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseDropColumn(tableName, columnName string) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", tableName, columnName)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseAddColumn(tableName, columnName, dataType string) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s NOT NULL DEFAULT `` ;", tableName, columnName, dataType)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseAddColumnTinyInt(tableName, columnName string) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` TINYINT UNSIGNED NOT NULL DEFAULT 0 ;", tableName, columnName)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseAddColumnVarChar255(tableName, columnName string) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` VARCHAR(255) NOT NULL;", tableName, columnName)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseAddColumnVarChar(tableName, columnName string, size int64) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` VARCHAR(%d) NOT NULL;", tableName, columnName, size)
	_, err := db.DBBaseExecRawSQL(sql)
	if err != nil {
		return err
	}

	return nil
}

func (db *BaseDataBase)DBBaseExecRawSQL(asql string) (int64, error) {
	return db.DBBaseExecSQL(asql)
}
// insert 		返回num代表插入的ID
// update, delete 	返回num代表影响的行数
// create table		返回num为0
// select		不要使用这个语句
func (db *BaseDataBase)DBBaseExecSQL(asql string, params ...interface{}) (int64, error) {
	var res sql.Result
	var err error
	if len(params) > 0 {
		res, err = db.Orm.Raw(asql, params...).Exec()
	} else {
		res, err = db.Orm.Raw(asql).Exec()
	}

	if err != nil {
		return 0, err
	}

	rowid, _ := res.LastInsertId()
	if rowid > 0 {
		return rowid, nil
	}

	num, aerr := res.RowsAffected()
	if aerr != nil {
		return 0, aerr
	}

	if num > 0 {
		return num, nil
	}

	return 0, nil
}



