package sn;

type NetworkSimulation struct {
  simulation *Simulation;
  neurons []*NetworkNeuron;
  masterNeuron *NetworkNeuron;
  writeableChannels []chan float64;
};

func NewNetworkSimulation(simulation *Simulation, neurons []*NetworkNeuron) *NetworkSimulation {
  // Create N channels for neuron simulation.
  masterNeuron := &NetworkNeuron{};
  numberOfWritableConnections := 0;

  // Get all writable channels.
  for _, neuron := range neurons {
    for _, connection := range neuron.GetConnections() {
      if connection.IsWriteable() {
        if numberOfWritableConnections == 0 {
          masterNeuron = neuron;
        }
        numberOfWritableConnections++;
      }
    }
  }

  writeableChannels := make([]chan float64, numberOfWritableConnections);

  // Share the output of one neuron to the output of another.
  for i := 0; i < numberOfWritableConnections; i++ {
    writeableChannels[i] = make(chan float64);
  }

  return &NetworkSimulation{
    simulation: simulation,
    neurons: neurons,
    masterNeuron: masterNeuron,
    writeableChannels: writeableChannels,
  };
};

func (this *NetworkSimulation) Simulate() {

};
