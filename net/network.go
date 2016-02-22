package net;

import (
  "time";
  // "fmt";
  "sync";
  "github.com/wenkesj/sn/sn";
  "github.com/wenkesj/sn/sim";
  "github.com/wenkesj/sn/vars";
  "github.com/wenkesj/sn/group";
);

type Network struct {
  neurons []*sn.SpikingNeuron;
};

func NewNetwork(neurons []*sn.SpikingNeuron) *Network {
  return &Network{
    neurons: neurons,
  };
};

func (this *Network) GetNeurons() []*sn.SpikingNeuron {
  return this.neurons;
};

func (_this *Network) Simulate(simulation *sim.Simulation) {
  // Share the simulation across all neurons.
  for index, neuron := range _this.neurons {
    // Every neuron except the first one.
    if index != 0 {
      neuron.SetInput(0.0);
    }

    // Conditions for recieveing input.
    neuron.SetInputPredicate(func (i int, t, T1 float64, this *sn.SpikingNeuron) bool {
      return t > T1;
    });

    neuron.SetInputSuccess(func (i int, t, T1 float64, this *sn.SpikingNeuron) float64 {
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

    neuron.SetInputFail(func (i int, t, T1 float64, this *sn.SpikingNeuron) float64 {
      return vars.GetDefaultInputCase();
    });

    // Assign predicates.
    neuron.SetPredicate(func (timeIndex float64, currentIndex int, this *sn.SpikingNeuron) bool {
      return this.GetV() > vars.GetDefaultVCutoff();
    });

    neuron.SetSuccess(func (timeIndex float64, currentIndex int, this *sn.SpikingNeuron) bool {
      // Set it's own output.
      this.SetOutput(currentIndex, vars.GetDefaultVCutoff());

      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(vars.GetDefaultOutputMembranePotentialSuccess());
          // fmt.Println("Fire: ",this.GetId(),", connection sent",connection.GetOutput());
        }
      }

      // For calculating mean spike rate.
      if timeIndex > vars.GetDefaultMeasureStart() {
        this.SetSpikes(this.GetSpikes() + 1);
        this.SetTimeSpike(timeIndex);
      }
      return true;
    });

    neuron.SetFail(func (timeIndex float64, currentIndex int, this *sn.SpikingNeuron) bool {
      // Set its own output unit
      this.SetOutput(currentIndex, this.GetV());

      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(vars.GetDefaultOutputMembranePotentialFail());
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
  neuronManager := group.NewNeuronManager(&simulationWaitGroup, &innerWaitGroup, mutexSignal, len(_this.neurons));

  // The first neuron starts the simulation ahead of all the others.
  // It grabs the outer lock, blocking all other neurons.
  // It then calculates its input based on the external connection inputting that are ready.
  // It then increments an atomic counter.
  // It unlocks the outer lock so the next neuron can go through.
  // It then calculates its output and atomically sets the output connections values and sets them to be read.
  // The neuron then unlocks the next neuron to go through and waits for it to finish.
  // After all the neurons finish, the first neuron goes again first.
  // This then repeats over the time series...
  for _, neuron := range _this.neurons {
    neuronManager.AddWaitGroup(1);
    go neuron.Simulate(simulation, neuronManager);
    time.Sleep(time.Millisecond);
  }

  // Wait for the simulation to complete.
  neuronManager.FinishWaitGroup();
  // fmt.Println("Finally Finished...");
};
