package sn;

import (
  "fmt";
  "time";
  "sync";
);

type AtomicNeuron struct {
  simulationWaitGroup *sync.WaitGroup;
  mutexSignal *sync.Mutex;
  innerSignal *sync.WaitGroup;
  numberOfOtherNeurons int;
};

func NewAtomicNeuron(
  simulationWaitGroup *sync.WaitGroup,
  mutexSignal *sync.Mutex,
  innerSignal *sync.WaitGroup,
  numberOfOtherNeurons int,
) *AtomicNeuron {

  return &AtomicNeuron{
    simulationWaitGroup: simulationWaitGroup,
    mutexSignal: mutexSignal,
    innerSignal: innerSignal,
    numberOfOtherNeurons: numberOfOtherNeurons,
  };
};

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
  // Wait for the other one to increment the counter first so this neuron doesn't go through.
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
