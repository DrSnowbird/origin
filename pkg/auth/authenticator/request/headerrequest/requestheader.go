package headerrequest

import (
	"net/http"
	"strings"

	"github.com/golang/glog"

	"github.com/openshift/origin/pkg/auth/api"
	authapi "github.com/openshift/origin/pkg/auth/api"
)

type Config struct {
	// UserNameHeaders lists the headers to check (in order, case-insensitively) for a username. The first header with a value wins.
	UserNameHeaders []string
}

func NewDefaultConfig() *Config {
	return &Config{
		UserNameHeaders: []string{"X-Remote-User"},
	}
}

type Authenticator struct {
	config *Config
	mapper authapi.UserIdentityMapper
}

func NewAuthenticator(config *Config, mapper authapi.UserIdentityMapper) *Authenticator {
	return &Authenticator{config, mapper}
}

func (a *Authenticator) AuthenticateRequest(req *http.Request) (api.UserInfo, bool, error) {
	username := ""
	for _, header := range a.config.UserNameHeaders {
		header = strings.TrimSpace(header)
		if len(header) == 0 {
			continue
		}
		username = req.Header.Get(header)
		if len(username) != 0 {
			break
		}
	}
	if len(username) == 0 {
		return nil, false, nil
	}

	identity := &authapi.DefaultUserIdentityInfo{
		UserName: username,
	}
	user, err := a.mapper.UserFor(identity)
	if err != nil {
		return nil, false, err
	}
	glog.V(4).Infof("Got userIdentityMapping: %#v", user)

	return user, true, nil
}
