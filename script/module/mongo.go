package script

import (
	"context"
	"fmt"
	"time"

	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MgoModule struct {
	Url    string
	DBName string
	client *mongo.Client
}

func NewMgoModule() *MgoModule {
	return &MgoModule{}
}

func (m *MgoModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"conn":       m.conn,
		"insert_one": m.insert_one,
		"find":       m.find,
		"find_one":   m.find_one,
	})
	registerHttpResponseType(mod, L)
	L.Push(mod)
	return 1
}

func (m *MgoModule) conn(L *lua.LState) int {
	err := m._conn(L.ToString(1), L.ToString(2))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (m *MgoModule) _conn(db string, url string) error {

	clientOpt := options.Client()
	clientOpt.ApplyURI(url)
	clientOpt.SetConnectTimeout(5 * time.Second)
	clientOpt.SetMaxPoolSize(128)

	client, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	m.client = client
	m.DBName = db
	m.Url = url

	return nil
}

func (m *MgoModule) insert(L *lua.LState) int {

	return 0
}

func (m *MgoModule) _insert() {

}

func (m *MgoModule) insert_one(L *lua.LState) int {
	err := m._insert_one(L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (m *MgoModule) _insert_one(collection string, doc *lua.LTable) error {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {
		m, err := utils.Table2Map(doc)
		if err != nil {
			return err
		}

		_, err = coll.InsertOne(context.TODO(), m)
		return err
	}

	return fmt.Errorf("insert one failed to get collection %v", collection)
}

func (m *MgoModule) find(L *lua.LState) int {

	v, err := m._find(L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(lua.LString(v))
	L.Push(lua.LString("succ"))
	return 2
}

func (m *MgoModule) _find(collection string, doc *lua.LTable) (string, error) {

	coll := m.client.Database(m.DBName).Collection(collection)
	raws := "["
	if coll != nil {
		m, err := utils.Table2Map(doc)
		if err != nil {
			return "", err
		}
		filter := bson.M{}
		for k := range m {
			filter[k] = m[k]
		}
		cur, err := coll.Find(context.TODO(), m)
		if err != nil {
			return "", err
		}

		for {

			if !cur.Next(context.Background()) {
				break
			}

			raws += cur.Current.String() + ","
		}

		raws = raws[0:len(raws)-1] + "]"

		return raws, nil
	}

	return "", fmt.Errorf("find one failed to get collection %v", collection)
}

func (m *MgoModule) find_one(L *lua.LState) int {

	v, err := m._find_one(L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(lua.LString(v))
	L.Push(lua.LString("succ"))
	return 2
}

func (m *MgoModule) _find_one(collection string, doc *lua.LTable) (string, error) {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {
		m, err := utils.Table2Map(doc)
		if err != nil {
			return "", err
		}
		filter := bson.M{}
		for k := range m {
			filter[k] = m[k]
		}
		res := coll.FindOne(context.TODO(), filter)
		raw, err := res.DecodeBytes()
		if err != nil {
			return "", err
		}
		return raw.String(), nil

	}

	return "", fmt.Errorf("find one failed to get collection %v", collection)
}

func (m *MgoModule) update(L *lua.LState) int {
	return 0
}

func (m *MgoModule) _update() int {
	return 0
}

func (m *MgoModule) update_one(L *lua.LState) int {
	return 0
}

func (m *MgoModule) _update_one() int {
	return 0
}

func (m *MgoModule) delete(L *lua.LState) int {

	return 0
}

func (m *MgoModule) _delete() int {
	return 0
}

func (m *MgoModule) delete_one(L *lua.LState) int {

	return 0
}

func (m *MgoModule) _delete_one(L *lua.LState) int {

	return 0
}
