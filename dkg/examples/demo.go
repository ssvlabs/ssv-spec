package main

import "C"
import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/herumi/bls-eth-go-binary/bls"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type M = base.Message

func init() {
	log.SetFormatter(&log.TextFormatter{})

	//Set log output to standard output (default output is stderr, standard error)
	//Log message output can be any io.writer type
	log.SetOutput(os.Stdout)

	log.SetLevel(log.TraceLevel)
}

func main() {

	t := 2
	n := 4
	var (
		ins       []chan M
		outs      []chan M
		kMachines []*keygen.Runner
	)
	defer func() {
		for _, ch := range ins {
			close(ch)
		}
		for _, ch := range outs {
			close(ch)
		}
	}()

	id := dkg.RequestID{}
	for i, _ := range id {
		id[i] = 1
	}
	for i := 1; i < n+1; i++ {
		in := make(chan M, n)
		out := make(chan M, n)
		keygen, _ := keygen.NewRunner(id, uint32(i), uint32(t), uint32(n), in, out)
		ins = append(ins, in)
		outs = append(outs, out)
		kMachines = append(kMachines, keygen)
	}

	go func(o1 <-chan M, o2 <-chan M, o3 <-chan M, o4 <-chan M, i1 chan<- M, i2 chan<- M, i3 chan<- M, i4 chan<- M) {
		send := func(msg base.Message) {
			if msg.Header.Receiver == 0 {
				i1 <- msg
				i2 <- msg
				i3 <- msg
				i4 <- msg
			} else if msg.Header.Receiver == 1 {
				i1 <- msg
			} else if msg.Header.Receiver == 2 {
				i2 <- msg
			} else if msg.Header.Receiver == 3 {
				i3 <- msg
			} else if msg.Header.Receiver == 4 {
				i4 <- msg
			}
		}
		for {
			select {
			case m, ok := <-o1:
				if ok {
					send(m)
				}
			case m, ok := <-o2:
				if ok {
					send(m)
				}
			case m, ok := <-o3:
				if ok {
					send(m)
				}
			case m, ok := <-o4:
				if ok {
					send(m)
				}
			case <-time.After(1 * time.Second):
			}
		}
	}(outs[0], outs[1], outs[2], outs[3], ins[0], ins[1], ins[2], ins[3])

	log.Debug("Starting keygen")
	go kMachines[0].ProcessLoop()
	go kMachines[1].ProcessLoop()
	go kMachines[2].ProcessLoop()
	go kMachines[3].ProcessLoop()

	kMachines[0].Initialize()
	kMachines[1].Initialize()
	kMachines[2].Initialize()
	kMachines[3].Initialize()
	log.Debug("KeygenSimple started")

	var allFinished bool
	for !allFinished {
		select {
		case <-time.After(1 * time.Second):
			allFinished = true
			for _, machine := range kMachines {
				allFinished = allFinished && machine.Keygen.Output != nil
			}
			if allFinished {
				break
			}
		}
	}
	log.Debug("KeygenSimple completed")

	msgHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 32 * "a"
	log.Infof("PublicKey is: %v", hex.EncodeToString(kMachines[0].Keygen.Output.PublicKey))
	log.Infof("Msg is: %v", hex.EncodeToString([]byte(msgHash)))

	pSigs := make([]bls.G2, n)
	for i, machine := range kMachines {
		sk := new(bls.SecretKey)
		sk.Deserialize(machine.Keygen.Output.SecretShare)
		pSigs[i].Deserialize(sk.SignByte([]byte(msgHash)).Serialize())
	}

	ids := make([]bls.Fr, n)
	for i := 0; i < n; i++ {
		ids[i].SetInt64(int64(i + 1))
	}

	sig := bls.G2{}

	bls.G2LagrangeInterpolation(&sig, ids[1:], pSigs[1:])
	log.Infof("Signature is: %v\n", hex.EncodeToString(sig.Serialize()))
	bls.G2LagrangeInterpolation(&sig, ids[:3], pSigs[:3])
	log.Infof("Signature is: %v\n", hex.EncodeToString(sig.Serialize()))
}
