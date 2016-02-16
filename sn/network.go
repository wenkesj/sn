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
      // Set locks and wait for signals.
      return this.GetInput();
      // Wait for signals to get correct input.
    });

    neuron.SetInputFail(func (t, T1 float64, this *SpikingNeuron) float64 {
      // Set locks and wait for signals.
      return defaultInputCase;
      // Wait for signals to get correct input.
    });

    // Assign predicates.
    neuron.SetPredicate(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      return this.GetV() > defaultVCutoff;
    });

    neuron.SetSuccess(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.GetTarget().SetInput(connection.GetWeight() * defaultOutputMembranePotentialSuccess);
        }
      }

      // For calculating mean spike rate.
      if timeIndex > defaultMeasureStart {
        this.SetSpikes(this.GetSpikes() + 1);
      }
      return true;
    });

    neuron.SetFail(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      // Send the output to the connected neurons.
      // Get the lock as the first neuron with external input.
      mutex.Lock();
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          // Once the first neuron reaches this lock, send it's response to the next neuron.
          connection.GetTarget().SetInput(connection.GetWeight() * defaultOutputMembranePotentialFail);
        }
      }
      // Unlock the next neuron that is connected to this one.
      mutex.Unlock();
      return true;
    });
  }

  // Create a new WaitGroup for simulation to complete after all have been completed.
  var waitGroup sync.WaitGroup;
  startSimulation := make(chan struct{});

  for _, neuron := range this.neurons {
    waitGroup.Add(1);
    go neuron.Simulate(simulation, startSimulation, &waitGroup);
  }

  // Wait for the simulation to complete.
  close(startSimulation);
  waitGroup.Wait();
};
