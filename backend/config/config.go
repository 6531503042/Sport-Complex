package config

type (

	Config struct {
		App App
		Db   Db
		Grpc Grpc
	}

	App struct {
		Name string
		Url string
		Stage string
	}

	Db struct {
		Url string
	}

	Grpc struct {
		AuthUrl string
		UserUrl string
		GymUrl string
		BadmintonUrl string
		SwimmingUrl string
		FootballUrl string
		PaymentUrl string
	}
)