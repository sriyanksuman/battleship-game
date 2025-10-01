package service

import (
	"battleship-game/enums"
	"battleship-game/models"
	"battleship-game/strategy"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type GameInput struct {
	N              int                    `json:"n"`
	Players        int                    `json:"players"`
	PlayersToShips map[string][]ShipInput `json:"playersToShips"`
	FiringStrategy string                 `json:"firingStrategy"`
}

type ShipInput struct {
	Size   int   `json:"size"`
	Center []int `json:"center"`
}

type GameService struct {
	Battlefield    *models.Battlefield
	Players        []*models.Player
	PlayersMap     map[int]*models.Player
	FiringStrategy strategy.FiringStrategy
	currentIndex   int
}

func NewGameService() *GameService {
	return &GameService{
		Players:    make([]*models.Player, 0),
		PlayersMap: make(map[int]*models.Player),
	}
}

func (gs *GameService) InitGameWithSize(n int) error {
	if n <= 0 {
		return fmt.Errorf("battlefield size must be positive")
	}

	gs.Battlefield = models.NewBattlefield(n)

	for i := 0; i < 2; i++ {
		playerID := i + 1
		p := models.NewPlayer(playerID, fmt.Sprintf("Player%d", playerID))
		gs.Players = append(gs.Players, p)
		gs.PlayersMap[playerID] = p
	}

	gs.divideGridAmongPlayers(n)

	fmt.Printf("Game initialized with %dx%d battlefield\n", n, n)
	return nil
}

func (gs *GameService) AddShipForPlayer(shipID string, size, x, y, playerID int) error {
	if gs.Battlefield == nil {
		return fmt.Errorf("game not initialized. Call InitGame first")
	}

	center := models.Point{X: x, Y: y}

	if !models.ValidateBattleship(gs.Battlefield.N, size, center) {
		return fmt.Errorf("invalid ship placement: center at (%d,%d) with size %d exceeds battlefield boundaries", x, y, size)
	}

	player, exists := gs.PlayersMap[playerID]
	if !exists {
		return fmt.Errorf("player %d does not exist", playerID)
	}

	if !gs.isWithinPlayersHalf(playerID, center, size) {
		return fmt.Errorf("ship placement center (%d,%d) size %d not within Player%d's half", x, y, size, playerID)
	}

	bs := models.NewBattleship(shipID, size, center)
	if err := gs.addBattleshipInGrid(bs); err != nil {
		return err
	}
	player.AddBattleship(bs)

	fmt.Printf("Ship %s (size %d) added to Player%d at position (%d,%d)\n", shipID, size, playerID, x, y)
	return nil
}

func (gs *GameService) StartGame(firingStrategyType string) error {
	if gs.Battlefield == nil {
		return fmt.Errorf("game not initialized. Call InitGame first")
	}

	for _, player := range gs.Players {
		if len(player.Battleships) == 0 {
			return fmt.Errorf("player %d has no ships, add ships before starting the game", player.ID)
		}
	}

	allPoints := make([]models.Point, 0)
	for i := 0; i < gs.Battlefield.N; i++ {
		for j := 0; j < gs.Battlefield.N; j++ {
			allPoints = append(allPoints, models.Point{X: i, Y: j})
		}
	}

	if firingStrategyType == string(enums.RandomFireStrategy) || firingStrategyType == "" {
		gs.FiringStrategy = strategy.NewRandomFiringStrategy(allPoints)
	} else {
		return fmt.Errorf("unsupported firing strategy: %s", firingStrategyType)
	}

	gs.PlayGame()
	return nil
}

func (gs *GameService) ViewBattlefield() error {
	if gs.Battlefield == nil {
		return fmt.Errorf("game not initialized. Call InitGame first")
	}

	gs.Battlefield.PrintBattlefield()
	return nil
}

func (gs *GameService) ReadInput(inputPath string) (*GameInput, error) {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	var input GameInput
	err = json.Unmarshal(data, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func (gs *GameService) InitializeStrategy(input *GameInput) error {
	allPoints := make([]models.Point, 0)
	for i := 0; i < gs.Battlefield.N; i++ {
		for j := 0; j < gs.Battlefield.N; j++ {
			allPoints = append(allPoints, models.Point{X: i, Y: j})
		}
	}

	if input.FiringStrategy == string(enums.RandomFireStrategy) {
		gs.FiringStrategy = strategy.NewRandomFiringStrategy(allPoints)
		return nil
	}

	return fmt.Errorf("incorrect firing strategy: %s", input.FiringStrategy)
}

func (gs *GameService) InitGame(path string) error {
	input, err := gs.ReadInput(path)
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	gs.Battlefield = models.NewBattlefield(input.N)

	for i := 0; i < input.Players; i++ {
		playerID := i + 1
		p := models.NewPlayer(playerID, fmt.Sprintf("Player%d", playerID))
		gs.Players = append(gs.Players, p)
		gs.PlayersMap[playerID] = p
	}

	for playerIDStr, ships := range input.PlayersToShips {
		var playerID int
		fmt.Sscanf(playerIDStr, "%d", &playerID)

		for shipIdx, ship := range ships {
			center := models.Point{X: ship.Center[0] + ship.Size/2, Y: ship.Center[1] + ship.Size/2}
			if models.ValidateBattleship(input.N, ship.Size, center) {
				bs := models.NewBattleship(fmt.Sprintf("SH%d", shipIdx+1), ship.Size, center)
				player := gs.PlayersMap[playerID]
				if err := gs.addBattleshipInGrid(bs); err != nil {
					return fmt.Errorf("overlap while placing player %d ship %s: %v", playerID, bs.ID, err)
				}
				player.AddBattleship(bs)
			} else {
				return fmt.Errorf("ship input invalid for player %d, ship %d", playerID, shipIdx+1)
			}
		}
	}

	gs.divideGridAmongPlayers(input.N)

	err = gs.InitializeStrategy(input)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) addBattleshipInGrid(battleship *models.Battleship) error {
	for point := range battleship.ShipPoints {
		gridPoint := gs.Battlefield.Grid[point.X][point.Y]
		if gridPoint.AssignedObject != nil {
			return fmt.Errorf("ship %s overlaps at (%d,%d)", battleship.ID, point.X, point.Y)
		}
	}
	for point := range battleship.ShipPoints {
		gridPoint := gs.Battlefield.Grid[point.X][point.Y]
		gridPoint.AssignedObject = battleship
	}
	return nil
}

func (gs *GameService) divideGridAmongPlayers(n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			point := gs.Battlefield.Grid[i][j]
			if j < n/2 {
				playerID := gs.Players[0].ID
				point.AssignedPlayer = &playerID
			} else {
				playerID := gs.Players[1].ID
				point.AssignedPlayer = &playerID
			}
		}
	}
}

func (gs *GameService) CheckGameEnd() bool {
	for _, player := range gs.Players {
		if len(player.ActiveBattleships) == 0 {
			fmt.Printf("Player%d has lost all battleships. Ending the game\n", player.ID)
			return true
		}
	}
	return false
}

func (gs *GameService) PlayGame() {
	fmt.Println("\n=== Starting Battleship Game ===")
	gs.Battlefield.PrintBattlefield()

	previousPoint := models.Point{X: -1, Y: -1}

	for !gs.CheckGameEnd() {
		currentPlayer := gs.Players[gs.currentIndex]
		firingPosition := gs.FiringStrategy.Fire(currentPlayer, gs.Battlefield, previousPoint)

		previousPoint = models.Point{X: firingPosition.X, Y: firingPosition.Y}

		if firingPosition.AssignedObject != nil {
			fmt.Printf("Player%d's turn: Missile fired at (%d, %d) : \"Hit\" : %s destroyed",
				currentPlayer.ID, firingPosition.X, firingPosition.Y, firingPosition.AssignedObject.ID)

			destroyedShip := firingPosition.AssignedObject
			gs.dfs(firingPosition.X, firingPosition.Y, destroyedShip)

			player := gs.PlayersMap[*firingPosition.AssignedPlayer]
			delete(player.ActiveBattleships, destroyedShip.ID)
			firingPosition.AssignedObject = nil
		} else {
			fmt.Printf("Player%d's turn: Missile fired at (%d, %d) : \"Miss\"",
				currentPlayer.ID, firingPosition.X, firingPosition.Y)
		}

		fmt.Printf(" : Ships Remaining - PlayerA:%d, PlayerB:%d\n",
			len(gs.PlayersMap[1].ActiveBattleships), len(gs.PlayersMap[2].ActiveBattleships))

		gs.currentIndex = (gs.currentIndex + 1) % len(gs.Players)
	}

	for _, player := range gs.Players {
		if len(player.ActiveBattleships) > 0 {
			fmt.Printf("\nPlayer%d is the WINNER! \n", player.ID)
		}
	}
}

func (gs *GameService) dfs(x, y int, ship *models.Battleship) {
	if x < 0 || y < 0 || x >= gs.Battlefield.N || y >= gs.Battlefield.N {
		return
	}

	point := gs.Battlefield.Grid[x][y]
	if point.AssignedObject == nil || point.AssignedObject.ID != ship.ID {
		return
	}

	point.AssignedObject = nil
	gs.dfs(x+1, y, ship)
	gs.dfs(x-1, y, ship)
	gs.dfs(x, y+1, ship)
	gs.dfs(x, y-1, ship)
}

func (gs *GameService) InitGameSize(n int) error { // alias to clear name
	return gs.InitGameWithSize(n)
}

func (gs *GameService) AddShip(id string, size, xA, yA, xB, yB int) error {
	if gs.Battlefield == nil {
		return fmt.Errorf("game not initialized. Call InitGame first")
	}
	convertCenter := func(x, y int) (int, int) {
		half := size / 2
		row := (gs.Battlefield.N - 1 - y) + half
		col := x
		return row, col
	}
	ax, ay := convertCenter(xA, yA)
	bx, by := convertCenter(xB, yB)

	if err := gs.AddShipForPlayer("A-"+id, size, ax, ay, 1); err != nil {
		return err
	}
	if err := gs.AddShipForPlayer("B-"+id, size, bx, by, 2); err != nil {
		return err
	}
	return nil
}

func (gs *GameService) Start() error {
	return gs.StartGame(string(enums.RandomFireStrategy))
}

func (gs *GameService) isWithinPlayersHalf(playerID int, center models.Point, size int) bool {
	half := size / 2
	topLeftX := center.X - half
	topLeftY := center.Y - half
	for i := topLeftX; i < topLeftX+size; i++ {
		for j := topLeftY; j < topLeftY+size; j++ {
			if j < 0 || j >= gs.Battlefield.N {
				return false
			}
			if playerID == 1 {
				if j >= gs.Battlefield.N/2 {
					return false
				}
			} else {
				if j < gs.Battlefield.N/2 {
					return false
				}
			}
		}
	}
	return true
}
