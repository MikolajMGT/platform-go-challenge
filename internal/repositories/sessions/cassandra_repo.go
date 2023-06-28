package sessions_db

import (
	"assets/pkg/logging"
	"context"
	"encoding/base32"
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

type CassandraRepo struct {
	logger  logging.Logger
	session *gocql.Session
}

type Session struct {
	Id     string
	Values map[interface{}]interface{}
}

func NewCassandraRepo(logger logging.Logger, session *gocql.Session) (repo *CassandraRepo) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := session.Query(CreateTableQuery()).WithContext(ctx).Exec(); err != nil {
		panic(errors.Wrap(err, "failed to inspect/create sessions table"))
	}

	return &CassandraRepo{logger: logger, session: session}
}

func (s *CassandraRepo) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *CassandraRepo) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.IsNew = true

	c := &Session{}
	if err := SelectRecordsById(s.session, name).Scan(&c.Id, &c.Values); err == nil {
		session.ID = c.Id
		session.Values = c.Values
		session.IsNew = false
	}

	return session, nil
}

func (s *CassandraRepo) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) (err error) {
	if session.ID == "" {
		session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}

	temp := make(map[string]interface{})
	for k, v := range session.Values {
		key, ok := k.(string)
		if !ok {
			return errors.New("Non-string key found in map")
		}
		temp[key] = v
	}

	var data []byte
	if data, err = json.Marshal(temp); err != nil {
		return err
	}

	if err := InsertQuery(s.session, session.ID, string(data)).Exec(); err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))

	return nil
}
