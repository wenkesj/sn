package sn;

type Network struct {
  neurons []*SpikingNeuron;
};

func NewNetwork(neurons []*SpikingNeuron) *Network {
  return &Network{
    neurons: neurons,
  };
};

func (this *Network) Simulate(input float64, simulation *Simulation) {

  // Share the simulation across all neurons.
  for _, neuron := range this.neurons {
    go neuron.Simulate(input, simulation);
  }
};
