package main

// only need mysql OR sqlite
// both are included here for reference
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

var db *gorm.DB
var err error

type Category struct {
	ID           uint   `json:"id"`
	Name 	     string `json:"name"`
	Description  string `json:"description"`
}


type BattleField struct {
	Tanks [2]*Tank   `json:"tanks"`
	Range Position `json:"range"`
	Next  uint    `json:"next"`
}

func (game *BattleField) getTank(id uint) *Tank {

	return game.Tanks[id-1] //TODO
}

func (game *BattleField) changeNext() {
	if game.Next==1 {
		game.Next=2
	}else {
		game.Next=1
	}
}

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Result struct {
	Result string `json:"x"`
}

const (
	DIRECTION_UP string  = "UP"
	DIRECTION_DOWN string = "DOWN"
	DIRECTION_RIGHT string="RIGHT"
	DIRECTION_LEFT string="LEFT"

	ENEMY_DAMAGED = "enemy Damaged"
	MISSED string= "Missed"
	OUT_OF_RANGE string= "OUT_OF_RANGE"
	MOVED string= "MOVED"
	NOT_YOUR_TURN string= "NOT_YOUR_TURN"
	YOU_WIN string= "YOU_WIN"

)
type Tank struct {
	Life      int64    `json:"life"`
	Position  Position `json:"Position"`
	Direction string   `json:"Direction"`
	ID        uint    `json:"id"`
}

func (tank *Tank) move(position Position) {
	tank.Position=position
}


func (tank *Tank) shoot(direction string) bool {
	enemy := game.getTank(game.Next)
	enemyPos := enemy.Position
	myPos:=tank.Position
	enemyDamaged := false
	if (direction == DIRECTION_UP  && enemyPos.X==myPos.X && enemyPos.Y>myPos.Y)||
		(direction ==DIRECTION_DOWN  && enemyPos.X==myPos.X && enemyPos.Y<myPos.Y)||
		(direction ==DIRECTION_RIGHT  && enemyPos.Y==myPos.Y && enemyPos.X>myPos.X)||
		(direction ==DIRECTION_LEFT  && enemyPos.Y==myPos.Y && enemyPos.X<myPos.X){
		enemyDamaged = true
	}
	return enemyDamaged

}

var game *BattleField
func mover(id uint, direction string) string {
	thisTank, notYourTurnError := getTankAndUpdateTurn(id)
	if notYourTurnError {
		return NOT_YOUR_TURN
	}

	directions := make(map[string]Position)
	directions[DIRECTION_UP] = Position{X: 0, Y: 1}
	directions[DIRECTION_DOWN] = Position{X: 0, Y: -1}
	directions[DIRECTION_RIGHT] = Position{X: 1, Y: 0}
	directions[DIRECTION_LEFT] = Position{X: -1, Y: 0}


	probableX := thisTank.Position.X + directions[direction].X
	probableY := thisTank.Position.Y + directions[direction].Y
	if probableX > game.Range.X || probableX < 0 || probableY > game.Range.Y || probableY < 0 {
		return OUT_OF_RANGE
	}

	thisTank.move(Position{X: probableX, Y: probableY})
	thisTank.Direction=direction
	return MOVED
}

func shoot(id uint ) string{
	thisTank, notYourTurnError := getTankAndUpdateTurn(id)
	if notYourTurnError {
		return NOT_YOUR_TURN
	}
	enemyDamaged := thisTank.shoot(thisTank.Direction)
	enemy := game.getTank(game.Next)

	if enemyDamaged {
		enemy.Life = enemy.Life-1
		if enemy.Life <1 {
			return YOU_WIN
		}
		return ENEMY_DAMAGED
	}
	return MISSED
}



func getTankAndUpdateTurn(id uint) (*Tank, bool) {
	if game.Next != id {
		return nil, true//TODO
	}
	thisTank := game.getTank(id)
	game.changeNext()
	secondsPlayed = 0
	return thisTank, false
}
var secondsPlayed =0

func main() {

	position := &Position{X: 3000, Y: 3000}
	tanks := [2]*Tank{{Position:Position{X: 1, Y: 1},Direction:DIRECTION_UP,ID:uint(1),Life:3},
	        {Position:Position{X: 1, Y: 10},Direction:DIRECTION_UP,ID:uint(2),Life:3}}
	game = &BattleField{Range: *position , Next:uint(1),Tanks:tanks}

	go func() {
		for {
			time.Sleep(1000)
			secondsPlayed++
			if secondsPlayed == 2 {
				game.changeNext()
			}
		}
	}()

	fmt.Println(shoot(1))
//	time.Sleep(50000)
	fmt.Println(shoot(1))
	fmt.Println(mover(2,DIRECTION_RIGHT))
	fmt.Println(shoot(1))
	fmt.Print(game.Tanks[0])
	fmt.Print(game.Tanks[1])

	db, err = gorm.Open("sqlite3", "./gorm.db"); if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Auto genera la tabla en referencia a la struct
	db.AutoMigrate(&Category{})

	//estos son los endpoints
	r := gin.Default()
	r.GET("/categories/", GetCategories)
	r.GET("/categories/:id", GetCategory)
	r.POST("/categories", CreateCategory)
	r.PUT("/categories/:id", UpdateCategory)
	r.DELETE("/categories/:id", DeleteCategory)

	r.Run(":8080")
}

func DeleteCategory(c *gin.Context) {
	id := c.Params.ByName("id")
	var category Category
	d := db.Where("id = ?", id).Delete(&category)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdateCategory(c *gin.Context) {

	var category Category
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&category)

	db.Save(&category)
	c.JSON(200, category)

}

func CreateCategory(c *gin.Context) {

	var category Category
	c.BindJSON(&category)

	db.Create(&category)
	c.JSON(200, category)
}

//   "/category/:id"
func GetCategory(c *gin.Context) {
	id := c.Params.ByName("id")
	var category Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, category)
	}
}

// crea un array de categorias. Lo pasa por referencia y devuelve not found si no lo encuentra
func GetCategories(c *gin.Context) {
	var category []Category
	if err := db.Find(&category).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, category)
	}

}
func Fibbonacci(n int) int{
	if n < 2{
		 return n
	}
	return Fibbonacci(n-1) + Fibbonacci(n-2)
}
