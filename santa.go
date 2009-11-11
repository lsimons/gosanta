// Copyright 2009 Leo Simons. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Implementation of the santa claus concurrencly problem in Go

import "log"
import "time"

const NUMBER_OF_REINDEER = 9
const REQUIRED_REINDEER = 9
const NUMBER_OF_ELVES = 10
const REQUIRED_ELVES = 4
const NUMBER_OF_RUNS = 10000

//
// Messaging code
//

type Message int
const (
	READY = iota;
	WORK;
	DONE;
	STOP;
	STOPPED;
)

//
// Helper
//

type Helper struct {
	i int;
	t string;
	cmd chan Message;
	status chan Message;
}

func (h Helper) Exhaust() {
	//log.Stdoutf("Helper.Exhaust() -- %v:%v is being listened to\n", h.t, h.i);
	for {
		select {
			case msg := <- h.status:
				switch msg {
					case DONE:
						//log.Stdoutf("Helper.Exhaust() -- %v:%v is done\n", h.t, h.i);
				}
			default:
				return;
		}
	}
}

func (h Helper) Work() {
	//log.Stdoutf("Helper.Work() -- %v:%v is being told to work\n", h.t, h.i);
	h.cmd <- WORK;
}

func (h Helper) Stop() {
	//log.Stdoutf("Helper.Stop() -- %v:%v is being told to stop\n", h.t, h.i);
	h.cmd <- STOP;
}

func (h Helper) IsWaiting() bool {
	select {
		case msg := <- h.status:
			switch msg {
				case READY:
					//log.Stdoutf("Helper.CheckReady() -- %v:%v is ready\n", h.t, h.i);
					return true;
				default:
					//log.Exitf("Helper.CheckReady() -- %v:%v sent message I don't understand: %v", h.t, h.i, msg);
			}
		default:
			//log.Stdoutf("Helper.CheckReady() -- %v:%v didn't say it was ready\n", h.t, h.i);
	}
	return false;
}

func (h Helper) Wait() {
	for {
		select {
			case msg := <- h.status:
				switch msg {
					case DONE:
						//log.Stdoutf("Helper.Join() -- %v:%v is done\n", h.t, h.i);
						return;
					default:
						log.Exitf("Helper.Wait() -- %v:%v sent message I don't understand: %v", h.t, h.i, msg);
				}
		}
	}
}

func (h Helper) Join() {
	for {
		select {
			case msg := <- h.status:
				switch msg {
					case STOPPED:
						//log.Stdoutf("Helper.Wait() -- %v:%v has stopped\n", h.t, h.i);
						return;
					case DONE:
						//log.Stdoutf("Helper.Wait() -- %v:%v is done\n", h.t, h.i);
					default:
						log.Exitf("Helper.Wait() -- %v:%v sent message I don't understand: %v", h.t, h.i, msg);
				}
		}
	}
}

func (h Helper) start() {
	for {
		select {
			case h.status <- READY:
				// we told santa we are ready, wait for a response
			case msg := <- h.cmd:
				// santa told us something
				switch msg {
					case STOP:
						//log.Stdoutf("Helper.start() -- %v:%v is stopping\n", h.t, h.i);
						h.status <- STOPPED;
						return;
					case WORK:
						//log.Stdoutf("Helper.start() -- %v:%v is doing work\n", h.t, h.i);
						h.status <- DONE;
						time.Sleep(1000 * 100);
					default:
						log.Exitf("Helper.start() -- %v:%v received message I don't understand", h.t, h.i, msg);
				}
		}
	}
}

func NewHelper(i int, t string) *Helper {
	h := new(Helper);
	h.i = i;
	h.t = t;
	h.cmd = make(chan Message);
	h.status = make(chan Message);
	go h.start();
	return h;
}

func NewHelpers(num int, t string) []Helper {
	helpers := make([]Helper, num);
	for i := 0; i < len(helpers); i++ {
		helpers[i] = *NewHelper(i, t);
	}
	return helpers;
}

func WorkHelpers(helpers []Helper) {
	for _, h := range helpers { h.Work() }
}
func ExhaustHelpers(helpers []Helper) {
	for _, h := range helpers { h.Exhaust() }
}
func StopHelpers(helpers []Helper) {
	for _, h := range helpers { h.Stop() }
}
func JoinHelpers(helpers []Helper) {
	for _, h := range helpers { h.Join() }
}
func WaitHelpers(helpers []Helper) {
	for _, h := range helpers { h.Wait() }
}

//
// Santa
//

type Santa struct {
	reindeer []Helper;
	elves []Helper;
}

func (s Santa) DeliverPresents(reindeer []Helper) {
	//log.Stdout("Santa is delivering presents!\n");
	WorkHelpers(reindeer);
	WaitHelpers(reindeer);
	//log.Stdout("Santa is done delivering presents!\n");
}

func (s Santa) MakePresents(elves []Helper) {
	//log.Stdout("Santa is making presents!\n");
	WorkHelpers(elves);
	WaitHelpers(elves);
	//log.Stdout("Santa is making presents!\n");
}

func (s Santa) DoWork(numRuns int) {

	var wReindeer []Helper;
	var wReindeerNo int;
	wReindeerReset := func() {
		wReindeer = make([]Helper, REQUIRED_REINDEER);
		wReindeerNo = 0;
	};
	wReindeerReset();

	var wElves []Helper;
	var wElfNo int;
	wElvesReset := func() {
		wElves = make([]Helper, REQUIRED_ELVES);
		wElfNo = 0;
	};
	wElvesReset();
	
	runs := 0;
	checkDone := func() bool {
		runs++;
		return runs >= numRuns;
	};
	DOWORK_LOOP: for {
		for _, r := range s.reindeer {
			time.Sleep(1000 * 100);
			
			if r.IsWaiting() {
				wReindeer[wReindeerNo] = r;
				wReindeerNo++;
			}
			if wReindeerNo == REQUIRED_REINDEER {
				s.DeliverPresents(wReindeer);
				wReindeerReset();
				if checkDone() { break DOWORK_LOOP }
			}
		}


		for _, e := range s.elves {
			time.Sleep(1000 * 100);

			if e.IsWaiting() {
				wElves[wElfNo] = e;
				wElfNo++;
			}
			if wElfNo == REQUIRED_ELVES {
				s.MakePresents(wElves);
				wElvesReset();
				if checkDone() { break DOWORK_LOOP }
			}
		}
	}
	log.Stdoutf("Done running around after %v runs, cleaning up!\n", runs);
	ExhaustHelpers(s.reindeer);
	ExhaustHelpers(s.elves);

	StopHelpers(s.reindeer);
	StopHelpers(s.elves);
	JoinHelpers(s.reindeer);
	JoinHelpers(s.elves);
}

func NewSanta(reindeer []Helper, elves []Helper) *Santa {
	s := new(Santa);
	s.reindeer = reindeer;
	s.elves = elves;
	return s;
}

func main() {
	reindeer := NewHelpers(NUMBER_OF_REINDEER, "reindeer");
	elves := NewHelpers(NUMBER_OF_ELVES, "elves");
	santa := NewSanta(reindeer, elves);
	santa.DoWork(NUMBER_OF_RUNS);
}
