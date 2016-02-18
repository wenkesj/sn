package sn;

import (
  "time";
  "fmt"
  "sync";
);

var defaultInputCase = float64(0);
var defaultMeasureStart = float64(300);
var defaultVCutoff = float64(30);
var defaultOutputMembranePotentialSuccess = float64(1.0);
var defaultOutputMembranePotentialFail = float64(0.0);

type StartChannel chan struct{};

type AtomicNeuron struct {
  startSimulation StartChannel;
  simulationWaitGroup *sync.WaitGroup;
  mutexSignal *sync.Mutex;
  innerSignal *sync.WaitGroup;
  atomicSignalCondition *uint64;
  numberOfOtherNeurons int;
};

func NewAtomicNeuron(
  startSimulation StartChannel,
  simulationWaitGroup *sync.WaitGroup,
  mutexSignal *sync.Mutex,
  innerSignal *sync.WaitGroup,
  atomicSignalCondition *uint64,
  numberOfOtherNeurons int,
) *AtomicNeuron {

  return &AtomicNeuron{
    startSimulation: startSimulation,
    simulationWaitGroup: simulationWaitGroup,
    mutexSignal: mutexSignal,
    innerSignal: innerSignal,
    atomicSignalCondition: atomicSignalCondition,
    numberOfOtherNeurons: numberOfOtherNeurons,
  };
};

func (this *AtomicNeuron) GetStartChannel() StartChannel {
  return this.startSimulation;
}

func (this *AtomicNeuron) GetNumber() int {
  return this.numberOfOtherNeurons;
};

func (this *AtomicNeuron) AddWaitGroup(num int) {
  this.simulationWaitGroup.Add(num);
};

func (this *AtomicNeuron) OuterLock() {
  this.mutexSignal.Lock();
};

func (this *AtomicNeuron) OuterUnlock() {
  this.mutexSignal.Unlock();
};

func (this *AtomicNeuron) Wait(neuron *SpikingNeuron) {
  // Next, this neuron lets the next neuron go through.
  this.innerSignal.Done();
  time.Sleep(time.Millisecond);
  this.OuterUnlock();
  fmt.Println(neuron.GetId(),"waiting for other neurons...");
  this.innerSignal.Wait();
  fmt.Println(neuron.GetId(),"is no longer waiting!");
};

func (this *AtomicNeuron) FinishWaitGroup() {
  this.simulationWaitGroup.Wait();
};

func (this *AtomicNeuron) DoneWaitGroup() {
  this.simulationWaitGroup.Done();
};

func (this *AtomicNeuron) GetInnerWaitGroup() *sync.WaitGroup {
  return this.innerSignal;
};

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
      inputSum := float64(0);
      for _, connection := range this.GetConnections() {
        if !connection.IsWriteable() {
          // Sum the connections to the neuron.
          inputSum += connection.GetOutput();
          fmt.Println(this.GetId(),"connection recieved!");
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
          fmt.Println(this.GetId(),"connection sent!");
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
          fmt.Println(this.GetId(),"connection sent!");
        }
      }
      return true;
    });
  }

  // Create an atomic neuron.
  var simulationWaitGroup sync.WaitGroup;
  var innerSignal sync.WaitGroup;
  mutexSignal := new(sync.Mutex);
  startSimulation := make(chan struct{});
  var atomicSignalCondition uint64 = 0;
  atomicNeuron := NewAtomicNeuron(startSimulation, &simulationWaitGroup, mutexSignal, &innerSignal, &atomicSignalCondition, len(this.neurons));

  // The first neuron starts the simulation ahead of all the others.
  // It grabs the outer lock, blocking all other neurons.
  // It then calculates its input based on the external connection inputting that are ready.
  // It then increments an atomic counter.
  // It unlocks the outer lock so the next neuron can go through.
  // It then calculates its output and sets the output connections values and sets them to be read.
  // The neuron then unlocks the next neuron to go through and waits for it to finish.
  // After all the neurons finish, the first neuron goes again first.
  // This then repeats over the time series...
  for _, neuron := range this.neurons {
    atomicNeuron.AddWaitGroup(1);
    go neuron.Simulate(simulation, atomicNeuron);
    time.Sleep(time.Millisecond * 2);
  }

  // Wait for the simulation to complete.
  atomicNeuron.FinishWaitGroup();
  fmt.Println("Finally Finished...");
};
