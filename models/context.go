package models

type Context struct {
	User *Users
}

type Pagination struct {
	Page     int64 `query:"page" json:"page"`
	PageSize int64 `query:"page_size" json:"page_size"`
	Start    int64 `query:"-" json:"-"`
}

func (c *Pagination) Format() (err error) {
	if c.Page == 0 {
		c.Page = 1
	}
	if c.PageSize == 0 {
		c.PageSize = 15
	}
	c.Start = 0
	if c.Page > 1 {
		c.Start = (c.Page * c.PageSize) - c.PageSize
	}
	return
}
