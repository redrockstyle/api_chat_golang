package dbtestsuser

import (
	"api_chat/api/layers/controller/database"
	"api_chat/api/layers/domain/db"
	"fmt"
	"testing"
)

type TestsUser struct {
	dbc *database.DbController
	t   *testing.T
}

func NewTestsUser(dbc *database.DbController, t *testing.T) *TestsUser {
	return &TestsUser{dbc: dbc, t: t}
}

func (tku *TestsUser) TestsAll() {
	tku.TestCreateUser()
	tku.TestGetUser()
	tku.TestCountUsers()
	tku.TestModifyUser()
	tku.TestDeleteUser()
}

func (tku *TestsUser) TestCreateUser() {
	user := db.User{
		Id:        "1a2b3c4d5e6f",
		FirstName: "fname",
		LastName:  "lname",
		Login:     "somelogin",
		Password:  "somepass",
	}
	user2 := db.User{
		Id:        "1a2b3c4d5e",
		FirstName: "fname",
		LastName:  "lname",
		Login:     "comelogin",
		Password:  "somepass",
	}
	tku.t.Run(fmt.Sprintf("craete users id1:%s id2:%s", user.Id, user2.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeCreate("", &user)
		if err != nil {
			tku.t.Errorf("error create user: %v", err)
		} else if id != user.Id {
			tku.t.Errorf("error create user: error returned id '%v' != '%v'", id, user.Id)
		}
		id, err = tku.dbc.TypeCreate("", &user2)
		if err != nil {
			tku.t.Errorf("error create user: %v", err)
		} else if id != user2.Id {
			tku.t.Errorf("error create user: error returned id '%v' != '%v'", id, user2.Id)
		}
	})
}

func (tku *TestsUser) TestGetUser() {
	user := db.User{
		Id: "1a2b3c4d5e6f",
	}
	// usersCheck := []db.User{
	// 	{
	// 		Id: "1a2b3c4d5e6f",
	// 	},
	// 	{
	// 		Id: "1a2b3c4d5e",
	// 	},
	// }
	// usersCheck2 := []db.User{
	// 	{
	// 		Id: "1a2b3c4d5e",
	// 	},
	// 	{
	// 		Id: "1a2b3c4d5e6f",
	// 	},
	// }

	tku.t.Run(fmt.Sprintf("get user by id:%v", user.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeGet("", &user, nil)
		if err != nil {
			tku.t.Errorf("error get one user: %v", err)
		} else if id != user.Id {
			tku.t.Errorf("error get one user: error returned id '%v' != '%v'", id, user.Id)
		}

		users := make([]db.User, 2)
		_, err = tku.dbc.TypeGets("", &users, map[string]interface{}{"first_name": "fname"}, 0, 0)
		if err != nil {
			tku.t.Errorf("error get users (conds): %v", err)
		}
		// else {
		// 	if len == 0 {
		// 		tku.t.Error("error get users (conds): len()=0")
		// 	}
		// 	var count int
		// 	for i := 0; i < int(len); i++ {
		// 		if usersCheck[i].Id != users[i].Id {
		// 			count += 1
		// 		}
		// 	}
		// 	if count != 0 {
		// 		if count != int(len) {
		// 			tku.t.Errorf("error get users (conds): error returned id '%v' != '%v' len()=%v count=%v", usersCheck[0].Id, users[0].Id, len, count)

		// 		} else {
		// 			for i := 0; i < int(len); i++ {
		// 				if usersCheck2[i].Id != users[i].Id {
		// 					tku.t.Errorf("error get users (conds): error returned id '%v' != '%v' len()=%v", usersCheck2[i].Id, users[i].Id, len)

		// 				}
		// 			}
		// 		}
		// 	}

		// }

		// users2 := make([]db.User, 2)
		// len, err = tku.dbc.TypeGets("", &users2, nil, []string{"1a2b3c4d5e6f", "1a2b3c4d5e"})
		// if err != nil {
		// 	tku.t.Errorf("error get users (ids): %v", err)
		// } else {
		// 	if len == 0 {
		// 		tku.t.Error("error get users (ids): len()=0")
		// 	}
		// 	var count int
		// 	for i := 0; i < int(len); i++ {
		// 		if usersCheck[i].Id != users[i].Id {
		// 			count += 1
		// 		}
		// 	}
		// 	if count != 0 {
		// 		if count != int(len) {
		// 			tku.t.Errorf("error get users (ids): error returned id '%v' != '%v' len()=%v count=%v", usersCheck[0].Id, users[0].Id, len, count)

		// 		} else {
		// 			for i := 0; i < int(len); i++ {
		// 				if usersCheck2[i].Id != users[i].Id {
		// 					tku.t.Errorf("error get users (ids): error returned id '%v' != '%v' len()=%v", usersCheck2[i].Id, users[i].Id, len)

		// 				}
		// 			}
		// 		}
		// 	}

		// }
	})
}

func (tku *TestsUser) TestCountUsers() {
	tku.t.Run("get count users", func(t *testing.T) {
		_, err := tku.dbc.TypeCount(0, &db.User{}, nil)
		if err != nil {
			tku.t.Errorf("error get count users: %v", err)
		}
		// if count != 2 {
		// 	tku.t.Errorf("error get count:%v users", count)
		// }
	})
}

func (tku *TestsUser) TestModifyUser() {
	user := db.User{
		Id:        "1a2b3c4d5e6f",
		FirstName: "fname2",
		LastName:  "lname2",
		Login:     "somelogin2",
		Password:  "somepass2",
	}
	userGet := db.User{
		Id: "1a2b3c4d5e6f",
	}
	tku.t.Run(fmt.Sprintf("get modify by id:%v", user.Id), func(t *testing.T) {
		id, err := tku.dbc.TypeModify("", &user)
		if err != nil {
			tku.t.Errorf("error modify user: %v", err)
		}
		if id != user.Id {
			tku.t.Errorf("error modify user: error returned id '%v' != '%v'", id, user.Id)
		}

		if have := database.ValidateStruct(&userGet); have != nil {
			tku.t.Errorf("error validate: %v", have)
		}
		_, err = tku.dbc.TypeGet("", &userGet, nil)
		if err != nil {
			tku.t.Errorf("error get user: %v", err)
		}

		if user != userGet {
			tku.t.Errorf(
				"error modify: fields is not modified %v '%v'->'%v' '%v'->'%v'",
				userGet.Id, user.Login, userGet.Login, user.Password, userGet.Password,
			)
		}
	})
}

func (tku *TestsUser) TestDeleteUser() {
	users := []db.User{
		{
			Id: "1a2b3c4d5e6f",
		},
		{
			Id: "1a2b3c4d5e",
		},
	}
	tku.t.Run("delete user", func(t *testing.T) {
		id, err := tku.dbc.TypeDelete("", &users[0], nil)
		if err != nil {
			tku.t.Errorf("error gelete user: %v", err)
		} else if id != users[0].Id {
			tku.t.Errorf("error delete user: error returned id '%v' != '%v'", id, users[0].Id)
		}

		id, err = tku.dbc.TypeDelete("", &users[1], nil)
		if err != nil {
			tku.t.Errorf("error gelete user: %v", err)
		} else if id != users[1].Id {
			tku.t.Errorf("error delete user: error returned id '%v' != '%v'", id, users[1].Id)
		}
	})
}
