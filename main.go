package main

import (
	"github.com/mesilliac/pulse-simple" // pulse-simple
	"encoding/binary"
	"fmt"
	"math"
  "strconv"
  "os"
)

const sample_rate uint32 = 44100;
const tau float64 = 2 * math.Pi;

func main() {
  f_in, _ := strconv.Atoi(os.Args[1])
  f := uint32(f_in)
  ss := pulse.SampleSpec{pulse.SAMPLE_S16LE, sample_rate, 2}
	pb, err := pulse.Playback("pulse-simple test", "playback test", &ss)
	defer pb.Free()
	defer pb.Drain()
	if err != nil {
		fmt.Printf("Could not create playback stream: %s\n", err)
		return
	}
	playfreq(pb, &ss, f)
}

func playfreq(s *pulse.Stream, ss *pulse.SampleSpec, f uint32) {
  data := make([]byte, 2*44100)
  var i uint32;
  for i = 0; i < 44100; i++ {
    bits := sinewave(f, i)
    //bits := sawtooth(f, i)
		binary.LittleEndian.PutUint16(data[2*i:2*i+2], bits)
  }
  s.Write(data)
}

func sinewave (freq uint32, phase uint32) uint16{
  sample := math.Sin(tau*float64(freq*phase)/float64(sample_rate))
  return uint16(sample * 32000);
}

func sawtooth (freq uint32, phase uint32) uint16{
  //TODO fix this
  sample := ((phase % (sample_rate/freq))/(sample_rate/(2*freq)) - 1);
  return uint16(sample * 32000);
}
