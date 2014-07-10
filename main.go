package main

import "fmt"
import "runtime"

/*
The Ackermann function. See http://en.wikipedia.org/wiki/Ackermann_function

@param m unint64 magic!
@param n unint64 also magic!
@param calls *unint64 for keeping track of recursion

@return unint64 much magic

Does magic with `m` and `n`. Also takes `calls` to count the number of times
it recurses, Should be pointer to an unint64 containing 0 when called by
a human
*/
func ackermann(m, n uint64, calls *uint64) uint64 {
	// dereference calls and increment its value
	(*calls) = (*calls) + 1

	// GNDN
	if m == 0 {
		return n + 1
	} else if n == 0 {
		return ackermann(m-1, 1, calls)
	} else {
		return ackermann(m-1, ackermann(m, n-1, calls), calls)
	}
}

/*
Wrapper for calling ackermann()

@param m unint64 magic!
@param n unint64 also magic!
@param chan ackRes channel for sending the 'return'

Takes `m` and `n` and pass those into ackermann().
Sets calls to zero and passes that in too.
Packages the results of other info an ackRes and sends that down the channel.
*/
func caller(m, n uint64, ch chan ackRes) {
	calls := uint64(0)
	ch <- ackRes{m, n, ackermann(m, n, &calls), calls}
}

/*
Struct for holding Ackermann function results and meta-data.

@field m unint64 m used to invoke ackermann()
@field n unint64 n used to invoke ackermann()
@field res unint64 ackermann result
@field calls unint64 number of times ackermann() recursed
*/
type ackRes struct {
	m     uint64
	n     uint64
	res   uint64
	calls uint64
}

func main() {
	// @todo set these via cmd args
	runtime.GOMAXPROCS(1)
	mmax := uint64(5)
	nmax := uint64(5)

	// make the channel for capturing the results
	ch := make(chan ackRes)
	// compute the number of times we will call ackermann()
	numAcks := (mmax + 1) * (nmax + 1)
	// call ackermann() with all combination of 0 through m and n
	for m := uint64(0); m <= mmax; m++ {
		for n := uint64(0); n <= nmax; n++ {
			// launch ackermann() in its own go-routine
			go caller(m, n, ch)
		}
	}

	// keep pulling results until we have gotten as many as numAcks
	for ackDone := uint64(0); ackDone < numAcks; ackDone++ {
		r := <-ch
		fmt.Println("ackermann(", r.m, ",", r.n, ")", "=", r.res, "calls:", r.calls)
	}
}
