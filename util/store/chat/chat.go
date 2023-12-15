package UtilStore_Chat

import (
	"fmt"
	"gobasic/database"
	StoreChatStruct "gobasic/struct/store/chat"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetChatInfoByToken(ChatToken string) (StoreChatStruct.ChatData, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreChatStruct.ChatData{}, err
	}
	defer db.Close()

	var cart StoreChatStruct.ChatData
	query := `
		SELECT c.ChatId, c.ChatToken, 
		CASE WHEN c.Status = 'open' THEN true ELSE false END AS Open, 
		c.CreateDate, c.LastTalkDate FROM Chat c WHERE c.ChatToken = ?
    `
	err = db.QueryRow(query, ChatToken).Scan(
		&cart.ChatId,
		&cart.ChatToken,
		&cart.Open,
		&cart.CreateDate,
		&cart.LastTalkDate,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreChatStruct.ChatData{}, err
	}

	return cart, nil
}

func GetChatInfoById(ChatId int) (StoreChatStruct.ChatData, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreChatStruct.ChatData{}, err
	}
	defer db.Close()

	var cart StoreChatStruct.ChatData
	query := `
		SELECT c.ChatId, c.ChatToken, 
		CASE WHEN c.Status = 'open' THEN true ELSE false END AS Open, 
		c.CreateDate, c.LastTalkDate FROM Chat c WHERE c.ChatId = ?
    `
	err = db.QueryRow(query, ChatId).Scan(
		&cart.ChatId,
		&cart.ChatToken,
		&cart.Open,
		&cart.CreateDate,
		&cart.LastTalkDate,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreChatStruct.ChatData{}, err
	}

	return cart, nil
}

func GetWhoIMI(ChatId int, accountId any, accountType string) (StoreChatStruct.ChatProfile, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreChatStruct.ChatProfile{}, err
	}
	defer db.Close()

	var cart StoreChatStruct.ChatProfile
	query := `
		SELECT AccountId AS MemberId, MemberUUID AS AccountId FROM ChatMember WHERE ChatId = ? AND AccountId = ? AND AccountType = ?
    `
	err = db.QueryRow(query, ChatId, accountId, accountType).Scan(
		&cart.MemberId,
		&cart.MemberUUID,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreChatStruct.ChatProfile{}, err
	}

	return cart, nil
}

func GetMemberAccount(AccountId int, ChatId int, Type string) (StoreChatStruct.SourceData, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreChatStruct.SourceData{}, err
	}
	defer db.Close()

	var Member StoreChatStruct.SourceData
	query := `
		SELECT cm.ChatMemberId, cm.MemberUUID,
		CONCAT_WS('-', ac.Username, CONCAT(sl.Name, " ", sl.LastName)) AS AccountName,
		CONCAT_WS('-', ac.AccountUUID, sl.SellerUUID) AS AccountUUID,
		CONCAT_WS('-', ac.AccountId, sl.SellerId) AS AccountId, cm.AccountType
		FROM ChatMember cm
		LEFT JOIN Account ac ON ac.AccountId = cm.AccountId AND cm.AccountType = 'customer'
		LEFT JOIN Seller sl ON sl.SellerId = cm.AccountId AND cm.AccountType = 'seller' 
		WHERE (ac.AccountId = ? OR sl.SellerId = ?) AND AccountType = ? AND ChatId = ?;
    `
	err = db.QueryRow(query, AccountId, AccountId, Type, ChatId).Scan(
		&Member.ChatMemberId,
		&Member.MemberUUID,
		&Member.AccountName,
		&Member.AccountUUID,
		&Member.AccountId,
		&Member.AccountType,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreChatStruct.SourceData{}, err
	}

	return Member, nil
}

func GetMemberInChat(ChatId int) ([]StoreChatStruct.MemberDataV2, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return []StoreChatStruct.MemberDataV2{}, err
	}
	defer db.Close()

	query := `
		SELECT ChatMemberId, MemberUUID, AccountType FROM ChatMember WHERE ChatId = ?
    `
	rows, err := db.Query(query, ChatId)
	if err != nil {
		fmt.Print(err.Error())
		return []StoreChatStruct.MemberDataV2{}, err
	}
	defer rows.Close()

	var Members []StoreChatStruct.MemberDataV2

	for rows.Next() {
		var data StoreChatStruct.MemberDataV2

		err := rows.Scan(
			&data.ChatMemberId,
			&data.MemberUUID,
			&data.AccountType,
		)

		if err != nil {
			fmt.Print(err.Error())
			return []StoreChatStruct.MemberDataV2{}, err
		}
		Members = append(Members, data)
	}

	if len(Members) == 0 {
		fmt.Print(err.Error())
		return []StoreChatStruct.MemberDataV2{}, err
	}

	return Members, nil
}

func BoardcastSocket(Channel string, Message string) {
	url := "https://socket.patta.dev/send"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
	  "channel": "%s",
	  "message": %s,
	  "key": "124d6b6b-1a6c-449c-8932-0c8e5ba1d659"
	}`, Channel, Message))

	// fmt.Print(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(string(body))
}
