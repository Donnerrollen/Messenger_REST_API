package repository

import (
	"context"
	"fmt"
	_ "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	todo "rest_API"
)

type ChatPostgres struct {
	db *pgx.Conn
}

func NewChatPostgres(db *pgx.Conn) *ChatPostgres {
	return &ChatPostgres{db: db}
}

func (r *ChatPostgres) LoadChatList(id int) ([]todo.ChatList, error) {

	var list []todo.ChatList
	var row todo.ChatList

	query := `
				SELECT DISTINCT
					ch."id",
					ch."name",
					count(msg2."id") as countNewMsg,
					max(coalesce(msg."id", 0)) lastMsgChat
				FROM ` + chatsTable + ` ch INNER JOIN
					 ` + usersChatTable + ` u_c ON ch."id" = u_c."idChat" LEFT JOIN
					 ` + messagesTable + ` msg ON ch."id" = msg."idChat" LEFT JOIN
					 ` + messagesTable + ` msg2 ON ch."id" = msg2."idChat" AND
												   msg2."id" > u_c.last_msg
				WHERE u_c."idUser" = $1
				GROUP BY ch."id", ch."name", msg2."id" 
				ORDER BY max(coalesce(msg."id", 0)) DESC, ch."name", ch."id"
			 `

	rows, _ := r.db.Query(context.Background(), query, id)
	for rows.Next() {
		rows.Scan(&row.Id, &row.Name, &row.CountNewMsg, &row.LastMsgChat)
		list = append(list, row)
	}
	rows.Close()

	return list, nil

}

func (r *ChatPostgres) CreateChat(chat todo.Chat) (todo.Chat, error) {

	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return todo.Chat{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	query := `
				INSERT INTO ` + chatsTable + ` (name, owner)
					VALUES ($1, $2)
					RETURNING id
			 `
	row := tx.QueryRow(context.Background(), query, chat.Name, chat.UserId)
	row.Scan(&chat.Id)

	query = `
				INSERT INTO ` + usersChatTable + ` ("idChat", "idUser", last_msg)
					VALUES ($1, $2, 0)
			 `
	_, err = tx.Exec(context.Background(), query, chat.Id, chat.UserId)
	if err != nil {
		return todo.Chat{}, err
	}

	return chat, nil
}

func (r *ChatPostgres) SearchChatByName(name string) ([]todo.SearchChat, error) {
	var chat []todo.SearchChat
	var row todo.SearchChat
	query := fmt.Sprintf("SELECT id, owner FROM %s WHERE name = '%s'", chatsTable, name)
	rows, _ := r.db.Query(context.Background(), query)
	for rows.Next() {
		rows.Scan(&row.Id, &row.UserId)
		chat = append(chat, row)
	}
	rows.Close()
	return chat, nil
}

func (r *ChatPostgres) CheckChatAndUser(userCheck todo.UserVerification) (access bool) {
	var RepUserCheck todo.UserVerification
	query := `
				SELECT "idUser", "idChat"
				FROM users_chat
				WHERE "idUser" = $1 AND 
					  "idChat" = $2
			 `
	row := r.db.QueryRow(context.Background(), query, userCheck.UserId, userCheck.ChatId)
	row.Scan(&RepUserCheck.UserId, &RepUserCheck.ChatId)

	if (userCheck.UserId == RepUserCheck.UserId) && (userCheck.ChatId == RepUserCheck.ChatId) {
		return true
	} else {
		return false
	}
}

func (r *ChatPostgres) SendMessage(message todo.SendMessage) (int, error) {

	query := `
				INSERT INTO ` + messagesTable + ` ("idUser", "idChat", "message")
					VALUES ($1, $2, $3)
					RETURNING "id"
			 `
	row := r.db.QueryRow(context.Background(), query, message.UserId, message.ChatId, message.Message)
	row.Scan(&message.Id)

	return message.Id, nil

}

func (r *ChatPostgres) AdminAccessVerification(userCheck todo.UserVerification) (access bool) {
	var id int

	query := `
				SELECT "id"
				FROM chats
				WHERE 
					"id" = $1 AND
					"owner" = $2
				LIMIT 1
			 `
	//		 fmt.Sprintf("SELECT id FROM %s WHERE (id = '%d' AND owner = '%d')", chatsTable, userCheck.ChatId, userCheck.UserId)
	row := r.db.QueryRow(context.Background(), query, userCheck.ChatId, userCheck.UserId)
	row.Scan(&id)
	if userCheck.ChatId == id {
		return true
	} else {
		return false
	}
}

func (r *ChatPostgres) LoadNewMsgById(userId, chatId int) ([]todo.ReadMessage, error) {

	var messages []todo.ReadMessage
	var row todo.ReadMessage

	query := `
				SELECT tmsg.id, us.username || '#' || tmsg."idUser", tmsg.message, to_char(tmsg.date_create, 'DD.mm.YYYY HH24:MI:SS')
				FROM public.messages tmsg INNER JOIN
					 public.users_chat u_c ON tmsg."idChat" = u_c."idChat" and
											  	 tmsg."id" > u_c.last_msg
											 INNER JOIN
                     public.users us ON tmsg."idUser" = us."id"
				WHERE u_c."idUser" = $1 and
					  tmsg."idChat" = $2
				ORDER BY tmsg.id
		`
	rows, _ := r.db.Query(context.Background(), query, userId, chatId)

	for rows.Next() {
		rows.Scan(&row.Id, &row.User, &row.Message, &row.Data)
		messages = append(messages, row)
	}
	rows.Close()

	query2 := `
				UPDATE users_chat
				SET last_msg = $3
				WHERE "idUser" = $1 and
					  "idChat" = $2
			  `
	_, err := r.db.Exec(context.Background(), query2, userId, chatId, row.Id)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *ChatPostgres) LoadOldMsgById(chatId, msgId int) ([]todo.ReadMessage, error) {

	type mRM []todo.ReadMessage
	var messages, t mRM
	var row todo.ReadMessage

	query := `
				SELECT tmsg.id, us.username || '#' || tmsg."idUser", tmsg.message, to_char(tmsg.date_create, 'DD.mm.YYYY HH24:MI:SS')
				FROM public.messages as tmsg INNER JOIN
                     public.users us ON tmsg."idUser" = us."id"
				WHERE tmsg."idChat" = $1 AND
					  (tmsg.id < $2 OR 0 = $2)
				ORDER BY tmsg.id desc
				LIMIT 10
			 `
	rows, _ := r.db.Query(context.Background(), query, chatId, msgId)
	for rows.Next() {
		rows.Scan(&row.Id, &row.User, &row.Message, &row.Data)
		t = make(mRM, 0)
		t = append(t, row)
		messages = append(t, messages...)
	}
	rows.Close()

	return messages, nil
}

func (r *ChatPostgres) CheckUserInSystem(userId int) (bool, error) {
	var userId_DB int

	query := `SELECT id FROM users WHERE users.id = $1`
	row := r.db.QueryRow(context.Background(), query, userId)
	row.Scan(&userId_DB)

	if userId == userId_DB {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *ChatPostgres) IsUserInChat(chatId, userId int) (bool, error) {

	var rec todo.UserVerification
	query := `
				SELECT "idUser", "idChat"
				FROM users_chat
				WHERE "idUser" = $1 AND 
					  "idChat" = $2
			 `
	row := r.db.QueryRow(context.Background(), query, userId, chatId)
	row.Scan(&rec.UserId, &rec.ChatId)

	if (rec.UserId == userId) && (rec.ChatId == chatId) {
		return true, nil
	} else {
		return false, nil
	}

}

func (r *ChatPostgres) CreateInvite(chatId, userId int) (int, error) {

	var row todo.UserVerification
	query := `
				SELECT "idUser", "idChat"
				FROM invitations
				WHERE "idUser" = $1 and
					  "idChat" = $2
			 `
	queryrow := r.db.QueryRow(context.Background(), query, userId, chatId)
	queryrow.Scan(&row.UserId, &row.ChatId)
	if (row.UserId == userId) && (row.ChatId == chatId) {
		return 1, nil
	}

	query = `INSERT INTO invitations ("idUser", "idChat") VALUES($1, $2)`
	_, err := r.db.Exec(context.Background(), query, userId, chatId)
	if err != nil {
		return 2, err
	}

	return 0, nil

}

func (r *ChatPostgres) GetUsersOfChat(chatId int) ([]todo.UsersList, error) {

	var listUsers []todo.UsersList
	var row todo.UsersList

	query := `
				SELECT us."id", us.name, us."id" = ch."owner" as IsOwner
				FROM users_chat u_ch INNER JOIN
					 users us ON u_ch."idUser" = us."id" INNER JOIN
					 chats ch ON u_ch."idChat" = ch."id"
				WHERE u_ch."idChat" = $1
				ORDER BY IsOwner desc, us.name, us."id"
			 `
	rows, _ := r.db.Query(context.Background(), query, chatId)
	for rows.Next() {
		rows.Scan(&row.Id, &row.Name, &row.IsOwner)
		listUsers = append(listUsers, row)
	}
	rows.Close()

	return listUsers, nil
}

func (r *ChatPostgres) AcceptInvite(userId, chatId int) (int, error) {
	var User todo.UserVerification
	var last_msg int

	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()
	query := `
				SELECT * FROM invitations WHERE ("idUser" = $1 AND "idChat" = $2)
			`
	row := tx.QueryRow(context.Background(), query, userId, chatId)
	row.Scan(&User.UserId, &User.ChatId)
	if (userId != User.UserId) && (chatId != User.ChatId) {
		return 1, nil // Значит у пользователя нет этого приглашения
	}

	query = `
				DELETE FROM invitations WHERE ("idUser" = $1 AND "idChat" = $2)
			 `
	_, err = tx.Exec(context.Background(), query, userId, chatId)
	if err != nil {
		return 0, err
	}

	query = `
				SELECT id FROM messages WHERE ("idChat" = $1)
				ORDER BY id DESC
				LIMIT 1
			`
	row = tx.QueryRow(context.Background(), query, chatId)
	row.Scan(&last_msg)

	query = `
				INSERT INTO users_chat VALUES ($1, $2, $3)

			`
	_, err = tx.Exec(context.Background(), query, chatId, userId, last_msg)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (r *ChatPostgres) DenyInvite(UserId, ChatId int) (int, error) {
	var User todo.UserVerification

	query := `
				SELECT * FROM invitations WHERE ("idUser" = $1 AND "idChat" = $2)
			`
	row := r.db.QueryRow(context.Background(), query, UserId, ChatId)
	row.Scan(&User.UserId, &User.ChatId)
	if (UserId != User.UserId) && (ChatId != User.ChatId) {
		return 1, nil // Значит у пользователя нет этого приглашения
	}

	query = `
				DELETE FROM invitations WHERE ("idUser" = $1 AND "idChat" = $2)
			`
	_, err := r.db.Exec(context.Background(), query, UserId, ChatId)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

func (r *ChatPostgres) LogOut(UserId, ChatId int) (int, error) {
	var User todo.UserVerification

	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()
	query := `
				SELECT id, owner FROM chats WHERE (id = $1 AND owner = $2)
			 `
	row := tx.QueryRow(context.Background(), query, ChatId, UserId)
	row.Scan(&User.ChatId, &User.UserId)
	if (UserId == User.UserId) && (ChatId == User.ChatId) {
		return 1, nil // Значит, что пользователь администратор и пытается удалить сам себя
	}

	query = `
				DELETE FROM users_chat
				WHERE ("idUser" = $1 AND
					   "idChat" = $2)
			`

	_, err = tx.Exec(context.Background(), query, UserId, ChatId)
	if err != nil {
		return 0, err
	}
	return 2, nil
}

func (r *ChatPostgres) DeleteUser(UserId, ChatId int) (int, error) {
	var User todo.UserVerification

	query := `
				SELECT "idChat", "idUser"
				FROM users_chat
				WHERE ("idChat" = $1 AND
					   "idUser" = $2)
			 `
	row := r.db.QueryRow(context.Background(), query, ChatId, UserId)
	row.Scan(&User.ChatId, &User.UserId)

	if (User.ChatId != ChatId) && (User.UserId != UserId) {
		return 2, nil
	}

	query = `
				DELETE FROM users_chat
				WHERE ("idChat" = $1 AND
					   "idUser" = $2)
			`
	_, err := r.db.Exec(context.Background(), query, ChatId, UserId)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (r *ChatPostgres) RenameChat(RenameChat todo.RenameChat) error {
	query := `
				UPDATE public.chats
				SET name = '$1'
				WHERE id = $2 AND
					  owner = $3;
			 `

	_, err := r.db.Exec(context.Background(), query, RenameChat.ChatName, RenameChat.ChatId, RenameChat.Owner)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatPostgres) DeleteChat(ChatId int) error {

	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	query := `
				DELETE FROM chats
				WHERE (id = $1)
			 `
	_, err = tx.Exec(context.Background(), query, ChatId)
	if err != nil {
		return err
	}

	query = `
				DELETE FROM invitations
				WHERE (id = $1)
			 `
	_, err = tx.Exec(context.Background(), query, ChatId)
	if err != nil {
		return err
	}

	query = `
				DELETE FROM messages
				WHERE (id = $1)
			 `
	_, err = tx.Exec(context.Background(), query, ChatId)
	if err != nil {
		return err
	}

	query = `
				DELETE FROM users_chat
				WHERE (id = $1)
			 `
	_, err = tx.Exec(context.Background(), query, ChatId)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatPostgres) GiveAdminStatus(ChatId, UserId int) error {
	query := `
				UPDATE public.chats
				SET owner = $1
				WHERE id = $2
			 `
	_, err := r.db.Exec(context.Background(), query, UserId, ChatId)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatPostgres) LoadNick(UserId int) (string, int, error) {
	var NickName string

	query := `
				SELECT name
				FROM users
				WHERE id = $1
			`
	row := r.db.QueryRow(context.Background(), query, UserId)
	row.Scan(&NickName)
	return NickName, UserId, nil
}
