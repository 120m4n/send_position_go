package singleton

var instance *Singleton

type Singleton struct {
	Fleet  string
	Userid string
	Uniqueid string
}

func GetInstance() *Singleton {
	if instance == nil {
		instance = &Singleton{}
	}
	return instance
}

func (s *Singleton) SetFleet(fleet string) {
	s.Fleet = fleet
}

func (s *Singleton) SetUserid(userid string) {
	s.Userid = userid
}

func (s *Singleton) SetUniqueid(uniqueid string) {
	s.Uniqueid = uniqueid
}