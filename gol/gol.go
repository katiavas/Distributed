package gol

type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}
// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {
	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	//create a filename channel
	ioFilename := make(chan string)
	ioInput := make(chan uint8)
	ioOutput := make(chan uint8)

	c := distributorChannels{
		events,
		ioCommand,
		ioIdle,
		ioFilename,
		ioInput,
		ioOutput,
	}
		go client(p, c)


	//go distributor(p, c)
	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: ioFilename,
		output:   ioOutput,
		input:    ioInput,
	}

	/*	server := flag.String("server","127.0.0.1:8030","IP:port string to connect to as server")
		flag.Parse()
		//create an rcp client/Its going to dial the server address the server has passed in the command line
		client1, _ := rpc.Dial("tcp", *server)
		defer client1.Close()
		makeCall(*client1,p,c)*/
	//form a filename
	//%dx%d bc images follow this specific format (e.g 16x16 ,512x512) and %d for base10
	go startIo(p, ioChannels)
	//send the filename, i want to read , down the ioFilename channel

}
