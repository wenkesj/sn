package group;

import (
  "sync";
);

type NeuronManager struct {
  simulationWaitGroup *sync.WaitGroup;
  innerWaitGroup *sync.WaitGroup;
  mutexSignal *sync.Mutex;
  numberOfOtherNeurons int;
};

func NewNeuronManager(
  simulationWaitGroup *sync.WaitGroup,
  innerWaitGroup *sync.WaitGroup,
  mutexSignal *sync.Mutex,
  numberOfOtherNeurons int,
) *NeuronManager {

  return &NeuronManager{
    simulationWaitGroup: simulationWaitGroup,
    innerWaitGroup: innerWaitGroup,
    mutexSignal: mutexSignal,
    numberOfOtherNeurons: numberOfOtherNeurons,
  };
};

func (this *NeuronManager) GetNumber() int {
  return this.numberOfOtherNeurons;
};

func (this *NeuronManager) AddWaitGroup(num int) {
  this.simulationWaitGroup.Add(num);
};

func (this *NeuronManager) OuterLock() {
  this.mutexSignal.Lock();
};

func (this *NeuronManager) OuterUnlock() {
  this.mutexSignal.Unlock();
};

func (this *NeuronManager) GetInnerWaitGroup() *sync.WaitGroup {
  return this.innerWaitGroup;
};

func (this *NeuronManager) FinishWaitGroup() {
  this.simulationWaitGroup.Wait();
};

func (this *NeuronManager) DoneWaitGroup() {
  this.simulationWaitGroup.Done();
};
