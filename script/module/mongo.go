package script

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type MgoModule struct {
	Url    string
	DBName string
	client *mongo.Client
}

func NewMgoModule() *MgoModule {
	return &MgoModule{}
}

func _decodeBsonElement(elemt bson.RawElement, notkey bool) string {

	k := elemt.Key()
	v := elemt.Value()

	var buf bytes.Buffer

	splicfunc := func(buf *bytes.Buffer, k, v interface{}) {
		if notkey {
			fmt.Fprintf(buf, `%v`, v)
		} else {
			fmt.Fprintf(buf, `"%s": %v`, k, v)
		}
	}

	switch v.Type {
	case bsontype.ObjectID:
		splicfunc(&buf, k, v.String())
	case bsontype.Double:
		val, _, _ := bsoncore.ReadValue(v.Value, v.Type)
		dval, _, _ := bsoncore.ReadDouble(val.Data)
		splicfunc(&buf, k, dval)
	case bsontype.String:
		splicfunc(&buf, k, v.String())
	case bsontype.Boolean:
		val, _, _ := bsoncore.ReadValue(v.Value, v.Type)
		dval, _, _ := bsoncore.ReadBoolean(val.Data)
		splicfunc(&buf, k, dval)
	case bsontype.EmbeddedDocument:
		val, _, _ := bsoncore.ReadValue(v.Value, v.Type)
		doc, _, _ := bsoncore.ReadDocument(val.Data)
		_elemts, err := bson.Raw(doc).Elements()
		if err != nil {
			fmt.Println("document parse err", err.Error())
		} else {
			splicfunc(&buf, k, _decodeBson(_elemts))
		}
	case bson.TypeArray:
		val, _, _ := bsoncore.ReadValue(v.Value, v.Type)
		arr, _, _ := bsoncore.ReadArray(val.Data)

		if len(arr) < 5 {
			fmt.Println("arr len < 5")
			return ""
		}

		var arrbuf bytes.Buffer
		arrbuf.WriteByte('[')

		length, rem, _ := bsoncore.ReadLength(arr)
		length -= 4

		var elem bsoncore.Element
		var ok bool
		for length > 1 {
			elem, rem, ok = bsoncore.ReadElement(rem)
			length -= int32(len(elem))
			if !ok {
				fmt.Println("ReadElement err")
				return ""
			}

			arrelm := _decodeBsonElement(bson.RawElement(elem), true)
			fmt.Fprintf(&arrbuf, "%s", arrelm)
			if length > 1 {
				arrbuf.WriteByte(',')
			}
		}
		if length != 1 { // Missing final null byte or inaccurate length
			fmt.Println("Missing final null byte or inaccurate length")
			return ""
		}

		arrbuf.WriteByte(']')
		splicfunc(&buf, k, arrbuf.String())
	default:
		fmt.Println("unknow tye", v.Type)
	}

	return buf.String()
}

func _decodeBson(elms []bson.RawElement) string {
	var buf bytes.Buffer
	buf.WriteByte('{')

	for i := 0; i < len(elms); i++ {

		val := _decodeBsonElement(elms[i], false)
		if val == "" {
			continue
		}

		buf.WriteString(val)

		if i < len(elms)-1 {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte('}')

	return buf.String()
}

func (m *MgoModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"conn":        m.conn,
		"disconn":     m.disconn,
		"insert_one":  m.insert_one,
		"insert_many": m.insert_many,
		"find":        m.find,
		"find_one":    m.find_one,
		"update_one":  m.update_one,
		"update_many": m.update_many,
		"delete_one":  m.delete_one,
		"delete_many": m.delete_many,
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
	clientOpt.SetMaxPoolSize(1)
	clientOpt.SetMaxConnIdleTime(60 * time.Second)

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

func (m *MgoModule) disconn(L *lua.LState) int {
	err := m.client.Disconnect(context.TODO())
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func _luaTable2Filter(doc *lua.LTable) (bson.M, error) {
	filter := bson.M{}

	m, err := utils.Table2MgoMap(doc)
	if err != nil {
		return filter, err
	}

	for k := range m {
		filter[k] = m[k]
	}

	return filter, nil
}

func (m *MgoModule) insert_many(L *lua.LState) int {
	err := m._insert_many(L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (m *MgoModule) _insert_many(collection string, doc *lua.LTable) error {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {
		v, err := utils.Table2MgoArr(doc)
		if err != nil {
			return err
		}

		_, err = coll.InsertMany(context.TODO(), v)
		return err
	}

	return nil
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
		m, err := utils.Table2MgoMap(doc)
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
		filter, err := _luaTable2Filter(doc)
		if err != nil {
			return "", err
		}

		cur, err := coll.Find(context.TODO(), filter)
		if err != nil {
			return "", err
		}

		for {

			if !cur.Next(context.Background()) {
				break
			}

			elms, _ := cur.Current.Elements()

			raws += _decodeBson(elms) + ","
		}

		raws = raws[0:len(raws)-1] + "]"
		fmt.Println(raws)

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

		filter, err := _luaTable2Filter(doc)
		if err != nil {
			return "", err
		}

		res := coll.FindOne(context.TODO(), filter)
		raw, err := res.DecodeBytes()
		if err != nil {
			return "", err
		}

		rawelement, err := raw.Elements()
		if err == nil {
			return _decodeBson(rawelement), nil
		} else {
			return "", fmt.Errorf("find one failed %v", err.Error())
		}
	}

	return "", fmt.Errorf("find one failed to get collection %v", collection)
}

func (m *MgoModule) update_many(L *lua.LState) int {

	err := m._update_many(L.ToString(1), L.ToTable(2), L.ToTable(3))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1

}

func (m *MgoModule) _update_many(collection string, f *lua.LTable, up *lua.LTable) error {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {

		filter, err := _luaTable2Filter(f)
		if err != nil {
			return err
		}

		update, err := _luaTable2Filter(up)
		if err != nil {
			return err
		}

		_, err = coll.UpdateMany(context.TODO(), filter, update)
		return err
	}

	return nil
}

func (m *MgoModule) update_one(L *lua.LState) int {

	err := m._update_one(L.ToString(1), L.ToTable(2), L.ToTable(3))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1
}

func (m *MgoModule) _update_one(collection string, f *lua.LTable, up *lua.LTable) error {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {

		filter, err := _luaTable2Filter(f)
		if err != nil {
			return err
		}

		update, err := _luaTable2Filter(up)
		if err != nil {
			return err
		}

		_, err = coll.UpdateOne(context.TODO(), filter, update)
		return err
	}

	return nil
}

func (m *MgoModule) delete_many(L *lua.LState) int {

	err := m._delete_many(L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 1
	}

	L.Push(lua.LString("succ"))
	return 1

}

func (m *MgoModule) _delete_many(collection string, doc *lua.LTable) error {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {
		filter, err := _luaTable2Filter(doc)
		if err != nil {
			return err
		}

		_, err = coll.DeleteMany(context.TODO(), filter)
		return err
	}

	return nil
}

func (m *MgoModule) delete_one(L *lua.LState) int {

	err := m._delete_one(L.ToString(1), L.ToTable(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		return 2
	}

	L.Push(lua.LString("succ"))
	return 2

}

func (m *MgoModule) _delete_one(collection string, doc *lua.LTable) error {

	coll := m.client.Database(m.DBName).Collection(collection)
	if coll != nil {
		filter, err := _luaTable2Filter(doc)
		if err != nil {
			return err
		}

		_, err = coll.DeleteOne(context.TODO(), filter)
		return err
	}

	return nil
}
