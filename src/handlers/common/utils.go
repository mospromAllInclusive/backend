package common

import (
	"backend/src/domains/entities"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"context"
	"log"
	"strconv"
)

func Ints64ToStrings(arr []int64) []string {
	res := make([]string, 0, len(arr))
	for _, v := range arr {
		res = append(res, strconv.FormatInt(v, 10))
	}
	return res
}

func SendActionToDBUsers(
	ctx context.Context,
	databasesService services.IDatabasesService,
	usersHub *web_sockets.Hub,
	dbID int64,
	action string,
) error {
	userIDs, err := databasesService.GetDatabasesUsersIDs(ctx, dbID)
	if err != nil {
		log.Printf("Error SendActionToDBUsers: %v", err)
		return err
	}
	usersHub.BroadcastMany(Ints64ToStrings(userIDs), action, nil)
	return nil
}

func ThrowUserFromDBTables(
	ctx context.Context,
	tablesService services.ITablesService,
	usersHub *web_sockets.Hub,
	userID int64,
	dbID int64,
) error {
	tableIDs, err := tablesService.ListIDsByDatabaseID(ctx, dbID)
	if err != nil {
		log.Printf("Error ThrowUserFromDBTables: %v", err)
		return err
	}

	for _, id := range tableIDs {
		usersHub.Broadcast(strconv.FormatInt(userID, 10), entities.EventActionGoAwayFromTable, &entities.GoAwayFromTableMessage{TableID: id})
	}
	return nil
}
