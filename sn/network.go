package sn;

import (
  "time";
  "fmt"
  "sync";
  "sync/atomic";
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
  conditionalSignal *sync.Cond;
  atomicSignalCondition *uint64;
  numberOfOtherNeurons int;
};

func NewAtomicNeuron(
  startSimulation StartChannel,
  simulationWaitGroup *sync.WaitGroup,
  conditionalSignal *sync.Cond,
  atomicSignalCondition *uint64,
  numberOfOtherNeurons int,
) *AtomicNeuron {

  return &AtomicNeuron{
    startSimulation: startSimulation,
    simulationWaitGroup: simulationWaitGroup,
    conditionalSignal: conditionalSignal,
    atomicSignalCondition: atomicSignalCondition,
    numberOfOtherNeurons: numberOfOtherNeurons,
  };
};

func (this *AtomicNeuron) GetStartChannel() StartChannel {
  return this.startSimulation;
}

func (this *AtomicNeuron) AddWaitGroup(num int) {
  this.simulationWaitGroup.Add(num);
};

func (this *AtomicNeuron) Lock() {
  this.conditionalSignal.L.Lock();
};

func (this *AtomicNeuron) Unlock() {
  this.conditionalSignal.L.Unlock();
};

func (this *AtomicNeuron) Wait() {
  // First one gets here and increments a counter.
  // This says the neuron has updated its value and has set the input of another neuron.
  this.IncrementSignal();

  // Next, this neuron lets the next neuron go through.
  this.Unlock();

  // Next, the neuron will wait until all of the other neurons recieved the signal and processed it.
  for {
    // The first neuron should be doing his thing first...
    loadedAtomicValue := atomic.LoadUint64(this.atomicSignalCondition);
    if loadedAtomicValue != uint64(this.numberOfOtherNeurons) {
      // Do nothing...
    } else if loadedAtomicValue == uint64(0) {
      break;
    } else {
      // The first neuron should go through, reset the atomic counter and grab the lock before anyone else.
      // Reset the counter.
      atomic.StoreUint64(this.atomicSignalCondition, uint64(0));
      break;
    }
  }
};

func (this *AtomicNeuron) FinishWaitGroup() {
  this.simulationWaitGroup.Wait();
};

func (this *AtomicNeuron) IncrementSignal() {
  atomic.AddUint64(this.atomicSignalCondition, 1);
};

func (this *AtomicNeuron) DoneWaitGroup() {
  this.simulationWaitGroup.Done();
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
      // Wait for each connection that is not a writeable connection.
      inputSum := 0;
      for _, connection := range this.GetConnections() {
        if !connection.IsWriteable() {
          // Wait for the connection to be ready for the neuron to recieve input.
          for !connection.IsReady() {
            // ...
          }
          // Sum the connections to the neuron.
          inputSum += connection.GetOutput(i);
          connection.SetReady(false);
        }
      }

      // Ready the next neuron.
      this.SetInput(inputSum);
      return this.GetInput();
    });

    neuron.SetInputFail(func (i int, t, T1 float64, this *SpikingNeuron) float64 {
      this.SetInput(defaultInputCase);
      return defaultInputCase;
    });

    // Assign predicates.
    neuron.SetPredicate(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      return this.GetV() > defaultVCutoff;
    });

    neuron.SetSuccess(func (timeIndex float64, currentIndex int, this *SpikingNeuron) bool {
      this.SetOutput(currentIndex, defaultVCutoff);

      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(defaultOutputMembranePotentialSuccess);
          connection.SetReady(true);
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
      this.SetOutput(currentIndex, this.GetV());

      // Send the output to the connected neuron.
      for _, connection := range this.GetConnections() {
        if connection.IsWriteable() {
          connection.SetOutput(defaultOutputMembranePotentialFail);
          connection.SetReady(true);
        }
      }
      return true;
    });
  }

  // Create an atomic neuron.
  var simulationWaitGroup sync.WaitGroup;
  conditionalMutex := new(sync.Mutex);
  conditionalSignal := sync.NewCond(conditionalMutex);
  startSimulation := make(chan struct{});
  var atomicSignalCondition uint64 = 0;
  atomicNeuron := NewAtomicNeuron(startSimulation, &simulationWaitGroup, conditionalSignal, &atomicSignalCondition, len(this.neurons));

  // The first neuron starts the simulation ahead of all the others.
  // It grabs the lock, blocking all other neurons.
  // It then calculates its input based on the external connection inputting that are ready.
  // It then increments an atomic counter.
  // It then calculates its output and sets the output connections values and sets them to be read.
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
  fmt.Println("Finally Finished...");
};
