package main

import (
	"github.com/mesilliac/pulse-simple" // pulse-simple
	"encoding/binary"
	"fmt"
	"math"
  "strconv"
  "os"
)

func main() {
  f, _ := strconv.Atoi(os.Args[1])
	ss := pulse.SampleSpec{pulse.SAMPLE_FLOAT32LE, 44100, 1}
	pb, err := pulse.Playback("pulse-simple test", "playback test", &ss)
	defer pb.Free()
	defer pb.Drain()
	if err != nil {
		fmt.Printf("Could not create playback stream: %s\n", err)
		return
	}
	playfreq(pb, &ss, f)
}

func playfreq(s *pulse.Stream, ss *pulse.SampleSpec, f int) {
  tau := 2 * math.Pi
  data := make([]byte, 8*ss.Rate)
  r := float64(ss.Rate)
  for i := 0; i < 2*int(ss.Rate); i++ {
    sample := float32((math.Sin(tau*float64(f*i)/r) / 3.0)  * float64(i)/(r))
    bits := math.Float32bits(sample)
		binary.LittleEndian.PutUint32(data[4*i:4*i+4], bits)
  }
  s.Write(data)
  for i := 0; i < 2 * int(ss.Rate); i++ {
    //sawtooth wav
    sample := float32(i % (44100/150)) * 0.0001
    //sample := float32(math.Sin(tau*float64(f*i)/r))
    bits := math.Float32bits(sample)
		binary.LittleEndian.PutUint32(data[4*i:4*i+4], bits)
  }
  s.Write(data)
  s.Write(data)
  for i := 0; i < 2 * int(ss.Rate); i++ {
    //sawtooth wav
    sample := float32(i % (44100/180)) * 0.0001
    //sample := float32(math.Sin(tau*float64(f*i)/r))
    bits := math.Float32bits(sample)
		binary.LittleEndian.PutUint32(data[4*i:4*i+4], bits)
  }
  s.Write(data)
}

func time2sample (t int){

}


