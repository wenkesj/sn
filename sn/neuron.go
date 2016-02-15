package sn;

import (
  uuid "github.com/satori/go.uuid";
);

// Simulation constants.
var defaultInputCase = float64(0);
var defaultV = float64(-64);
var constantV1 = float64(0.04);
var constantV2 = float64(5);
var constantV3 = float64(140);

type DecisionFunction func(float64, int, *SpikingNeuron) bool;

type SpikingNeuron struct {
  a float64;
  b float64;
  c float64;
  d float64;
  V float64;
  u float64;
  id string;
  predicate DecisionFunction;
  success DecisionFunction;
  fail DecisionFunction;
  spikes int;
  input float64;
  connections []*Connection;
};

func NewSpikingNeuron(a, b, c, d float64) *SpikingNeuron {
  id := uuid.NewV4().String();
  return &SpikingNeuron{
    a: a,
    b: b,
    c: c,
    d: d,
    V: defaultV,
    u: b * defaultV,
    input: 0,
    spikes: 0,
    predicate: nil,
    success: nil,
    fail: nil,
    id: id,
    connections: nil,
  };
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

func (this *SpikingNeuron) Simulate(input float64, simulation *Simulation) {
  I := float64(0);
  this.SetInput(input);

  steps := simulation.GetSteps();
  tau := simulation.GetTau();
  timeSeries := simulation.GetTimeSeries();
  start := simulation.GetStart();
  T1 := simulation.GetT();

  uu := make([]float64, len(timeSeries));

  for t, i := start, 0; t < steps; t, i = t + tau, i + 1 {
    timeSeries[i] = t;

    if t > T1 {
      I = this.GetInput();
    } else {
      I = defaultInputCase;
    }

    this.SetV(this.GetV() + tau * (constantV1 * (this.GetV() * this.GetV()) + constantV2 * this.GetV() + constantV3 - this.GetU() + I));
    this.SetU(this.GetU() + tau * this.GetA() * (this.GetB() * this.GetV() - this.GetU()));

    if this.GetPredicate()(t, i, this) {
      // Fire...
      this.GetSuccess()(t, i, this);
    } else {
      // Don't Fire...
      this.GetFail()(t, i, this);
    }

    uu[i] = this.GetU();
  }

  simulation.SetTimeSeries(timeSeries);
};
