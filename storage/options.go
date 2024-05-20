package storage 

type options struct {
	id int
	uuid string
}

type Option func(c *options)
var Options options


func (o options) Id(id int) Option {
	return func(c *options) {
		c.id = id
	}
}


func (o options) Uuid(uuid string) Option {
	return func(c *options) {
		c.uuid = uuid
	}
}

func (o options) apply(opts ...Option) options {
	ret := options{}
	for _,opt := range opts {
		opt(&ret)
	}
	if ret.id == 0 {
		ret.id = 1
	}
	if ret.uuid == "" {
		ret.uuid = "DEFAULT VALUE"
	}
	return ret
}