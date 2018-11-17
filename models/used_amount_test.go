package models

import (
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"

	"schannel-qt5/parser"
)

const (
	amountPath = "/tmp/used_amount.db"
)

func init() {
	orm.Debug = true
	orm.RegisterDataBase("testAmount", "sqlite3", amountPath)
}

// 用于amount插入测试的结构
type dummyAmount struct {
	service                 *parser.Service
	total, upload, download int
	date                    time.Time
	user                    string
}

// 初始化测试环境
func initAmount(t *testing.T) (orm.Ormer, []*User, []*dummyAmount) {
	os.Truncate(amountPath, 0)
	err := orm.RunSyncdb("testAmount", false, true)
	if err != nil {
		t.Fatal(err)
	}

	users := []*User{
		{
			Name:   "test@test.com",
			Passwd: genPassword(),
		},
		{
			Name:   "example",
			Passwd: genPassword(),
		},
	}

	db := orm.NewOrm()
	db.Using("testAmount")
	for _, v := range users {
		if err := SetUserPassword(db, v.Name, v.Passwd); err != nil {
			t.Fatalf("initdb user error: %v\n", err)
		}
	}

	amounts := []*dummyAmount{
		{
			service:  &parser.Service{Name: "A"},
			total:    10000,
			upload:   1500,
			download: 9000,
			date:     parser.GetCurrentDay(),
			user:     users[0].Name,
		},
		{
			service:  &parser.Service{Name: "A"},
			total:    10000,
			upload:   1000,
			download: 9000,
			date:     parser.GetCurrentDay().Add(time.Duration(-1) * time.Hour * 24),
			user:     users[0].Name,
		},
		{
			service:  &parser.Service{Name: "B"},
			total:    20000,
			upload:   5000,
			download: 9000,
			date:     parser.GetCurrentDay(),
			user:     users[0].Name,
		},
		{
			service:  &parser.Service{Name: "B"},
			total:    20000,
			upload:   1400,
			download: 13000,
			date:     parser.GetCurrentDay(),
			user:     users[1].Name,
		},
	}

	for _, v := range amounts {
		if err := SetUsedAmount(db, v.user, v.service, v.total, v.upload, v.download, v.date); err != nil {
			t.Fatalf("initdb used_amount error: %v\n", err)
		}
	}

	return db, users, amounts
}

func TestSetUsedAmount(t *testing.T) {
	db, users, _ := initAmount(t)
	testData := []*struct {
		data     *dummyAmount
		inserted bool
	}{
		{
			data: &dummyAmount{
				service:  &parser.Service{Name: "test"},
				total:    7000,
				upload:   0,
				download: 2000,
				date:     parser.GetCurrentDay().Add(time.Duration(-2) * time.Hour * 24),
				user:     users[0].Name,
			},
			inserted: true,
		},
		// 测试更新
		{
			data: &dummyAmount{
				service:  &parser.Service{Name: "A"},
				total:    10000,
				upload:   0,
				download: 0,
				date:     parser.GetCurrentDay(),
				user:     users[0].Name,
			},
			inserted: true,
		},
		{
			data: &dummyAmount{
				service:  &parser.Service{Name: "B"},
				total:    7000,
				upload:   2000,
				download: 0,
				date:     parser.GetCurrentDay().Add(time.Duration(-1) * time.Hour * 24),
				user:     users[1].Name,
			},
			inserted: true,
		},
		{
			data: &dummyAmount{
				service:  &parser.Service{Name: "test"},
				total:    7000,
				upload:   0,
				download: 2000,
				date:     time.Time{},
				user:     "notExists",
			},
			inserted: false,
		},
	}

	for _, v := range testData {
		d := v.data
		err := SetUsedAmount(db, d.user, d.service, d.total, d.upload, d.download, d.date)
		if (err == nil) != v.inserted {
			format := "insert error:\n\twant: %v\n\thave: %v\n\tinfo: %+v\n"
			t.Errorf(format, v.inserted, err, d)
		}
	}
}

func TestGetRecentUsedAmount(t *testing.T) {
	db, users, amounts := initAmount(t)
	testData := []*struct {
		service string
		user    string
		date    time.Time
		res     []*UsedAmount
	}{
		{
			service: "A",
			user:    users[0].Name,
			date:    parser.GetCurrentDay(),
			res: []*UsedAmount{
				{
					Service:  "A",
					Total:    amounts[0].total,
					Upload:   amounts[0].upload,
					Download: amounts[0].download,
					Date:     amounts[0].date,
					User:     users[0],
				},
				{
					Service:  "A",
					Total:    amounts[1].total,
					Upload:   amounts[1].upload,
					Download: amounts[1].download,
					Date:     amounts[1].date,
					User:     users[0],
				},
				{
					Service:  "A",
					Total:    amounts[1].total,
					Upload:   amounts[1].upload,
					Download: amounts[1].download,
					Date:     amounts[1].date.Add(time.Duration(-1) * time.Hour * 24),
					User:     users[0],
				},
				{
					Service:  "A",
					Total:    amounts[1].total,
					Upload:   amounts[1].upload,
					Download: amounts[1].download,
					Date:     amounts[1].date.Add(time.Duration(-2) * time.Hour * 24),
					User:     users[0],
				},
				{
					Service:  "A",
					Total:    amounts[1].total,
					Upload:   amounts[1].upload,
					Download: amounts[1].download,
					Date:     amounts[1].date.Add(time.Duration(-3) * time.Hour * 24),
					User:     users[0],
				},
			},
		},
	}

	for _, v := range testData {
		res, err := GetRecentUsedAmount(db, v.user, v.service, v.date)
		if err != nil {
			t.Fatalf("GetRecentUsedAmount error: %v\n", err)
		}
		if !sliceEqual(res, v.res) {
			t.Error("content wrong in GetRecentUsedAmount")
		}
	}
}

func TestPaddingDate(t *testing.T) {
	testData := []*struct {
		origin []*UsedAmount
		max    int
		res    []*UsedAmount
	}{
		{
			origin: []*UsedAmount{
				{
					Service:  "A",
					Total:    70000,
					Upload:   4500,
					Download: 34600,
					Date:     time.Date(2018, 11, 2, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
				{
					Service:  "A",
					Total:    70000,
					Upload:   4660,
					Download: 37100,
					Date:     time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
			},
			max: 2,
			res: []*UsedAmount{
				{
					Service:  "A",
					Total:    70000,
					Upload:   4500,
					Download: 34600,
					Date:     time.Date(2018, 11, 2, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
				{
					Service:  "A",
					Total:    70000,
					Upload:   4660,
					Download: 37100,
					Date:     time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
			},
		},
		{
			origin: []*UsedAmount{
				{
					Service:  "A",
					Total:    70000,
					Upload:   4500,
					Download: 34600,
					Date:     time.Date(2018, 11, 2, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
				{
					Service:  "A",
					Total:    70000,
					Upload:   4660,
					Download: 37100,
					Date:     time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
			},
			max: 4,
			res: []*UsedAmount{
				{
					Service:  "A",
					Total:    70000,
					Upload:   4500,
					Download: 34600,
					Date:     time.Date(2018, 11, 2, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
				{
					Service:  "A",
					Total:    70000,
					Upload:   4660,
					Download: 37100,
					Date:     time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
				{
					Service:  "A",
					Total:    70000,
					Upload:   4660,
					Download: 37100,
					Date:     time.Date(2018, 10, 31, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
				{
					Service:  "A",
					Total:    70000,
					Upload:   4660,
					Download: 37100,
					Date:     time.Date(2018, 10, 30, 0, 0, 0, 0, time.UTC),
					User:     nil,
				},
			},
		},
	}

	for _, v := range testData {
		res := paddingDate(v.origin, len(v.origin), v.max)
		if len(res) != len(v.res) {
			format := "length wrong\n\twant: %v\n\thave :%v\n"
			t.Errorf(format, v.max, len(res))
		}
		if !sliceEqual(res, v.res) {
			t.Errorf("content wrong in padding %d length slice\n", v.max)
		}
	}
}

// 测试两个slice中的每个指针指向的内容是否相同
// a, b长度必须相等
func sliceEqual(a, b []*UsedAmount) bool {
	for i := 0; i < len(a); i++ {
		if (a[i].User == nil) && (b[i].User == nil) {
			return true
		}

		if a[i].Download != b[i].Download ||
			a[i].Upload != b[i].Upload ||
			a[i].Total != b[i].Total ||
			a[i].Service != b[i].Service ||
			a[i].Date != b[i].Date ||
			a[i].User.Name != b[i].User.Name {
			return false
		}
	}

	return true
}
