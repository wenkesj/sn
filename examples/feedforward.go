package main;

import (
  "strconv";
  "github.com/wenkesj/sn/plots";
  "github.com/wenkesj/sn/sn";
  "github.com/wenkesj/sn/sim";
  "github.com/wenkesj/sn/net";
  "github.com/garyburd/redigo/redis";
);

func main() {
  // Create a new simulation with some parameters.
  var alphabet = []string{"A","B","C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "O", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"};

  numberOfSteps := float64(500);
  interval := float64(0.25);
  start := float64(0);
  turnOn := float64(50);
  simulation := sim.NewSimulation(numberOfSteps, interval, start, turnOn);

  // Connect to the redis instance.
  redisConnection, err := redis.Dial("tcp", ":6379");
  if err != nil {
    panic(err);
  }
  defer redisConnection.Close();

  // Allocate the neurons.
  numberOfNeurons := 6;
  a := float64(0.02);
  b := float64(0.25);
  c := float64(-65);
  d := float64(6);
  neurons := make([]*sn.SpikingNeuron, numberOfNeurons);
  for i := 0; i < numberOfNeurons; i++ {
    neurons[i] = sn.NewSpikingNeuron(a, b, c, d, int64(i), redisConnection);
  }

  // Create external sources to provide input to the neurons.
  numberOfExternalSources := 4;
  externalSources := make([]*sn.SpikingNeuron, numberOfExternalSources);
  for i := 0; i < numberOfExternalSources; i++ {
    externalSources[i] = sn.NewSpikingNeuron(0, 0, 0, 0, int64(-i), redisConnection);
  }

  // Manually connect the external sources to the first layer neurons.
  // Create one way connections.
  externalSources[0].CreateConnection(neurons[0], 1.0, true, 0);
  externalSources[1].CreateConnection(neurons[0], 1.0, true, 0);
  externalSources[1].CreateConnection(neurons[1], 1.0, true, 0);
  externalSources[2].CreateConnection(neurons[1], 1.0, true, 0);
  externalSources[2].CreateConnection(neurons[2], 1.0, true, 0);
  externalSources[3].CreateConnection(neurons[2], 1.0, true, 0);

  // Manually connect the neurons to feed forward.
  // Layer 1...
  neurons[0].CreateConnection(neurons[3], 100.0, true, 0);
  neurons[1].CreateConnection(neurons[3], 100.0, true, 0);
  neurons[1].CreateConnection(neurons[4], 100.0, true, 0);
  neurons[2].CreateConnection(neurons[4], 100.0, true, 0);

  // Layer 2...
  neurons[3].CreateConnection(neurons[5], 100.0, true, 0);
  neurons[4].CreateConnection(neurons[5], 100.0, true, 0);

  // Create a some random data for the connections.
  externalInputs := [][]float64{
    {15.0, 15.0},
    {15.0, 15.0},
    {15.0, 15.0},
    {15.0, 15.0},
  };

  // Apply the inputs of each external source.
  for i, externalSource := range externalSources {
    for j, connection := range externalSource.GetConnections() {
      if connection.IsWriteable() {
        connection.SetOutput(externalInputs[i][j]);
      }
    }
  }

  // Create the network simulation and run it.
  network := net.NewNetwork(neurons);
  network.Simulate(simulation);

  for plotIndex, neuron := range network.GetNeurons() {
    plotIndexString := alphabet[plotIndex];
    numberOfConnections := strconv.Itoa(len(neuron.GetConnections()));
    title := "Phasic Spiking Neuron " + plotIndexString + " with connections " + numberOfConnections;
    xLabel := "Time Series";
    yLabel := "Membrane Potential";
    legendLabel := "Membrane Potential over Time";
    fileName := "plots/spiking-neuron-" + plotIndexString + "-with-" + numberOfConnections + "-connections.png";
    plots.GeneratePlot(simulation.GetTimeSeries(), neuron.GetOutputs(), title, xLabel, yLabel, legendLabel, fileName);
  }
};
