package main

import (
	"github.com/mesilliac/pulse-simple" // pulse-simple
	"encoding/binary"
	"fmt"
	"math"
  "strconv"
  "os"
  "gonum.org/v1/plot"
  "gonum.org/v1/plot/plotter"
)

const sample_rate uint32 = 44100;
const tau float64 = 2 * math.Pi;

func main() {
  f_in, _ := strconv.Atoi(os.Args[1]);
  f := uint32(f_in);
  ss := pulse.SampleSpec{pulse.SAMPLE_S16LE, sample_rate, 2};
	pb, err := pulse.Playback("pulse-simple test", "playback test", &ss);
	defer pb.Free();
	defer pb.Drain();
	if err != nil {
		fmt.Printf("Could not create playback stream: %s\n", err);
		return;
	}
	playfreq(pb, &ss, f);
}

func playfreq(s *pulse.Stream, ss *pulse.SampleSpec, f uint32) {
  data := make([]byte, 2*44100);
  plot := make([]uint16, 44100);
  var i uint32;
  for i = 0; i < 44100; i++ {
    //bits := sinewave(f, i);
    //bits := sawtooth(f, i);
    bits := triangle(f, i);
    plot[i] = bits;
		binary.LittleEndian.PutUint16(data[2*i:2*i+2], bits);
  }
  s.Write(data);
  int2plot(plot);
}

func sinewave (freq uint32, phase uint32) uint16{
  sample := math.Sin(tau*float64(freq*phase)/float64(sample_rate));
  return uint16(sample * 30000);
}

func sawtooth (freq uint32, phase uint32) uint16{
  var sample float32 = 0;
  spo := sample_rate/freq;
  sample = (float32(phase % spo)/float32(spo)) * 2 -1;
  return uint16(sample * 30000);
}

func triangle (freq uint32, phase uint32) uint16{
  var sample float32 = 0;
  spo := sample_rate/freq;
  spo2 := float32(spo / 2);
  spo4 := float32(spo / 4);
  pspo := float32(phase % spo);
  if pspo <= spo2 {
    sample = (pspo / spo4) - 1;
  } else {
    sample = 3 - (pspo / spo4);
  }
  return uint16(sample * 30000);
}

func int2XYs(intbuff []uint16) (plotter.XYs) {
  var xys plotter.XYs;
  for i:=0;i<1000;i++{
    xys = append(xys, struct{ X, Y float64 }{float64(i), float64(intbuff[i])});
  }
  return xys;
}

func int2plot(intbuff []uint16) {
  path := "plot.png";
  f, _ := os.Create(path);
  p, _ := plot.New();
  xys := int2XYs(intbuff);
  line, _ := plotter.NewLine(xys);
  line.StepStyle = plotter.PostStep;
  p.Add(line);
  wt, _ := p.WriterTo(2000,1500,"png");
  wt.WriteTo(f);
  f.Close();
}
