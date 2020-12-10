package gol

import (
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events    chan<- Event
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool
	ioFilename chan <- string
	ioInput <- chan uint8
	ioOutput chan<- uint8
}


const alive = 255
const dead = 0

func mod(x, m int) int {
	return (x + m) % m
}

func calculateNeighbours(p stubs.Parameters, x, y int, world [][]byte) int {
	neighbours := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i != 0 || j != 0 {
				if world[mod(y+i, p.ImageHeight)][mod(x+j, p.ImageWidth)] == alive {
					neighbours++
				}
			}
		}
	}
	return neighbours
}


func calculateAliveCells(p stubs.Parameters, world [][]byte) []util.Cell {
	aliveCells := []util.Cell{}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 255 {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}

	return aliveCells
}

func calculateNextState(p stubs.Parameters, world [][]byte) [][]byte {
	newWorld := make([][]byte, p.ImageHeight)
	for i := range newWorld {
		newWorld[i] = make([]byte, p.ImageWidth)
	}
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			neighbours := calculateNeighbours(p, x, y, world)
			if world[y][x] == alive {
				if neighbours == 2 || neighbours == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead

				}
			} else {
				if neighbours == 3 {
					newWorld[y][x] = alive
				} else {
					newWorld[y][x] = dead
				}
			}
		}
	}
	return newWorld
}



// distributor divides the work between workers and interacts with other goroutines.
func distributor(p stubs.Parameters, newWorld [][]byte) [][]byte{

	// TODO: Create a 2D slice to store the world.
	World := make([][]byte, p.ImageHeight)
	for i := range World {
		World[i] = make([]byte, p.ImageWidth)
	}



	turn:=0
	if p.Turns != 0 {
		// Created another 2D slice to store the world that has cache.
		World := make([][]byte, p.ImageHeight)
		for i := range World {
			World[i] = make([]byte, p.ImageWidth)
		}
		for turn = 1; turn <= p.Turns; turn++ {
			World = calculateNextState(p, newWorld)

			for y := 0; y < p.ImageHeight; y++ {
				for x := 0; x < p.ImageWidth; x++ {
					// Replace placeholder World[y][x] with the real newWorld[y][x]
					newWorld[y][x] = World[y][x]
				}
			}
		}
	}

	return newWorld
}

