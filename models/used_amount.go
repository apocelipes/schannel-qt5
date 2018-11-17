package models

import (
	"time"

	"github.com/astaxie/beego/orm"

	"schannel-qt5/parser"
)

// 用户的每日的流量使用记录，以service进行区分
type UsedAmount struct {
	Id       int       `orm:"auto"`
	Service  string    `orm:"size(50)"`
	Total    int
	Upload   int
	Download int
	Date     time.Time `orm:"type(date)"`
	User     *User     `orm:"rel(fk);on_delete(cascade)"`
}

// 设置数据库时区为UTC
func init() {
	orm.RegisterModel(&UsedAmount{})
	orm.DefaultTimeLoc = time.UTC
}

// SetUsedAmount 根据sevice插入使用量信息.
// 若date已经存在，则更新数据
func SetUsedAmount(db orm.Ormer,
	user string,
	service *parser.Service,
	total, upload, download int,
	date time.Time) error {
	amount := &UsedAmount{
		Service:  service.Name,
		Total:    total,
		Upload:   upload,
		Download: download,
		Date:     date,
		User: &User{
			Name: user,
		},
	}
	// 用户不存在无法insert
	err := db.QueryTable(amount.User).Filter("Name", user).One(amount.User)
	if err != nil {
		return err
	}

	cond := orm.NewCondition()
	cond = cond.And("Date", amount.Date).
		And("Service", amount.Service).
		And("User__Name", amount.User.Name)
	if db.QueryTable(amount).SetCond(cond).Exist() {
		db.QueryTable(amount).SetCond(cond).Update(orm.Params{
			"Upload":   amount.Upload,
			"Download": amount.Download,
		})
		return nil
	}

	if _, err := db.Insert(amount); err != nil {
		return err
	}

	return nil
}

const (
	// 需要获取的天数
	maxDays = 5
)

// GetRecentUsedAmount 返回date开始最近maxDays天的使用量数据
// 如果数量不足5天，则以最早一天的数据进行复制补足
func GetRecentUsedAmount(db orm.Ormer,
	user, service string,
	date time.Time) ([]*UsedAmount, error) {
	amounts := make([]*UsedAmount, 0)
	cond := orm.NewCondition()
	cond = cond.And("User__Name", user).
		And("Service", service).
		And("Date__lte", date)
	n, err := db.QueryTable(&UsedAmount{}).
		SetCond(cond).
		OrderBy("-Date").
		Limit(maxDays).All(&amounts)
	if err != nil {
		return nil, err
	}

	// 将时间转回本地时区
	for _, v := range amounts {
		v.Date = v.Date.In(time.Local)
	}

	if n < maxDays {
		amounts = paddingDate(amounts, int(n), maxDays)
	}

	return amounts, nil
}

// paddingDate 填充UsedAmount至长度为max
func paddingDate(amounts []*UsedAmount, length, max int) []*UsedAmount {
	origin := amounts[length-1]
	// 需要减去的天数
	dayCount := 1
	for i := length; i < max; i++ {
		duplicate := &UsedAmount{
			Service:  origin.Service,
			Total:    origin.Total,
			Upload:   origin.Upload,
			Download: origin.Download,
			Date:     origin.Date.Add(time.Duration(-dayCount) * time.Hour * 24),
			User:     origin.User,
		}
		dayCount++
		amounts = append(amounts, duplicate)
	}

	return amounts
}
