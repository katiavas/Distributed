package gol



import (
	"uk.ac.bris.cs/gameoflife/stubs"
)



type EngineOperations struct {}


//the function:Which accepts requests of the types that we defined and we do this by importing the interface
func (s *EngineOperations) CalculateNextState(req stubs.Request, res *stubs.Response) (err error) {
	res.W = make([][]byte, req.Param.ImageHeight)
	for i := range res.W {
		res.W[i] = make([]byte, req.Param.ImageWidth)
	}
	res.W = distributor(req.Param, req.W)
	return

}