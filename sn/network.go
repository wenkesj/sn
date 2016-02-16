package sn;

import (
  "sync"
);

type Network struct {
  neurons []*SpikingNeuron;
};

func NewNetwork(neurons []*SpikingNeuron) *Network {
  return &Network{
    neurons: neurons,
  };
};

func (this *Network) Simulate(input float64, simulation *Simulation) {
  var awaitGroup sync.WaitGroup;
  awaitGroup.Add(len(this.neurons));

  // Share the simulation across all neurons.
  for _, neuron := range this.neurons {
    go neuron.Simulate(input, simulation);
    defer awaitGroup.Done();
  }
  awaitGroup.Wait();
};
