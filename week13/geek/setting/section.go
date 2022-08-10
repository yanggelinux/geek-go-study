package setting

type APPSettingS struct {
	Name        string
	LogPath     string
	LogLevel    string
	Development bool
}

type ServerSettingS struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  int
	WriteTimeout int
}

type MySQLSettingS struct {
	Host     string
	Username string
	Password string
	Port     int
	DBName   string
}

type JWTSettingS struct {
	JwtSecret string
	JwtExpire int
	JwtIssuer string
}

type KafkaSettingS struct {
	ConsumerBrokers string
	ConsumerTopic   string
	ConsumerGroupID string
	ConsumerVersion string
	ConsumerOffset  string
	ProducerBrokers string
	ProducerTopic   string
	ProducerVersion string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
