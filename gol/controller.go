package gol

import (
	"fmt"
	"log"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

// Params provides the details of how to run the Game of Life and which image to load.

//query the game logic engine for the information it needs to respond to the test.



// Params provides the details of how to run the Game of Life and which image to load.

//query the game logic engine for the information it needs to respond to the test.

func readImage(p Params, c distributorChannels)[][]byte {
	newWorld := make([][]byte, p.ImageHeight)
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	// Request the io goroutine to read in the image with the given filename.
	c.ioCommand <- ioInput
	filename := fmt.Sprintf("%dx%d",p.ImageHeight,p.ImageWidth)
	c.ioFilename <- filename
	// The io goroutine sends the requested image byte by byte, in rows.
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			val := <-c.ioInput
			if val != 0 {
				newWorld[y][x] = val
			}
		}
	}
	//fmt.Println(newWorld)
	//World = newWorld
	return newWorld
}

var (
	MyServer string
)

func client(p Params,c distributorChannels) {
	var err error

	//servr:= flag.Lookup("servr").Value.(flag.Getter).Get().(string)
	//servr := flag.String("servr","127.0.0.1:8030","IP:port string to connect to as server")
	//flag.Parse()
	//create an rcp client/Its going to dial the server address the server has passed in the command line
	client, err := rpc.Dial("tcp", "127.0.0.1:8030")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer client.Close()

	World := readImage(p, c)
	if p.Turns == 0 {
		var final []util.Cell
		for y := 0; y < p.ImageHeight; y++ {
			for x := 0; x < p.ImageWidth; x++ {
				if World[y][x] == 255 {
					final = append(final, util.Cell{X: x, Y: y})
				}
			}
		}
		c.events <- FinalTurnComplete{CompletedTurns: p.Turns, Alive: final}

	}

	r := stubs.Parameters{Turns: p.Turns, Threads: p.Threads, ImageWidth: p.ImageWidth, ImageHeight: p.ImageHeight}
	request := stubs.Request{W: World, Param: r}
	response := new(stubs.Response)
	//client calls a procedure to send a request to the server
	err = client.Call(stubs.NextState, request, response)
	if err != nil {
		log.Fatal(err)
	}

	if p.Turns != 0 {
		var final []util.Cell
		for y := 0; y < p.ImageHeight; y++ {
			for x := 0; x < p.ImageWidth; x++ {
				if response.W[y][x] == 255 {
					final = append(final, util.Cell{X: x, Y: y})
				}
			}
		}
		//final = calculateAliveCells(, newWorld)
		c.events <- FinalTurnComplete{CompletedTurns: p.Turns, Alive: final}
	}
	c.ioCommand <- ioOutput
	filename := fmt.Sprintf("%dx%d",p.ImageHeight,p.ImageWidth)
	c.ioFilename <- filename
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			out := response.W[y][x]
			c.ioOutput <- out
		}
	}
	c.events <- ImageOutputComplete{CompletedTurns: p.Turns, Filename: filename}
	c.ioCommand <- ioCheckIdle
	<- c.ioIdle
	c.events <- StateChange{p.Turns, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}