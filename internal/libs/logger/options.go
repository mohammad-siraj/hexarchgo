package logger

type IConfigOptionsFunction func(*configs)

type Iconfigs interface {
	WithFilename(filename string)
	WithMaxBackUp(maxBackUp int)
	WithMaxAge(maxAge int)
	WithIsLocalTime(isLocalTime bool)
	WithIsCompressed(isCompressed bool)
	IsProduction() bool
	FileName() string
	MaxAge() int
	MaxSize() int
	MaxBackups() int
	IslocalTime() bool
	IsCompressed() bool
}

type configs struct {
	fileName     string `default:"../test.log"`
	maxSize      int    `default:"1024"`  //max size in bytes
	maxBackups   int    `default:"30"`    //max number of backup files
	maxAge       int    `default:"90"`    //age in days
	localTime    bool   `default:"true"`  //time used for local time
	compress     bool   `default:"true"`  //should the data be compressed
	isProduction bool   `default:"false"` //if the logger is set for production
	logRotation  bool   `default:"false"` //log rotation enabler
}

func NewlogConfigOptions(isProduction bool, opts ...IConfigOptionsFunction) Iconfigs {

	config := &configs{
		isProduction: isProduction,
		maxSize:      1024,
		maxBackups:   30,
		maxAge:       90,
		localTime:    true,
		compress:     true,
	}

	for _, operation := range opts {
		operation(config)
	}

	return config
}

func (c *configs) WithFilename(filename string) {
	c.fileName = filename
}

func (c *configs) WithMaxSize(maxSize int) {
	c.maxSize = maxSize
}

func (c *configs) WithMaxBackUp(maxBackUp int) {
	c.maxBackups = maxBackUp
}
func (c *configs) WithMaxAge(maxAge int) {
	c.maxAge = maxAge
}

func (c *configs) WithIsLocalTime(isLocalTime bool) {
	c.localTime = isLocalTime
}

func (c *configs) WithIsCompressed(isCompressed bool) {
	c.compress = isCompressed
}

func (c *configs) IsProduction() bool {
	return c.isProduction
}
func (c *configs) FileName() string {
	return c.fileName
}

func (c *configs) MaxAge() int {
	return c.maxAge
}
func (c *configs) MaxSize() int {
	return c.maxSize
}
func (c *configs) MaxBackups() int {
	return c.maxBackups
}
func (c *configs) IslocalTime() bool {
	return c.localTime
}

func (c *configs) IsCompressed() bool {
	return c.compress
}

func (c *configs) IsLogRotated() bool {
	return c.logRotation
}
