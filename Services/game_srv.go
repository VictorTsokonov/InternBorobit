package Services

import (
	"InternBorobitApp/Interfaces"
	"InternBorobitApp/Models"
)

type GameService struct {
	repository Interfaces.GameRepository
}

func NewGameService(repo Interfaces.GameRepository) *GameService {
	return &GameService{repository: repo}
}

func (s *GameService) CreateGame(game *Models.Game) error {
	return s.repository.Create(game)
}

func (s *GameService) GetGameByID(id string) (*Models.Game, error) {
	return s.repository.GetByID(id)
}

func (s *GameService) UpdateGame(game *Models.Game) error {
	return s.repository.Update(game)
}

func (s *GameService) DeleteGame(id string) error {
	return s.repository.Delete(id)
}

func (s *GameService) ListGames() ([]Models.Game, error) {
	return s.repository.List()
}
