package sn;

import (
  "sync";
);

type AtomicNeuron struct {
  simulationWaitGroup *sync.WaitGroup;
  innerWaitGroup *sync.WaitGroup;
  mutexSignal *sync.Mutex;
  numberOfOtherNeurons int;
};

func NewAtomicNeuron(
  simulationWaitGroup *sync.WaitGroup,
  innerWaitGroup *sync.WaitGroup,
  mutexSignal *sync.Mutex,
  numberOfOtherNeurons int,
) *AtomicNeuron {

  return &AtomicNeuron{
    simulationWaitGroup: simulationWaitGroup,
    innerWaitGroup: innerWaitGroup,
    mutexSignal: mutexSignal,
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

func (this *AtomicNeuron) GetInnerWaitGroup() *sync.WaitGroup {
  return this.innerWaitGroup;
};

func (this *AtomicNeuron) FinishWaitGroup() {
  this.simulationWaitGroup.Wait();
};

func (this *AtomicNeuron) DoneWaitGroup() {
  this.simulationWaitGroup.Done();
};
