package Interfaces

import "InternBorobitApp/Models"

type GameRepository interface {
	Create(game *Models.Game) error
	GetByID(id string) (*Models.Game, error)
	Update(game *Models.Game) error
	Delete(id string) error
	List() ([]Models.Game, error) // []*Models.Game, pointer not required
}
