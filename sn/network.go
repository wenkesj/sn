package sn;

import (
  "time";
  // "fmt";
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
  // Share the simulation across all neurons.
  for index, neuron := range this.neurons {
    // Every neuron except the first one.
    if index != 0 {
      neuron.SetInput(0.0);
    }

    // Conditions for recieveing input.
    neuron.SetInputPredicate(func (i int, t, T1 float64, this *SpikingNeuron) bool {
      return t > T1;
    });

    neuron.SetInputSuccess(func (i int, t, T1 float64, this *SpikingNeuron) float64 {
      inputSum := this.GetInput();
      for _, connection := range this.GetConnections() {
        if !connection.IsWriteable() {
          // Sum the connections to the neuron.
          inputSum += connection.GetOutput();
          // fmt.Println(this.GetId(),"connection recieved with value of", inputSum);
        }
      }
      return inputSum;
    });

    neuron.SetInputFail(func (i int, t, T1 float64, this *SpikingNeuron) float64 {
      return defaultInputCase;
    });

    // Assign predicates.
    neuron.SetPredicate(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      return this.GetV() > defaultVCutoff;
    });

    neuron.SetSuccess(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      // Set it's own output.
      this.SetOutput(currentIndex, defaultVCutoff);

      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(defaultOutputMembranePotentialSuccess);
          // fmt.Println("Fire: ",this.GetId(),", connection sent",connection.GetOutput());
        }
      }

      // For calculating mean spike rate.
      if timeIndex > defaultMeasureStart {
        this.SetSpikes(this.GetSpikes() + 1);
      }
      return true;
    });

    neuron.SetFail(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      // Set its own output unit
      this.SetOutput(currentIndex, this.GetV());

      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(defaultOutputMembranePotentialFail);
          // fmt.Println("Fail: ",this.GetId(),", connection sent",connection.GetOutput());
        }
      }
      return true;
    });
  }

  // Create an atomic neuron.
  var simulationWaitGroup sync.WaitGroup;
  var innerWaitGroup sync.WaitGroup;
  mutexSignal := new(sync.Mutex);
  atomicNeuron := NewAtomicNeuron(&simulationWaitGroup, &innerWaitGroup, mutexSignal, len(this.neurons));

  // The first neuron starts the simulation ahead of all the others.
  // It grabs the outer lock, blocking all other neurons.
  // It then calculates its input based on the external connection inputting that are ready.
  // It then increments an atomic counter.
  // It unlocks the outer lock so the next neuron can go through.
  // It then calculates its output and atomically sets the output connections values and sets them to be read.
  // The neuron then unlocks the next neuron to go through and waits for it to finish.
  // After all the neurons finish, the first neuron goes again first.
  // This then repeats over the time series...
  for _, neuron := range this.neurons {
    atomicNeuron.AddWaitGroup(1);
    go neuron.Simulate(simulation, atomicNeuron);
    time.Sleep(time.Millisecond);
  }

  // Wait for the simulation to complete.
  atomicNeuron.FinishWaitGroup();
  // fmt.Println("Finally Finished...");
};
