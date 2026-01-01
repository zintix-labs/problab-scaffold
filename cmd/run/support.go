// Copyright 2025 Zintix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"github.com/zintix-labs/problab-scaffold/pkg/engine"
	"github.com/zintix-labs/problab/spec"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var cfg *config = new(config)

type config struct {
	name      string
	id        spec.GID
	worker    int
	player    int
	bets      int
	spins     int
	betMode   int
	seed      int64
	pprofmode string
}

type gidFlag struct{ p *spec.GID }

func (f gidFlag) String() string { return fmt.Sprint(uint(*f.p)) }
func (f gidFlag) Set(s string) error {
	u, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return err
	}
	*f.p = spec.GID(uint(u))
	return nil
}

func bindVar() {
	flag.Var(gidFlag{&cfg.id}, "game", "target game id")
	flag.IntVar(&cfg.worker, "worker", 1, "number of workers")
	flag.IntVar(&cfg.player, "player", 1, "number of players")
	flag.IntVar(&cfg.bets, "bets", 200, "initial bets")
	flag.IntVar(&cfg.spins, "spins", 10000000, "spins per player")
	flag.IntVar(&cfg.betMode, "mode", 0, "bet mode index")
	flag.Int64Var(&cfg.seed, "seed", -1, "int64 seed for random number generator")
	flag.StringVar(&cfg.pprofmode, "p", "", "pprof: '', cpu, heap, allocs")

	flag.Parse()

	// given seed illeagel -> default seed
	if cfg.seed < 1 {
		seed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			log.Fatal(err)
		}
		cfg.seed = seed.Int64()
	}
}

func executeSimulator() {
	cfg.valid()

	lab := engine.MustNew()

	s, err := lab.NewSimulatorWithSeed(cfg.id, cfg.seed)
	if err != nil {
		log.Fatal(err)
	}
	ent, _ := lab.EntryById(cfg.id)
	cfg.name = ent.Name
	// able to execute
	green := "\033[1;32m"
	reset := "\033[0m"
	p := message.NewPrinter(language.English)

	if cfg.player == 1 { // sim machine
		if cfg.worker == 1 {
			// pure sim singo core
			p.Printf("%s[GAME:%s] [PLAYMODE:%d] [SPINS:%d]%s\n", green, cfg.name, cfg.betMode, cfg.spins, reset)
			st, used, _ := s.Sim(cfg.betMode, cfg.spins, true)
			st.StdOut(used)
		} else {
			// pure sim multi core
			p.Printf("%s[WORKERS:%d] [GAME:%s] [PLAYMODE:%d] [SPINS:%d]%s\n", green, cfg.worker, cfg.name, cfg.betMode, cfg.worker*cfg.spins, reset)
			st, used, _ := s.SimMP(cfg.betMode, cfg.spins, cfg.worker, true) // 併發
			st.StdOut(used)
		}
	} else {
		// sim by player's experenece statemant
		p.Printf("%s[WORKERS:%d] [GAME:%s] [PLAYERS:%d BALANCE:%d PLAYMODE:%d SPINS:%d]%s\n", green, cfg.worker, cfg.name, cfg.player, cfg.bets, cfg.betMode, cfg.spins, reset)
		st, est, used, _ := s.SimPlayers(cfg.worker, cfg.player, cfg.bets, cfg.betMode, cfg.spins, true)
		st.StdOut(used)
		est.Out()
	}
}

func (cfg *config) valid() {
	p := message.NewPrinter(language.English)

	if cfg.worker < 1 {
		log.Fatal("value err : workers must > 0")
	}

	if cfg.player < 1 {
		log.Fatal("value err : player must > 0")
	}
	// resize players
	if cfg.player > 100000 {
		p.Printf("too much players: %d resized to 100k players\n", cfg.player)
		cfg.player = 100000
	}

	if cfg.player > 1 && cfg.bets < 1 {
		log.Fatal("value err : balance must >= 1")
	}

	if cfg.spins < 1 {
		log.Fatal("value err : spins must > 0")
	}

	// When simulating player-based sessions, cap spins per player to 15,000.
	// This is an intentional business constraint rather than a technical limit.
	//
	// Even under aggressive turbo play (≈2–3 seconds per spin), a single player
	// can only reach ~1,200–1,800 spins per hour. 15,000 spins already represent
	// more than 10 hours of continuous high-intensity gameplay, which exceeds
	// what can be reasonably modeled as a short-term player experience.
	// For long-horizon behavior, machine-level simulation should be used instead.
	if cfg.player > 1 && cfg.spins > 15000 {
		p.Printf("too much spins for each players : %d resized to 15k spins for each player\n", cfg.spins)
		cfg.spins = 15000
	}
}
