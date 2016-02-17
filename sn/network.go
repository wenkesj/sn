package sn;

import (
  "sync";
);

var defaultInputCase = float64(0);
var defaultMeasureStart = float64(300);
var defaultVCutoff = float64(30);
var defaultOutputMembranePotentialSuccess = float64(1.0);
var defaultOutputMembranePotentialFail = float64(0.0);

type Network struct {
  neurons []*SpikingNeuron;
};

func NewNetwork(neurons []*SpikingNeuron) *Network {
  return &Network{
    neurons: neurons,
  };
};

func (this *Network) GetNeurons() []*SpikingNeuron {
  return this.neurons;
};

func (this *Network) Simulate(simulation *Simulation) {
  // Create a new shared lock.
  mutex := new(sync.Mutex);

  // Share the simulation across all neurons.
  for index, neuron := range this.neurons {
    // Every neuron except the first one.
    if index != 0 {
      neuron.SetInput(0.0);
    }

    // Conditions for recieveing input.
    neuron.SetInputPredicate(func (t, T1 float64, this *SpikingNeuron) bool {
      return t > T1;
    });

    neuron.SetInputSuccess(func (t, T1 float64, this *SpikingNeuron) float64 {
      return this.GetInput();
    });

    neuron.SetInputFail(func (t, T1 float64, this *SpikingNeuron) float64 {
      return defaultInputCase;
    });

    // Assign predicates.
    neuron.SetPredicate(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      return this.GetV() > defaultVCutoff;
    });

    neuron.SetSuccess(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      // Send the output to the connected neuron.
      mutex.Lock();
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          target := connection.GetTarget();
          // Reduce the product of connection weight and output.
          connectionInput := target.GetInput();
          target.SetInput(connection.GetWeight() * defaultOutputMembranePotentialSuccess + connectionInput);
        }
      }

      // For calculating mean spike rate.
      if timeIndex > defaultMeasureStart {
        this.SetSpikes(this.GetSpikes() + 1);
      }
      mutex.Unlock();
      return true;
    });

    neuron.SetFail(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      // Send the output to the connected neurons.
      // Get the lock as the first neuron with external input.
      mutex.Lock();
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          target := connection.GetTarget();
          // Once the first neuron reaches this lock, send it's response to the next neuron.
          connectionInput := target.GetInput();
          target.SetInput(connection.GetWeight() * defaultOutputMembranePotentialFail + connectionInput);
        }
      }
      // Unlock the next neuron that is connected to this one.
      mutex.Unlock();
      return true;
    });
  }

  // Create a new WaitGroup for simulation to complete after all have been completed.
  var simulationWaitGroup sync.WaitGroup;

  // Create a wait group for neurons.
  var neuronWaitGroup sync.WaitGroup;

  startSimulation := make(chan struct{});

  for _, neuron := range this.neurons {
    simulationWaitGroup.Add(1);
    go neuron.Simulate(simulation, startSimulation, &simulationWaitGroup, &neuronWaitGroup);
  }

  // Wait for the simulation to complete.
  close(startSimulation);
  simulationWaitGroup.Wait();
};
