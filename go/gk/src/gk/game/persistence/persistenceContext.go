package persistence

import (
	"sync"
)

import (
	"gk/database"
	"gk/game/config"
	"gk/gkerr"
)

type PersistenceContextDef struct {
	gameConfig      *config.GameConfigDef
	connectionMutex *sync.Mutex
	connection      *database.GkDbConDef
}

func NewPersistenceContext(gameConfig *config.GameConfigDef) (*PersistenceContextDef, *gkerr.GkErrDef) {
	var persistenceContext *PersistenceContextDef = new(PersistenceContextDef)
	var gkErr *gkerr.GkErrDef

	persistenceContext.gameConfig = gameConfig
	persistenceContext.connectionMutex = new(sync.Mutex)

	persistenceContext.connection, gkErr = database.NewGkDbCon(gameConfig.DatabaseUserName, gameConfig.DatabasePassword, gameConfig.DatabaseHost, gameConfig.DatabasePort, gameConfig.DatabaseDatabase)
	if gkErr != nil {
		return nil, gkErr
	}

	return persistenceContext, nil
}

func (persistenceContext *PersistenceContextDef) GetLastPodId(userName string) (int32, *gkerr.GkErrDef) {
	return 1, nil
}

func (persistenceContext *PersistenceContextDef) AddNewChatMessage(userName string, chatMessage string) *gkerr.GkErrDef {
	return persistenceContext.connection.AddNewChatMessage(userName, chatMessage)
}

func (persistenceContext *PersistenceContextDef) GetLastChatArchiveEntries(count int) ([]database.LugChatArchiveDef, *gkerr.GkErrDef) {
	return persistenceContext.connection.GetLastChatArchiveEntries(count)
}

func (persistenceContext *PersistenceContextDef) GetPodsList() ([]database.DbPodDef, *gkerr.GkErrDef) {
	return persistenceContext.connection.GetPodsList()
}

func (persistenceContext *PersistenceContextDef) SetUserPref(userName string, prefName string, prefValue string) *gkerr.GkErrDef {
	return persistenceContext.connection.SetUserPref(userName, prefName, prefValue)
}

func (persistenceContext *PersistenceContextDef) GetUserPrefsList(userName string) ([]database.DbUserPrefDef, *gkerr.GkErrDef) {
	return persistenceContext.connection.GetUserPrefsList(userName)
}
