package chip

type dummyLogger struct{}

func (d *dummyLogger) Errorf(string, ...interface{}) {}
func (d *dummyLogger) Warnf(string, ...interface{})  {}
func (d *dummyLogger) Infof(string, ...interface{})  {}
func (d *dummyLogger) Debugf(string, ...interface{}) {}
