package main

import "flag"
import "log"
import "os"
import "path"
import "strings"
import "io"
import "regexp"
import "strconv"
import "os/signal"
import "syscall"
import "github.com/yanke-guo/idy"
import "net/http"

// constants

const POOL_NAME_PATTERN = "^[0-9a-zA-Z._-]+$"

// command line args

var args struct {
	bind    string
	dataDir string
	pools   string
	shard   string
}

// pools

var pools = make(map[string]*idy.Pool)

func main() {
	// command line args
	flag.StringVar(&args.bind, "b", "127.0.0.1:8865", "HTTP service bind address")
	flag.StringVar(&args.dataDir, "d", "data", "location of data directory")
	flag.StringVar(&args.pools, "p", "user", "comma seperated pool names, only "+POOL_NAME_PATTERN+" is allowed")
	flag.StringVar(&args.shard, "s", "1:1:2048:4096", "sharding partten of this instance, see README.md")

	flag.Parse()

	// ensure data directory
	wd, _ := os.Getwd()
	args.dataDir = path.Join(wd, args.dataDir)

	err := os.MkdirAll(args.dataDir, 0700)
	if err != nil {
		log.Fatalln("failed to create data dir", args.dataDir, err)
	}

	// 'pools' validation
	poolNames := strings.Split(args.pools, ",")

	if len(poolNames) == 0 {
		log.Fatalln("pool names not specified")
	}

	for _, name := range poolNames {
		ok, err := regexp.MatchString(POOL_NAME_PATTERN, name)
		if !ok || err != nil {
			log.Fatalln("pool name", name, "illegal")
		}
	}

	// slice configuration
	config := idy.SliceConfig{}

	err = idy.DecodeSliceConfig(args.shard, &config)

	if err != nil {
		log.Fatalln(err)
	}

	// log
	log.Println("idyd started with", args)

	// create pools
	for _, n := range poolNames {
		pools[n] = idy.NewPool(n, args.dataDir, config)
	}

	// run pools
	for _, p := range pools {
		go p.Run()
	}

	// run HTTP
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		coms := strings.Split(req.URL.Path, "/")
		if len(coms) != 3 {
			res.WriteHeader(400)
			io.WriteString(res, "bad url")
		}

		name := coms[1]
		frmt := coms[2]

		pool := pools[name]
		if pool != nil {
			id := pool.NewId()
			if frmt == "_hex" {
				io.WriteString(res, strconv.FormatUint(id, 16))
			} else if frmt == "_dec" {
				io.WriteString(res, strconv.FormatUint(id, 10))
			} else {
				res.WriteHeader(400)
				io.WriteString(res, "fomat not supported"+frmt)
			}
		} else {
			res.WriteHeader(400)
			io.WriteString(res, "not found "+req.URL.Path)
		}
	})

	go http.ListenAndServe(args.bind, nil)

	// wait SIGINT, SIGTERM
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("shuting down due to", sig)
		// shutdown Pool
		for _, p := range pools {
			p.Shutdown()
		}
		done <- true
	}()

	<-done
}
