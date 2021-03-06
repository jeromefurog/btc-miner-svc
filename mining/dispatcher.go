package mining

import (
	"github.com/jeromefurog/btc-miner-svc/block"
	"github.com/jeromefurog/btc-miner-svc/logger"
	"runtime"
	"time"
)

var chunk_queue_capacity int = 300
var monitor logger.Logger
var Psize = poolsize()

//Dispatcher Entity.
//Contains a Pool of chans to send and receive from other miners.
//A queue of chunks to mine
//And a queue of chunks to validate and submit
type Dispatcher struct {
	MiningPool    chan chan Chunk
	ChunkQueueIn  chan Chunk
	ChunkQueueOut chan Chunk
}

//Make new Dispatcher
func NewDispatcher(log logger.Logger) *Dispatcher {
	pool := make(chan chan Chunk, Psize)
	chunkqueuein := make(chan Chunk, chunk_queue_capacity)
	chunkqueueout := make(chan Chunk, chunk_queue_capacity)
	monitor = log
	monitor.Print("info", "New Dispatcher created")
	return &Dispatcher{MiningPool: pool, ChunkQueueIn: chunkqueuein, ChunkQueueOut: chunkqueueout}
}

//Start the new dispatcher, create the miners, start them and begin dispatching.
func (dispatcher *Dispatcher) Run() {
	for i := 0; i < cap(dispatcher.MiningPool); i++ {
		NewMiner(i, dispatcher.MiningPool, dispatcher.ChunkQueueOut).Start()
		monitor.Print("info", "New Miner added to the pool")
	}
	go dispatcher.Dispatch()
}

//Dispatcher start the counter for monitoring. Waits for chunk and send it to an available miner
func (dispatcher *Dispatcher) Dispatch() {
	monitor.Print("info", "Starting time counter")
	monitor.BeginTime = time.Now()
	for {
		monitor.Print("info", "waiting for a new chunk")
		select {
		case chunk := <-dispatcher.ChunkQueueIn:
			AvailableMiner := <-dispatcher.MiningPool
			AvailableMiner <- chunk
		case chunk := <-dispatcher.ChunkQueueOut:
			//Verify received Chunk and must be sent back to the Bitcoin Core client through Websocket
			if verifyChunk(chunk) {
				chunk.Valid = true
			}
		}
	}
}

//Verify given chunk. To be completed with more checks related to Bitcoin.
func verifyChunk(chunk Chunk) bool {
	if hash := block.Doublesha256_BlockHeader(chunk.Block); hash < chunk.Target {
		return true
	}
	return false
}

//Set the number of miners depending on the number of threads of the machine.
//Made to optimize and reduce the overhead on multiplex scheduling
func poolsize() int {
	switch maxprocs := runtime.GOMAXPROCS(0); maxprocs {
	case 1:
		return 1
	case 2:
		return 1
	case 3:
		return 2
	default:
		return maxprocs - 2
	}
}
