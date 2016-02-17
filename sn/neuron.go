package sn;

import (
  "sync";
  uuid "github.com/satori/go.uuid";
);

// Simulation constants.
var defaultV = float64(-64);
var constantV1 = float64(0.04);
var constantV2 = float64(5);
var constantV3 = float64(140);

type DecisionFunction func(float64, int, *SpikingNeuron) bool;
type FloatDecisionFunction func(float64, float64, *SpikingNeuron) bool;
type ReturnFloatFunction func(float64, float64, *SpikingNeuron) float64;

type SpikingNeuron struct {
  a float64;
  b float64;
  c float64;
  d float64;
  V float64;
  u float64;
  id string;
  inputPredicate FloatDecisionFunction;
  inputSuccess ReturnFloatFunction;
  inputFail ReturnFloatFunction;
  predicate DecisionFunction;
  success DecisionFunction;
  fail DecisionFunction;
  spikes int;
  spikeRateMap map[float64]float64;
  input float64;
  connections []*Connection;
};

func NewSpikingNeuron(a, b, c, d float64) *SpikingNeuron {
  id := uuid.NewV4().String();
  spikeRateMap := make(map[float64]float64);
  return &SpikingNeuron{
    a: a,
    b: b,
    c: c,
    d: d,
    V: defaultV,
    u: b * defaultV,
    input: 0,
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

func (this *SpikingNeuron) GetId() string {
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
  newConnection := NewConnection(targetNeuron, weight, writeable);
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
    if connection.GetTarget().GetId() == targetNeuron.GetId() {
      this.connections = append(this.connections[:index], this.connections[index+1:]...);
      if once == 1 {
        return;
      }
      targetNeuron.RemoveConnection(this, 1);
      break;
    }
  }
};

func (this *SpikingNeuron) Simulate(simulation *Simulation, startSimulation chan struct{}, simulationWaitGroup *sync.WaitGroup, neuronWaitGroup *sync.WaitGroup) {
  I := float64(0);

  steps := simulation.GetSteps();
  tau := simulation.GetTau();
  timeSeries := simulation.GetTimeSeries();
  start := simulation.GetStart();
  T1 := simulation.GetT();

  uu := make([]float64, len(timeSeries));

  for t, i := start, 0; t < steps; t, i = t + tau, i + 1 {
    if startSimulation != nil && simulationWaitGroup != nil && neuronWaitGroup != nil {
      <- startSimulation;
      neuronWaitGroup.Add(1);
    }

    timeSeries[i] = t;

    if this.GetInputPredicate()(t, T1, this) {
      I = this.GetInputSuccess()(t, T1, this);
    } else {
      I = this.GetInputFail()(t, T1, this);
    }

    this.SetV(this.GetV() + tau * (constantV1 * (this.GetV() * this.GetV()) + constantV2 * this.GetV() + constantV3 - this.GetU() + I));
    this.SetU(this.GetU() + tau * this.GetA() * (this.GetB() * this.GetV() - this.GetU()));

    if this.GetPredicate()(t, i, this) {
      // Fire...
      this.GetSuccess()(t, i, this);
      // Default results.
      this.SetV(this.GetC());
      this.SetU(this.GetU() + this.GetD());
    } else {
      // Don't Fire...
      this.GetFail()(t, i, this);
    }

    if neuronWaitGroup != nil {
      neuronWaitGroup.Done();
      neuronWaitGroup.Wait();
    }

    uu[i] = this.GetU();
  }

  simulation.SetTimeSeries(timeSeries);

  if simulationWaitGroup != nil {
    simulationWaitGroup.Done();
  }
};
