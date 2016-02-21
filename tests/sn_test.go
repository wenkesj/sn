package tests;

import (
  "testing";
  "strconv";
  // Plots
  "github.com/wenkesj/sn/plots";
  // SN
  "github.com/wenkesj/sn/sn";
  "github.com/wenkesj/sn/sim";
  "github.com/wenkesj/sn/net";
  // Deps
  "github.com/garyburd/redigo/redis";
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
var numberOfSpikingNeurons = 3;
var numberOfExternalConnections = 3;

func IndexOf(list []float64, targetValue float64) int {
  for index, value := range list {
    if value == targetValue {
      return index;
    }
  }
  return -1;
};

func TestSingleSpikingNeuronSimulation(t *testing.T) {
  simulation := sim.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise);
  output := make([]float64, len(simulation.GetTimeSeries()));

  spikingNeuron := sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD, 0, nil);

  // Conditions for recieveing input.
  spikingNeuron.SetInputPredicate(func (i int, t, T1 float64, this *sn.SpikingNeuron) bool {
    return t > T1;
  });

  spikingNeuron.SetInputSuccess(func (i int, t, T1 float64, this *sn.SpikingNeuron) float64 {
    return this.GetInput();
  });

  spikingNeuron.SetInputFail(func (i int, t, T1 float64, this *sn.SpikingNeuron) float64 {
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
      plots.GeneratePlot(simulation.GetTimeSeries(), output, title, xLabel, yLabel, legendLabel, fileName);
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
  plots.GeneratePlot(inputMeasurements, means, title, xLabel, yLabel, legendLabel, fileName);
};

// External Source -> Neuron A -> Neuron B -> Neuron C ... example.
func TestSingleConnectionSpikingNeuronNetwork(t *testing.T) {
  // Create a group of neurons with the default parameters.
  // Connect to a redis client.
  redisConnection, err := redis.Dial("tcp", ":6379");
  if err != nil {
    // Handle error if any...
    panic(err);
  }
  defer redisConnection.Close();

  network := make([]*sn.SpikingNeuron, numberOfSpikingNeurons);

  // Create a feed-forward network.
  for i := 0; i < numberOfSpikingNeurons; i++ {
    // Define the neuron.
    network[i] = sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD, int64(i), redisConnection);

    if i > 0 {
      // Create a connection from one neuron to the next.
      network[i - 1].CreateConnection(network[i], defaultWeight, true, 0);
    }
  }

  // This is a hacky way to create an external connection...
  externalSource := sn.NewSpikingNeuron(0, 0, 0, 0, int64(-1), redisConnection);
  externalSource.CreateConnection(network[0], 1.0, true, 0);

  // Create a default simulation.
  simulation := sim.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise);

  var testNetwork *net.Network;

  for input := float64(0); input < maxInput + lazyIncrement; input = input + lazyIncrement {
    if input == 0 {
      input = 1;
    }

    for _, connection := range externalSource.GetConnections() {
      if connection.IsWriteable() {
        connection.SetOutput(input);
      }
    }

    // Create a new network simulation of connections feeding from/to neurons.
    testNetwork = net.NewNetwork(network);
    testNetwork.Simulate(simulation);

    for plotIndex, neuron := range testNetwork.GetNeurons() {
      plotIndexString := alphabet[plotIndex];
      inputString := strconv.FormatFloat(input, 'f', 6, 64);
      title := "Phasic Spiking Neuron " + plotIndexString + " @ I = " + inputString;
      xLabel := "Time Series";
      yLabel := "Membrane Potential";
      legendLabel := "Membrane Potential over Time";
      fileName := "plots/spiking-neuron-" + plotIndexString + "-" + inputString + ".png";
      plots.GeneratePlot(simulation.GetTimeSeries(), neuron.GetOutputs(), title, xLabel, yLabel, legendLabel, fileName);

      neuron.SetSpikeRate(input, float64(neuron.GetSpikes()) / defaultMeasureStart);
      neuron.ResetParameters(defaultA, defaultB, defaultC, defaultD);
    }

    if input == 1 {
      input = 0;
    }
  }

  // Measure and plot out mean spike rates.
  meansMap := make(map[int64][]float64);

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
    plots.GeneratePlot(inputMeasurements, value, title, xLabel, yLabel, legendLabel, fileName);
    plotIndex++;
  }
};

func TestMultipleConnectionSpikingNeuronNetwork(t *testing.T) {
  // Create a group of neurons with the default parameters.
  // Connect to a redis client.
  redisConnection, err := redis.Dial("tcp", ":6379");
  if err != nil {
    // Handle error if any...
    panic(err);
  }
  defer redisConnection.Close();

  network := make([]*sn.SpikingNeuron, numberOfSpikingNeurons);

  // Create a feed-forward network.
  for i := 0; i < numberOfSpikingNeurons; i++ {
    // Define the neuron.
    network[i] = sn.NewSpikingNeuron(defaultA, defaultB, defaultC, defaultD, int64(i), redisConnection);

    if i > 0 {
      // Create a connection from one neuron to the next.
      network[i - 1].CreateConnection(network[i], defaultWeight, true, 0);
    }
  }

  // This is a hacky way to create an external connection...
  externalConnections := make([]*sn.SpikingNeuron, numberOfExternalConnections);
  for index, _ := range externalConnections {
    externalConnections[index] = sn.NewSpikingNeuron(0, 0, 0, 0, int64(-index), redisConnection);
    externalConnections[index].CreateConnection(network[0], 1.0, true, 0);
  }

  // Create a default simulation.
  simulation := sim.NewSimulation(defaultSteps, defaultTau, defaultStart, defaultStepRise);

  var testNetwork *net.Network;

  for input := float64(0); input < maxInput + lazyIncrement; input = input + lazyIncrement {
    if input == 0 {
      input = 1;
    }

    for _, externalSource := range externalConnections {
      for _, connection := range externalSource.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(input);
        }
      }
    }

    // Create a new network simulation of connections feeding from/to neurons.
    testNetwork = net.NewNetwork(network);
    testNetwork.Simulate(simulation);

    for plotIndex, neuron := range testNetwork.GetNeurons() {
      plotIndexString := alphabet[plotIndex];
      numberOfConnections := strconv.Itoa(len(neuron.GetConnections()));
      inputString := strconv.FormatFloat(input, 'f', 6, 64);
      title := "Phasic Spiking Neuron " + plotIndexString + " with connections " + numberOfConnections + " @ I = " + inputString;
      xLabel := "Time Series";
      yLabel := "Membrane Potential";
      legendLabel := "Membrane Potential over Time";
      fileName := "plots/spiking-neuron-" + plotIndexString + "-" + inputString + "-with-" + numberOfConnections + "-connections.png";
      plots.GeneratePlot(simulation.GetTimeSeries(), neuron.GetOutputs(), title, xLabel, yLabel, legendLabel, fileName);

      neuron.SetSpikeRate(input, float64(neuron.GetSpikes()) / defaultMeasureStart);
      neuron.ResetParameters(defaultA, defaultB, defaultC, defaultD);
    }

    if input == 1 {
      input = 0;
    }
  }
};
