package sn;

import (
  "time";
  "fmt";
);

// Simulation constants.
var defaultV = float64(-64);
var constantV1 = float64(0.04);
var constantV2 = float64(5);
var constantV3 = float64(140);

type DecisionFunction func(float64, int, *SpikingNeuron) bool;
type FloatDecisionFunction func(int, float64, float64, *SpikingNeuron) bool;
type ReturnFloatFunction func(int, float64, float64, *SpikingNeuron) float64;

type SpikingNeuron struct {
  a float64;
  b float64;
  c float64;
  d float64;
  V float64;
  u float64;
  id int64;
  inputPredicate FloatDecisionFunction;
  inputSuccess ReturnFloatFunction;
  inputFail ReturnFloatFunction;
  predicate DecisionFunction;
  success DecisionFunction;
  fail DecisionFunction;
  spikes int;
  spikeRateMap map[float64]float64;
  input float64;
  outputs []float64;
  connections []*Connection;
};

func NewSpikingNeuron(a, b, c, d float64, id int64) *SpikingNeuron {
  spikeRateMap := make(map[float64]float64);
  return &SpikingNeuron{
    a: a,
    b: b,
    c: c,
    d: d,
    V: defaultV,
    u: b * defaultV,
    input: 0,
    outputs: nil,
    spikes: 0,
    spikeRateMap: spikeRateMap,
    predicate: nil,
    success: nil,
    fail: nil,
    inputPredicate: nil,
    inputSuccess: nil,
    inputFail: nil,
    id: id,
    connections: nil,
  };
};

func (this *SpikingNeuron) ResetParameters(a, b, c, d float64) {
  this.a = a;
  this.b = b;
  this.c = c;
  this.d = d;
  this.V = defaultV;
  this.u = b * defaultV;
  this.outputs = nil;
  this.input = 0;
  this.spikes = 0;
};

func (this *SpikingNeuron) SetV(V float64) {
  this.V = V;
};

func (this *SpikingNeuron) GetV() float64 {
  return this.V;
};

func (this *SpikingNeuron) SetU(u float64) {
  this.u = u;
};

func (this *SpikingNeuron) GetU() float64 {
  return this.u;
};

func (this *SpikingNeuron) GetA() float64 {
  return this.a;
};

func (this *SpikingNeuron) GetB() float64 {
  return this.b;
};

func (this *SpikingNeuron) GetC() float64 {
  return this.c;
};

func (this *SpikingNeuron) GetD() float64 {
  return this.d;
};

func (this *SpikingNeuron) SetOutputs(outputs []float64) {
  this.outputs = outputs;
};

func (this *SpikingNeuron) SetOutput(index int, output float64) {
  this.outputs[index] = output;
};

func (this *SpikingNeuron) GetOutputs() []float64 {
  return this.outputs;
};

func (this *SpikingNeuron) GetOutput(i int) float64 {
  return this.outputs[i];
};

func (this *SpikingNeuron) SetSpikeRate(key, val float64) {
  this.spikeRateMap[key] = val;
};

func (this *SpikingNeuron) GetSpikeRate() map[float64]float64 {
  return this.spikeRateMap;
};

func (this *SpikingNeuron) SetSpikes(spikes int) {
  this.spikes = spikes;
};

func (this *SpikingNeuron) GetSpikes() int {
  return this.spikes;
};

func (this *SpikingNeuron) GetInput() float64 {
  return this.input;
};

func (this *SpikingNeuron) SetInput(input float64) {
  this.input = input;
};

func (this *SpikingNeuron) GetPredicate() DecisionFunction {
  return this.predicate;
};

func (this *SpikingNeuron) SetPredicate(predicate DecisionFunction) {
  this.predicate = predicate;
};

func (this *SpikingNeuron) GetSuccess() DecisionFunction {
  return this.success;
};

func (this *SpikingNeuron) SetSuccess(success DecisionFunction) {
  this.success = success;
};

func (this *SpikingNeuron) GetFail() DecisionFunction {
  return this.fail;
};

func (this *SpikingNeuron) SetFail(fail DecisionFunction) {
  this.fail = fail;
};

func (this *SpikingNeuron) GetConnections() []*Connection {
  return this.connections;
};

func (this *SpikingNeuron) GetId() int64 {
  return this.id;
};

func (this *SpikingNeuron) GetInputPredicate() FloatDecisionFunction {
  return this.inputPredicate;
};

func (this *SpikingNeuron) GetInputSuccess() ReturnFloatFunction {
  return this.inputSuccess;
};

func (this *SpikingNeuron) GetInputFail() ReturnFloatFunction {
  return this.inputFail;
};

func (this *SpikingNeuron) SetInputPredicate(inputFunction FloatDecisionFunction) {
  this.inputPredicate = inputFunction;
};

func (this *SpikingNeuron) SetInputSuccess(inputFunction ReturnFloatFunction) {
  this.inputSuccess = inputFunction;
};

func (this *SpikingNeuron) SetInputFail(inputFunction ReturnFloatFunction) {
  this.inputFail = inputFunction;
};

func (this *SpikingNeuron) CreateConnection(targetNeuron *SpikingNeuron, weight float64, writeable bool, once int) {
  if this.connections == nil {
    this.connections = []*Connection{};
  }
  newConnection := NewConnection(targetNeuron, this, weight, writeable);
  this.connections = append(this.connections, newConnection);
  if once == 1 {
    return;
  }
  targetNeuron.CreateConnection(this, weight, !writeable, 1);
};

func (this *SpikingNeuron) RemoveConnection(targetNeuron *SpikingNeuron, once int) {
  if this.connections == nil {
    return;
  }
  for index, connection := range this.connections {
    if connection.GetTo().GetId() == targetNeuron.GetId() {
      this.connections = append(this.connections[:index], this.connections[index+1:]...);
      if once == 1 {
        return;
      }
      targetNeuron.RemoveConnection(this, 1);
      break;
    }
  }
};

func (this *SpikingNeuron) Simulate(simulation *Simulation, atomicNeuron *AtomicNeuron) {
  I := float64(0);

  steps := simulation.GetSteps();
  tau := simulation.GetTau();
  timeSeries := simulation.GetTimeSeries();
  start := simulation.GetStart();
  T1 := simulation.GetT();

  uu := make([]float64, len(timeSeries));

  this.SetOutputs(make([]float64, len(timeSeries)));

  for t, i := start, 0; t < steps; t, i = t + tau, i + 1 {
    timeSeries[i] = t;

    // Lock...
    if atomicNeuron != nil {
      // Increment to number of neurons...
      fmt.Println("Neuron", this.GetId(), "incremented the counter");
      atomicNeuron.GetInnerWaitGroup().Add(1);
      fmt.Println("Neuron", this.GetId(), "grabbed the lock");
      atomicNeuron.OuterLock();
      time.Sleep(time.Millisecond);
    }

    // Get all the outputs from each connection.
    if this.GetInputPredicate()(i, t, T1, this) {
      I = this.GetInputSuccess()(i, t, T1, this);
    } else {
      I = this.GetInputFail()(i, t, T1, this);
    }

    this.SetV(this.GetV() + tau * (constantV1 * (this.GetV() * this.GetV()) + constantV2 * this.GetV() + constantV3 - this.GetU() + I));
    this.SetU(this.GetU() + tau * this.GetA() * (this.GetB() * this.GetV() - this.GetU()));

    if this.GetPredicate()(t, i, this) {
      this.GetSuccess()(t, i, this);
      // Default results.
      this.SetV(this.GetC());
      this.SetU(this.GetU() + this.GetD());
    } else {
      // Don't Fire...
      this.GetFail()(t, i, this);
    }

    uu[i] = this.GetU();

    // Wait for all other Neurons to finish their computation.
    if atomicNeuron != nil {
      fmt.Println("Neuron", this.GetId(), "released the lock");
      atomicNeuron.Wait(this);
      fmt.Println("Neuron", this.GetId(), "is finished at time", t);
    }
  }

  simulation.SetTimeSeries(timeSeries);

  if atomicNeuron != nil {
    fmt.Println("Finished", this.GetId());
    atomicNeuron.DoneWaitGroup();
  }
};
