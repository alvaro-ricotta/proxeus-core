package service

import (
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	UserService interface {
		GetUser(auth model.Auth) (*model.User, error)
		GetById(auth model.Auth, id string) (*model.User, error)
		GetUserDataById(auth model.Auth, id string) (*model.UserDataItem, error)
		CreateApiKeyHandler(auth model.Auth, userId, apiKeyName string) (string, error)
		DeleteUser(auth model.Auth) error
		DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error
		DeleteUserData(auth model.Auth, id string) error
		GetByBCAddress(blockchainAddress string) (*model.User, error)
	}
	defaultUserService struct {
	}
)

func NewUserService() *defaultUserService {
	return &defaultUserService{}
}

func (me *defaultUserService) GetUser(auth model.Auth) (*model.User, error) {
	return userDB().Get(auth, auth.UserID())
}
func (me *defaultUserService) GetById(auth model.Auth, id string) (*model.User, error) {
	return userDB().Get(auth, id)
}

func (me *defaultUserService) GetUserDataById(auth model.Auth, id string) (*model.UserDataItem, error) {
	return userDataDB().Get(auth, id)
}

func (me *defaultUserService) DeleteUser(auth model.Auth) error {
	//remove documents / workflow instances of user
	workflowInstances, err := userDataDB().List(auth, "", storage.Options{}, false)
	if err != nil && !db.NotFound(err) {
		return err
	}
	for _, workflowInstance := range workflowInstances {
		//err = userDataDB().Delete(auth, c.System().DB.Files, workflowInstance.ID)
		err = me.DeleteUserData(auth, workflowInstance.ID)
		if err != nil {
			return err
		}
	}

	//set workflow templates to deactivated
	workflows, err := workflowDB().List(auth, "", storage.Options{})
	if err != nil && !db.NotFound(err) {
		return err
	}
	for _, workflow := range workflows {
		if workflow.OwnedBy(auth) {
			workflow.Deactivated = true
			err = workflowDB().Put(auth, workflow)
			if err != nil {
				return err
			}
		}
	}

	// unset user data and set inactive
	user, err := userDB().Get(auth, auth.UserID())
	if err != nil {
		return err
	}
	user.Active = false
	user.EthereumAddr = "0x"
	user.Email = ""
	user.Name = ""
	user.Photo = ""
	user.PhotoPath = ""
	user.WantToBeFound = false

	return userDB().Put(auth, user)
}

func (me *defaultUserService) DeleteUserData(auth model.Auth, id string) error {
	return userDataDB().Delete(auth, filesDB(), id)
}

func (me *defaultUserService) CreateApiKeyHandler(auth model.Auth, userId, apiKeyName string) (string, error) {
	return userDB().CreateApiKey(auth, userId, apiKeyName)
}

func (me *defaultUserService) DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error {
	return userDB().DeleteApiKey(auth, userId, hiddenApiKey)
}

func (me *defaultUserService) GetByBCAddress(blockchainAddress string) (*model.User, error) {
	return userDB().GetByBCAddress(blockchainAddress)
}
