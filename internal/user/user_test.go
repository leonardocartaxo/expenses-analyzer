package user_test

//import (
//	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
//	"reflect"
//	"testing"
//)
//
//func mockDB() *map[string]user.DTO {
//	usersArr := []user.DTO{
//		*user.New("John Doe"),
//		*user.New("Alice"),
//		*user.New("Bob"),
//	}
//	usersMap := map[string]user.DTO{}
//	for _, it := range usersArr {
//		usersMap[it.ID] = it
//	}
//
//	return &usersMap
//}
//
//func TestRepository_Save(t *testing.T) {
//	repo := user.NewRepository(&map[string]user.DTO{})
//	newUser := *user.New("Charlie")
//	savedUser := repo.Save(newUser)
//
//	if !reflect.DeepEqual(savedUser, newUser) {
//		t.Errorf("saved user does not match saved user")
//	}
//}
//
//func TestRepository_FindOne(t *testing.T) {
//	repo := user.NewRepository(mockDB())
//	repo.FindOne("Bob")
//}
//
//func TestRepository_FindAll(t *testing.T) {
//	mockUsers := mockDB()
//	repo := user.NewRepository(mockUsers)
//	repoUsers := repo.All()
//
//	//if !cmp.Equal(mockUsers, repoUsers) {
//	if !reflect.DeepEqual(repoUsers, mockUsers) {
//		t.Errorf("All users do not match")
//	}
//}
