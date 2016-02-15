package tests;

import (
  "testing";
  "strconv";
  "github.com/gonum/plot";
  "github.com/gonum/plot/plotter";
  "github.com/gonum/plot/plotutil";
  "github.com/gonum/plot/vg";
  "github.com/wenkesj/sn/sn";
);

// Simulation parameters.
var defaultSteps = float64(500);
var defaultTau = float64(0.25);
var defaultStart = float64(0);
var defaultStepRise = float64(50);
var defaultMeasureStart = float64(300);

// Neuron parameters.
var defaultA = float64(0.02);
var defaultB = float64(0.25);
var defaultC = float64(-65);
var defaultD = float64(6);

// Test parameters
var startInput = float64(1.0);
var maxInput = float64(20);
var inputIncrement = float64(0.25);
var inputMeasurements = []float64{1, 5, 10, 15, 20};

func GeneratePlot(x, y []float64, title, xLabel, yLabel, legendLabel, fileName string) {
  outPlotPoints := make(plotter.XYs, len(x));

  for i, _ := range x {
    outPlotPoints[i].X = x[i];
    outPlotPoints[i].Y = y[i];
  }

  outPlot, err := plot.New();
  if err != nil {
    panic(err);
  }

  outPlot.Title.Text = title;
  outPlot.X.Label.Text = xLabel;
  outPlot.Y.Label.Text = yLabel;

  err = plotutil.AddLines(outPlot,
    legendLabel, outPlotPoints);
  if err != nil {
    panic(err);
  }

  if err := outPlot.Save(6 * vg.Inch, 6 * vg.Inch, fileName); err != nil {
    panic(err);
  }
};

func IndexOf(list []float64, targetValue float64) int {
  for index, value := range list {
    if value == targetValue {
      return index;
    }
  }
  return -1;
};

func TestSpikingNeuronSimulation(t *testing.T) {
  simulation := sn.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise, defaultMeasureStart);
  spikingNeuron := sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD);

  for input := startInput; input < maxInput + inputIncrement; input = input + inputIncrement {
    VV := spikingNeuron.Simulate(input, simulation);

    if IndexOf(inputMeasurements, input) > -1 {
      inputString := strconv.FormatFloat(input, 'f', 6, 64);
      title := "Phasic Spiking Neuron @ I = " + inputString;
      xLabel := "Time Series";
      yLabel := "Membrane Potential";
      legendLabel := "Membrane Potential over Time";
      fileName := "spiking-neuron-" + inputString + ".png";
      GeneratePlot(simulation.GetTimeSeries(), VV, title, xLabel, yLabel, legendLabel, fileName);
    }

    simulation.SetSpikeRate(input, spikingNeuron.GetSpikeRate());
    spikingNeuron.Reset();
  }

  meanSketch := simulation.GetSketch();
  means := make([]float64, len(inputMeasurements));

  for index, val := range inputMeasurements {
    means[index] = meanSketch[val];
  }

  title := "Mean Spike Rate of Phasic Spiking Neuron";
  xLabel := "Input";
  yLabel := "Spike Rate";
  legendLabel := "Spike Rate as a Function of Input";
  fileName := "spiking-neuron-mean-spike-rate.png";
  GeneratePlot(inputMeasurements, means, title, xLabel, yLabel, legendLabel, fileName);
};
