package znakdb

import (
	"fmt"
)

// вызывается каждый раз при создании объекта
func (r *DbZnak) Check() (err error) {
	if !r.dbInfo.Exists {
		return fmt.Errorf("%s dbInfo.Exists false", modError)
	}
	r.dbSession, err = r.dbInfo.Connect()
	if err != nil {
		return fmt.Errorf("%s check ошибка подключения к БД %w", modError, err)
	}
	return nil
}
