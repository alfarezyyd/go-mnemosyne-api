package discord

type ServiceImpl struct {
	discordRepository Repository
}

func NewService(discordRepository Repository) *ServiceImpl {
	return &ServiceImpl{
		discordRepository: discordRepository,
	}
}
