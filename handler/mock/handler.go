package mock

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"reflect"
)

var (
	ErrParamIsEmpty = errors.New("param is empty")
)

// Mock for repository, make sure all is correct
type RepositoryMock struct{}

func (s *RepositoryMock) FindByAlias(alias string, result interface{}) error {

	if alias == "" {
		return ErrParamIsEmpty
	}

	if result == nil {
		return ErrParamIsEmpty
	}

	mockUpdateID(result)

	return nil
}

func (s *RepositoryMock) ExistByAlias(alias string) bool {
	if alias != "" {
		return true
	}
	return false
}

func (s *RepositoryMock) Create(result interface{}) error {
	if result == nil {
		return ErrParamIsEmpty
	}
	return nil
}

func (s *RepositoryMock) Update(id string, result interface{}) error {
	if id == "" {
		return ErrParamIsEmpty
	}
	if result == nil {
		return ErrParamIsEmpty
	}
	return nil
}

func (s *RepositoryMock) MoreVisited(result interface{}) error {
	if result == nil {
		return ErrParamIsEmpty
	}
	return nil
}

func mockUpdateID(v interface{}) {
	indirectStruct := reflect.Indirect(reflect.ValueOf(v))
	idField := indirectStruct.FieldByName("ID")
	if idField.IsValid() {
		if idField.CanSet() {
			objectId := bson.NewObjectId()
			idField.Set(reflect.ValueOf(objectId))
		}
	}
}
