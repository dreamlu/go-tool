// package gt

package gt

import (
	"fmt"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"github.com/dreamlu/gt/tool/type/time"
	"log"
	"net/url"
	"testing"
	time2 "time"
)

// params param.CMap is web request GET params
// in golang, it was url.values

// user model
type User struct {
	ID         uint64     `json:"id"`
	Name       string     `json:"name"`
	Createtime time.CTime `json:"createtime"`
}

// service model
type Service struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// order model
type Order struct {
	ID         uint64     `json:"id"`
	UserID     int64      `json:"user_id"`    // user id
	ServiceID  int64      `json:"service_id"` // service table id
	Createtime time.CTime `json:"createtime"` // createtime
}

// order detail
type OrderD struct {
	Order
	UserName    string `json:"user_name"`    // user table column name
	ServiceName string `json:"service_name"` // service table column `name`
}

// 局部
var crud = NewCrud()

func TestDB(t *testing.T) {

	var user = User{
		Name: "测试xx",
		//Createtime:JsonDate(time.Now()),
	}

	// return create id
	crud.DB().CreateData(&user)
	t.Log("user: ", user)
	t.Log(crud.DB().RowsAffected)
	user.Name = "haha"
	crud.DB().CreateData(&user)
	t.Log("user: ", user)
	t.Log(crud.DB().RowsAffected)
}

// 通用增删该查测试
// 传参可使用url.Values替代param.CMap操作方便
func TestCrud(t *testing.T) {

	// add
	var user = User{
		Name: "new name",
	}
	crud = NewCrud(
		Model(User{}),
		//Table("user"),
		Data(&user),
	).Create()

	// update
	user.ID = 2
	crud.Update()

	// get by id
	info := crud.GetByID(1)
	t.Log(info, "\n[GetByID]:", user)

	// get by search
	var args = url.Values{}
	args.Add("name", "梦")
	// get by search
	var users []*User
	crud.Params(
		Model(User{}),
		Data(&user),
		SubWhereSQL("1=1"),
	)
	//args["name"][0] = "梦"
	var params cmap.CMap
	params = params.CMap(args)
	crud.GetBySearch(params)
	t.Log("\n[User Info]:", users)

	// delete
	info2 := crud.Delete(12)
	t.Log(info2.Error())

	// update by form request
	args.Add("id", "4")
	args.Set("name", "梦4")
	err := crud.UpdateForm(cmap.CMap(args))
	log.Println(err)
}

// select sql
func TestCrudSQL(t *testing.T) {
	sql := "update `user` set name=? where id=?"
	t.Log("[Info]:", crud.Select(sql, "梦sql", 1).Select(" and 1=1").Exec())
	t.Log("[Info]:", crud.Select(sql, "梦sql", 1).Exec())
	t.Log("[Info]:", crud.DB().RowsAffected)
	var user []User
	sql = "select * from `user` where name=? and id=?"
	cd := NewCrud()
	t.Log("[Info]:", cd.Params(Data(&user)).Select(sql, "梦sql", 1).Select(" and 1=1").Exec())
	t.Log("[Info]:", cd.Params(Data(&user)).Select(sql, "梦sql", 1).Exec())
}

// 通用分页测试
// 如：
func TestSqlSearch(t *testing.T) {

	type UserInfo struct {
		ID       uint64     `json:"id"`
		UserID   int64      `json:"user_id"`   //用户id
		UserName string     `json:"user_name"` //用户名
		Userinfo json.CJSON `json:"userinfo"`
	}
	sql := fmt.Sprintf(`select a.id,a.user_id,a.userinfo,b.name as user_name from userinfo a inner join user b on a.user_id=b.id where 1=1 and `)
	sqlNt := `select
		count(distinct a.id) as total_num
		from userinfo a inner join user b on a.user_id=b.id
		where 1=1 and `
	var ui []UserInfo

	//页码,每页数量
	clientPage := int64(1) //默认第1页
	everyPage := int64(10) //默认10页

	sql = string([]byte(sql)[:len(sql)-4]) //去and
	sqlNt = string([]byte(sqlNt)[:len(sqlNt)-4])
	sql += "order by a.id "
	log.Println(crud.DB().GetDataBySQLSearch(&ui, sql, sqlNt, clientPage, everyPage, nil, nil))
	log.Println(ui[0].Userinfo.String())
}

// select 数据存在验证
func TestValidateData(t *testing.T) {
	sql := "select *from `user` where id=2"
	ss := crud.DB().ValidateSQL(sql)
	log.Println(ss)
}

// 分页搜索中key测试
func TestGetSearchSql(t *testing.T) {

	type UserDe struct {
		User
		Num int64 `json:"num" gt:"sub_sql"`
	}

	var args = make(cmap.CMap)
	args["clientPage"] = append(args["clientPage"], "1")
	args["everyPage"] = append(args["everyPage"], "2")
	//args["key"] = append(args["key"], "梦 嘿,伙计")
	//sub_sql := ",(select aa.name from shop aa where aa.user_id = a.id) as shop_name"
	sqlNt, sql, _, _, _ := GetSearchSQL(&GT{
		Params: nil,
		CMaps:  args,
		Select: "",
		From:   "",
		Group:  "",
		Args:   nil,
		ArgsNt: nil,
	})
	log.Println("SQLNOLIMIT:", sqlNt, "\nSQL:", sql)

	// 两张表，待重新测试
	// sqlNt, sql, _, _ = GetDoubleSearchSQL(UserInfo{}, "userinfo", "user", args)
	// log.Println("SQLNOLIMIT==>2:", sqlNt, "\nSQL==>2:", sql)

}

// 通用sql以及参数
func TestGetDataBySql(t *testing.T) {
	var sql = "select id,name,createtime from `user` where id = ?"

	var user User
	crud.DB().GetDataBySQL(&user, sql, "1")
	t.Log(user)
}

func TestGetDataBySearch(t *testing.T) {
	var args = make(cmap.CMap)
	args.Add("name", "梦")
	//args["name"] = append(args["name"], "梦")
	args["key"] = append(args["key"], "梦")
	args["clientPage"] = append(args["clientPage"], "1")
	args["everyPage"] = append(args["everyPage"], "2")
	var user []*User
	crud.DB().GetDataBySearch(&GT{
		CMaps: args,
		Params: &Params{
			Table: "user",
			Model: User{},
			Data:  &user,
		},
	})
	t.Log(user[0])
}

// 测试多表连接
func TestGetMoreDataBySearch(t *testing.T) {
	// 多表查询
	// get more search
	var params = make(cmap.CMap)
	params.Add("user_id", "1")
	//params.Add("key", "梦") // key work
	params.Add("clientPage", "1")
	params.Add("everyPage", "2")
	var or []*OrderD
	crud := NewCrud(
		InnerTable([]string{"order", "user", "order", "service"}),
		//LeftTable([]string{"order", "service"}),
		Model(OrderD{}),
		Data(&or),
		SubWhereSQL("1 = 1", "2 = 2", ""),
	)
	err := crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
	t.Log("\n[User Info]:", or[0])
}

func TestGetMoreSearchSQL(t *testing.T) {
	type ClientVipBehavior struct {
		ID          int64      `gorm:"type:bigint(20)" json:"id"`
		ClientVipID int64      `gorm:"type:bigint(20)" json:"client_vip_id"`
		ShopId      int64      `gorm:"type:bigint(20)" json:"shop_id"`
		StaffId     int64      `gorm:"type:bigint(20)" json:"staff_id"`
		Status      int64      `gorm:"type:tinyint(2);DEFAULT:0" json:"status"`
		Num         int64      `json:"num" gorm:"type:int(11)"` // 第几次参加
		Createtime  time.CTime `gorm:"type:datetime" json:"createtime"`
	}

	// TODO bug is_sp
	// 客户行为详情
	type ClientVipBehaviorDe struct {
		ClientVipBehavior
		ClientName    string `json:"client_name"`
		ClientHeadimg string `json:"client_headimg"`
		VipType       int64  `json:"vip_type" gt:"sub_sql"` // 0意向会员, 1会员
		IsSp          int64  `json:"-" gt:"field:is_sp"`    // 是否代言人, 0不是, 1是
	}
	gt := &GT{
		Params: &Params{
			InnerTable: []string{"client_vip_behavior", "client_vip", "client_vip", "client"},
			Model:      ClientVipBehaviorDe{},
		},
	}
	sqlNt, sql, _, _, _ := GetMoreSearchSQL(gt)
	t.Log(sqlNt)
	t.Log(sql)
}

// 批量创建
func TestCreateMoreData(t *testing.T) {

	type UserPar struct {
		Name       string     `json:"name"`
		Createtime time.CTime `json:"createtime"`
	}
	type User struct {
		ID uint64 `json:"id"`
		UserPar
	}

	var up = []UserPar{
		{Name: "测试1", Createtime: time.CTime(time2.Now())},
		{Name: "测试2"},
	}
	crud := NewCrud(
		//Table("user"),
		Model(UserPar{}),
		Data(up),
		//SubSQL("(asdf) as a","(asdfa) as b"),
	)

	err := crud.CreateMoreData()
	t.Log(err)
}

// 继承tag解析测试
func TestExtends(t *testing.T) {
	type UserDe struct {
		User
		Other string `json:"other"`
	}

	type UserDeX struct {
		a []string
		UserDe
		OtherX string `json:"other_x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name"`
		UserDeX
	}
	for i := 0; i < 3; i++ {
		t.Log(GetColSQL(UserDeX{}))
		t.Log(GetMoreTableColumnSQL(UserMore{}, []string{"user", "shop"}[:]...))
	}
}

// select test
func TestDBCrud_Select(t *testing.T) {
	var user []*User
	crud.Params(
		Data(&user),
		ClientPage(1),
		EveryPage(2),
	).
		Select("select *from user").
		Select("where id > 0")
	if true {
		crud.Select("and 1=1")
	}
	crud.Search()
	t.Log(crud.Pager())
	crud.Single()
}

// test update/delete
func TestDBCrud_Update(t *testing.T) {

	type UserPar struct {
		Name string `json:"name"`
	}
	crud := crud.Params(
		//Table("user"),
		Model(User{}),
		Data(&UserPar{
			//ID:   1,
			Name: "梦S",
		}),
	)
	t.Log(crud.Update().RowsAffected())
	t.Log(crud.Select("`name` = ?", "梦").Update().RowsAffected())
	t.Log(crud.Error())
}

// test update/delete
func TestDBCrud_Create(t *testing.T) {

	crud.Params(
		Table("user"),
		Data(&User{
			ID:   11234,
			Name: "梦S",
		}),
	)
	t.Log(crud.Error())
	t.Log(crud.Create().Error())
	crud.Params(
		Data(&User{
			Name: "梦SSS2",
		})).Create()
	t.Log(crud.Error())
}

// test Transcation
func TestTranscation(t *testing.T) {

	cd := crud.Begin()
	cd.Params(
		Table("user"),
		Data(&User{
			ID:   11234,
			Name: "梦S",
		}),
	).Create()
	if cd.Error() != nil {
		cd.Rollback()
	}
	cd.Params(
		Data(&User{
			Name: "梦SSS2",
		})).Create()
	if cd.Error() != nil {
		cd.Rollback()
	}
	// add select sql test
	var u []User
	cd.Params(Data(&u)).Select("select * from `user`").Select("where 1=1").Single()
	cd.Params(Data(&u)).Select("select * from `user`").Select("where 1=1").Single()
	//cd.DB().Raw("select * from `user`").Scan(&u)

	cd.Commit()
	if cd.Error() != nil {
		cd.Rollback()
	}
}

func TestGetReflectTagMore(t *testing.T) {
	//type GroupmealCategory struct {
	//	ID   int64  `gorm:"type:bigint(20) AUTO_INCREMENT;PRIMARY_KEY;" json:"id"` //编号
	//	Name string `gorm:"type:varchar(128);NOT NULL;" json:"name"`               //类型
	//}
	type Groupmeal struct {
		ID                  int64  `gorm:"type:bigint(20);AUTO_INCREMENT;PRIMARY_KEY;" json:"id"`
		GroupmealCategoryID string `gorm:"type:varchar(128);NOT NULL;" json:"groupmeal_category_id"`
	}
	type GroupmealModel struct {
		Groupmeal
		GroupmealCategoryName string `json:"groupmeal_category_name"`
	}
	var data []*GroupmealModel
	crud.Params(
		Data(&data),
		Model(GroupmealModel{}),
		InnerTable([]string{"groupmeal", "groupmeal_category"}))
	var params = make(cmap.CMap)
	crud.GetMoreBySearch(params)
}

func TestGetColSQLAlias(t *testing.T) {
	sql := GetColSQLAlias(User{}, "a")
	t.Log(sql)
}

func TestGetMoreSQL(t *testing.T) {
	// table: venuepricets
	// related table: venue/venuehomestay
	type Venuepricets struct {
		ID      int `json:"id"`
		VenueID int `json:"venue_id"` // different table
		// type, related table: venue/venuehomestay
		Type *int `json:"type" gorm:"type:tinyint(2);DEFAULT:0"`
	}
	// 后台 特价
	type VpsInfo struct {
		Venuepricets
		VenueName string `json:"venue_name" gt:"field:venuehomestay_name"`
	}
	var vsi []VpsInfo

	var param = cmap.CMap{}
	param.Add("key", "test")
	//param.Add()
	crud := NewCrud(
		Model(VpsInfo{}),
		Data(&vsi),
		InnerTable([]string{"venuepricets:venue_id", "venuehomestay"}),
	)
	for i := 0; i < 3; i++ {
		cd := crud.GetMoreBySearch(param)
		if cd.Error() != nil {
			t.Log(cd.Error())
		}
	}
}
