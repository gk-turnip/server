package persistence

import (
	"sync"
)

import (
	"gk/game/config"
	"gk/gkerr"
	"gk/database"
)

type PersistenceContextDef struct {
	gameConfig *config.GameConfigDef
	connectionMutex *sync.Mutex
	connection *database.GkDbConDef
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

func (persistenceContext *PersistenceContextDef) GetLastPodName(userName string) (string, *gkerr.GkErrDef) {
	return "aaa", nil
}

func (persistenceContext *PersistenceContextDef) AddNewChatMessage (userName string, chatMessage string) *gkerr.GkErrDef {
	return persistenceContext.connection.AddNewChatMessage(userName, chatMessage)
}

func (persistenceContext *PersistenceContextDef) GetLastChatArchiveEntries(count int) ([]database.LugChatArchiveDef, *gkerr.GkErrDef) {
	return persistenceContext.connection.GetLastChatArchiveEntries(count)
}

