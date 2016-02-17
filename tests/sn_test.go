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
var defaultInputCase = float64(0);
var startInput = float64(1.0);
var maxInput = float64(20);
var inputIncrement = float64(0.25);
var lazyIncrement = float64(5);
var inputMeasurements = []float64{1, 5, 10, 15, 20};
var defaultVCutoff = float64(30);
var alphabet = []string{"A","B","C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "O", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"};

// Network parameters
var defaultWeight = float64(100.0);
var defaultOutputMembranePotentialSuccess = float64(1.0);
var defaultOutputMembranePotentialFail = float64(0.0);
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
  output := make([]float64, len(simulation.GetTimeSeries()));

  spikingNeuron := sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD);

  // Conditions for recieveing input.
  spikingNeuron.SetInputPredicate(func (t, T1 float64, this *sn.SpikingNeuron) bool {
    return t > T1;
  });

  spikingNeuron.SetInputSuccess(func (t, T1 float64, this *sn.SpikingNeuron) float64 {
    return this.GetInput();
  });

  spikingNeuron.SetInputFail(func (t, T1 float64, this *sn.SpikingNeuron) float64 {
    return defaultInputCase;
  });

  // Conditions for firing the neuron.
  spikingNeuron.SetPredicate(func (t float64, i int, this *sn.SpikingNeuron) bool {
    return this.GetV() > defaultVCutoff;
  });

  spikingNeuron.SetSuccess(func (t float64, i int, this *sn.SpikingNeuron) bool {
    output[i] = defaultVCutoff;

    if t > defaultMeasureStart {
      this.SetSpikes(this.GetSpikes() + 1);
    }
    return true;
  });

  spikingNeuron.SetFail(func (t float64, i int, this *sn.SpikingNeuron) bool {
    output[i] = this.GetV();
    return true;
  });

  for input := startInput; input < maxInput + inputIncrement; input = input + inputIncrement {
    spikingNeuron.SetInput(input);
    spikingNeuron.Simulate(simulation, nil);

    if IndexOf(inputMeasurements, input) > -1 {
      inputString := strconv.FormatFloat(input, 'f', 6, 64);
      title := "Phasic Spiking Neuron @ I = " + inputString;
      xLabel := "Time Series";
      yLabel := "Membrane Potential";
      legendLabel := "Membrane Potential over Time";
      fileName := "plots/spiking-neuron-" + inputString + ".png";
      GeneratePlot(simulation.GetTimeSeries(), output, title, xLabel, yLabel, legendLabel, fileName);
    }

    spikingNeuron.SetSpikeRate(input, float64(spikingNeuron.GetSpikes()) / defaultMeasureStart);
    spikingNeuron.ResetParameters(defaultA, defaultB, defaultC, defaultD);
  }

  spikeRates := spikingNeuron.GetSpikeRate();
  means := make([]float64, len(inputMeasurements));

  for index, val := range inputMeasurements {
    means[index] = spikeRates[val];
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
  network := make([]*sn.SpikingNeuron, numberOfSpikingNeurons);

  // Create a feed-forward network.
  for i := 0; i < numberOfSpikingNeurons; i++ {
    // Define the neuron.
    network[i] = sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD);

    if i > 0 {
      // Create a connection from one neuron to the next.
      network[i - 1].CreateConnection(network[i], defaultWeight, true, 0);
    }
  }

  // Create a default simulation.
  simulation := sn.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise);

  var testNetwork *sn.Network;

  for input := float64(0); input < maxInput + lazyIncrement; input = input + lazyIncrement {
    if input == 0 {
      input = 1;
    }

    // Feed the first input to the externally connected neuron.
    network[0].SetInput(input);

    // Create a new network simulation of connections feeding from/to neurons.
    testNetwork = sn.NewNetwork(network);
    testNetwork.Simulate(simulation);

    for plotIndex, neuron := range testNetwork.GetNeurons() {
      plotIndexString := alphabet[plotIndex];
      inputString := strconv.FormatFloat(input, 'f', 6, 64);
      title := "Phasic Spiking Neuron " + plotIndexString + " @ I = " + inputString;
      xLabel := "Time Series";
      yLabel := "Membrane Potential";
      legendLabel := "Membrane Potential over Time";
      fileName := "plots/spiking-neuron-" + plotIndexString + "-" + inputString + ".png";
      GeneratePlot(simulation.GetTimeSeries(), neuron.GetOutputs(), title, xLabel, yLabel, legendLabel, fileName);

      neuron.SetSpikeRate(input, float64(neuron.GetSpikes()) / defaultMeasureStart);
      neuron.ResetParameters(defaultA, defaultB, defaultC, defaultD);
    }

    if input == 1 {
      input = 0;
    }
  }

  // Measure and plot out mean spike rates.
  meansMap := make(map[string][]float64);

  for _, neuron := range testNetwork.GetNeurons() {
    means := make([]float64, len(inputMeasurements));
    for i, selectedInput := range inputMeasurements {
      neuronSpikeRate := neuron.GetSpikeRate();
      means[i] = neuronSpikeRate[selectedInput];
    }
    meansMap[neuron.GetId()] = means;
  }

  plotIndex := int(0);
  for _, value := range meansMap {
    plotIndexString := alphabet[plotIndex];
    title := "Mean Spike Rate of Phasic Spiking Neuron " + plotIndexString;
    xLabel := "Input";
    yLabel := "Spike Rate";
    legendLabel := "Spike Rate as a Function of Input";
    fileName := "plots/spiking-neuron-network-mean-spike-rate-" + plotIndexString + ".png";
    GeneratePlot(inputMeasurements, value, title, xLabel, yLabel, legendLabel, fileName);
    plotIndex++;
  }
};
