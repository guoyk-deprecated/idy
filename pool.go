package idy

import "log"
import "os"
import "path"
import "strconv"
import "time"

type Pool struct {
	name string

	filePath string
	file     *os.File

	slice *Slice

	ticker *time.Ticker

	requests chan chan uint64
	shutdown chan chan bool
}

func NewPool(name string, dataDir string, c SliceConfig) *Pool {
	p := &Pool{
		name: name,

		filePath: path.Join(dataDir, name+".json"),
		file:     nil,

		slice: nil,

		requests: make(chan chan uint64, 32),
		shutdown: make(chan chan bool, 1),
	}
	p.setup(c)
	return p
}

func (p *Pool) NewId() uint64 {
	req := make(chan uint64, 1)
	p.requests <- req
	return <-req
}

func (p *Pool) setup(c SliceConfig) {
	var err error
	// open or create database file
	p.file, err = os.OpenFile(p.filePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln("failed to create database file", p.name, err)
	}
	// check file size
	stat, err := p.file.Stat()
	if err != nil {
		log.Fatalln("failed to read database file", p.name, err)
	}
	if stat.Size() == 0 {
		// init database file if file is empty
		p.initDatabase(c)
		log.Println("database ["+p.name+"] created at", p.filePath)
	} else {
		// load database file if file is not empty
		p.loadDatabase()
		log.Println("database ["+p.name+"] loaded from", p.filePath)
	}
}

func (p *Pool) Run() {
	p.ticker = time.NewTicker(5 * time.Second)
	for {
		select {
		case req := <-p.requests:
			{
				// get the id from slice
				id, moved := p.slice.NextId()
				// save database if slice moved
				if moved {
					p.saveDatabase()
				}
				// send back new id
				req <- id
			}
		case <-p.ticker.C:
			{
				// ticker fired
				p.saveDatabase()
			}
		case done := <-p.shutdown:
			{
				// save the database
				p.saveDatabase()
				// notify done
				done <- true
			}
		}
	}
	p.ticker.Stop()
	p.ticker = nil
}

func (p *Pool) initDatabase(c SliceConfig) {
	// create and init the slice
	p.slice = NewSlice(c)
	p.slice.NewSeed()
	p.slice.UpdateElements()
	// save database
	p.saveDatabase()
}

func (p *Pool) saveDatabase() {
	// seek to 0
	_, err := p.file.Seek(0, 0)
	if err != nil {
		log.Fatalln("failed to read database file", p.name, err)
	}

	// export database
	db := p.slice.toDatabase()

	// write database to file
	err = db.Encode(p.file)
	if err != nil {
		log.Fatalln("failed to write database file", p.name, err)
	}

	// get current cursor
	c, err := p.file.Seek(0, 1)
	if err != nil {
		log.Fatalln("failed to write database file", p.name, err)
	}

	// truncate file
	err = p.file.Truncate(c + 1)
	if err != nil {
		log.Fatalln("failed to write database file", p.name, err)
	}

	// sync file
	err = p.file.Sync()
	if err != nil {
		log.Fatalln("failed to write database file", p.name, err)
	}
}

func (p *Pool) loadDatabase() {
	// seek to 0
	_, err := p.file.Seek(0, 0)
	if err != nil {
		log.Fatalln("failed to read database file", p.name, err)
	}

	// unmarshal
	db := Database{}

	err = DecodeDatabase(&db, p.file)
	if err != nil {
		log.Fatalln("failed to read database file")
	}

	if db.Version != 1 {
		log.Fatalln("failed to read database file, version is not 1")
	}

	// decode slice config
	c := SliceConfig{}

	err = DecodeSliceConfig(db.Shard, &c)

	if err != nil {
		log.Fatalln("failed to read database file,", err)
	}

	// init slice
	p.slice = NewSlice(c)

	// load seed
	p.slice.Seed, err = strconv.ParseInt(db.Seed, 10, 64)
	if err != nil || p.slice.Seed == 0 {
		log.Fatalln("failed to read database file, could not load seed", err)
	}

	// load start
	p.slice.Start, err = strconv.ParseUint(db.Start, 10, 64)
	if err != nil {
		log.Fatalln("failed to read database file, could not load start", err)
	}

	// load index
	p.slice.Index = db.Index

	if p.slice.Index >= p.slice.Config.SliceEffectiveSize {
		log.Fatalln("failed to read database file, index exceeded")
	}

	// update elements
	p.slice.UpdateElements()
}

func (p *Pool) Shutdown() {
	done := make(chan bool, 1)
	p.shutdown <- done
	<-done
}
