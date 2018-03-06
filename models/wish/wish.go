package wish

import (
	"project256/models"
	"project256/util"
	"time"
	"errors"
	"log"
	"fmt"
)

type WishStruct struct {
	Id			int64
	WishId		string
	ParentWishId 	string
	WishContent string
	Status 		int
	CreateUser	string
	CreateTime	int64
	UpdateTime	int64
}

func InsertWish(wishData *map[string]interface{}) (int64, error) {
	user := util.GetUserInfo()
	if user["status"] == util.STATUS_INVALID {
		return 0, errors.New(util.GetErrorMessage(util.ERROR_USER_UNAUTHORIZED))
	}
	wishId, err :=util.GenUUID32()
	db := models.GetDbConn()
	ret, err := db.Exec("INSERT INTO wish (wish_id, parent_wish_id, wish_content, status, create_user, create_time) VALUES (?,?,?,?,?,?)",
		wishId,
		(*wishData)["parent_wish_id"],
		(*wishData)["wish_content"],
		util.STATUS_VALID,
		user["user_id"],
		time.Now().Unix(),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Insert Wish Error: %s", err))
		return 0, err
	}
	row, _ := ret.RowsAffected()
	return row, err
}

func GetListByUser(userId string) (*[]WishStruct, error){
	db := models.GetDbConn()
	ret, err := db.Query("SELECT * FROM wish WHERE create_user=?",
		userId,
	)
	defer ret.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Get List By User Error: %s", err))
		return nil, err
	}
	_, err = ret.Columns()
	if err != nil {
		log.Fatal(fmt.Sprintf("Get List By User Error: %s", err))
		return nil, err
	}

	// 初始化结构
	var wishData WishStruct
	var wishDataList []WishStruct
	for ret.Next() {
		err = ret.Scan(&wishData.Id, &wishData.WishId, &wishData.ParentWishId, &wishData.WishContent, &wishData.Status, &wishData.CreateUser, &wishData.CreateTime, &wishData.UpdateTime)
		if err != nil {
			log.Fatal(fmt.Sprintf("Scan Data Error: %s", err))
			return nil, err
		}
		wishDataList = append(wishDataList, wishData)
	}

	return &wishDataList, err
}