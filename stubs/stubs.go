
package stubs

import (
"uk.ac.bris.cs/gameoflife/util"
)

var NextState = "EngineOperations.CalculateNextState"
//var Par = "EngineOperations.AcceptP"
var Distributor = "ControllerOperations.Controller"


type Response struct {
	W [][]byte
	Param Parameters
}

type Request struct {
	W [][]byte
	Param Parameters
}

type Parameters struct{
	Turns int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

type ReqC struct {
	CompletedTurns int
	Alive []util.Cell

}

type ResC struct{
	CompletedTurns int
	Alive []util.Cell
}

/*type Respons struct {
	Message string
}

type Reques struct {
	Message string
}*/
