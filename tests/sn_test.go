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
var defaultVCutoff = float64(30);

// Network parameters
var defaultWeight = float64(100.0);
var defaultOutputMembranePotential = float64(1.0);
var numberOfSpikingNeurons = 2;

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

func TestSingleSpikingNeuronSimulation(t *testing.T) {
  simulation := sn.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise);
  spikingNeuron := sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD);

  output := make([]float64, len(simulation.GetTimeSeries()));

  for input := startInput; input < maxInput + inputIncrement; input = input + inputIncrement {

    // Conditions for firing the neuron.
    predicate := func (t float64, i int, this *sn.SpikingNeuron) bool {
      return this.GetV() > defaultVCutoff;
    };

    success := func (t float64, i int, this *sn.SpikingNeuron) bool {
      output[i] = defaultVCutoff;
      this.SetV(this.GetC());
      this.SetU(this.GetU() + this.GetD());
      if t > defaultMeasureStart {
        this.SetSpikes(this.GetSpikes() + 1);
      }
      return true;
    };

    fail := func (t float64, i int, this *sn.SpikingNeuron) bool {
      output[i] = this.GetV();
      return true;
    };

    spikingNeuron.Simulate(input, simulation, predicate, success, fail);

    if IndexOf(inputMeasurements, input) > -1 {
      inputString := strconv.FormatFloat(input, 'f', 6, 64);
      title := "Phasic Spiking Neuron @ I = " + inputString;
      xLabel := "Time Series";
      yLabel := "Membrane Potential";
      legendLabel := "Membrane Potential over Time";
      fileName := "plots/spiking-neuron-" + inputString + ".png";
      GeneratePlot(simulation.GetTimeSeries(), output, title, xLabel, yLabel, legendLabel, fileName);
    }

    simulation.SetSpikeRate(input, float64(spikingNeuron.GetSpikes()) / defaultMeasureStart);
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
  fileName := "plots/spiking-neuron-mean-spike-rate.png";
  GeneratePlot(inputMeasurements, means, title, xLabel, yLabel, legendLabel, fileName);
};

// Neuron A and B example.
func TestSpikingNeuronNetwork(t *testing.T) {
  // Create a group of neurons with the default parameters.
  networkNeuronGroup := make([]*sn.NetworkNeuron, numberOfSpikingNeurons);

  for i := 0; i < numberOfSpikingNeurons; i++ {
    spikingNeuron := sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD);
    networkNeuronGroup[i] = sn.NewNetworkNeuron(spikingNeuron);
  }

  neuronA := networkNeuronGroup[0];
  neuronB := networkNeuronGroup[1];

  // Create a one way connection from A to B with a weight.
  neuronA.CreateConnection(neuronB, defaultWeight, true, 0);

  // Create a default simulation.
  simulation := sn.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise);

  // Create a new network simulation of connections feeding from/to neurons.
  networkSimulation := sn.NewNetworkSimulation(simulation, networkNeuronGroup);
  networkSimulation.Simulate();
};
