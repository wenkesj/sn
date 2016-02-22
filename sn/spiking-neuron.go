package sn;

import (
  // "fmt";
  "time";
  "strconv";
  "github.com/wenkesj/sn/sim";
  "github.com/wenkesj/sn/vars";
  "github.com/wenkesj/sn/group";
  "github.com/garyburd/redigo/redis";
);

type DecisionFunction func(float64, int, *SpikingNeuron) bool;
type FloatDecisionFunction func(int, float64, float64, *SpikingNeuron) bool;
type ReturnFloatFunction func(int, float64, float64, *SpikingNeuron) float64;

type Connection struct {
  weight float64;
  to *SpikingNeuron;
  from *SpikingNeuron;
  writeable bool;
};

func NewConnection(to *SpikingNeuron, from *SpikingNeuron, weight float64, writeable bool) *Connection {
  return &Connection{
    weight: weight,
    to: to,
    from: from,
    writeable: writeable,
  };
};

func (this *Connection) GetOutput() float64 {
  // Atomically load the output of the connection.
  // fmt.Println("GET: " + strconv.FormatInt(this.GetFrom().GetId(), 10) + ".to." + strconv.FormatInt(this.GetTo().GetId(), 10));
  output, err := this.GetFrom().GetStore().Do("GET", strconv.FormatInt(this.GetFrom().GetId(), 10) + "." + strconv.FormatInt(this.GetTo().GetId(), 10));
  if err != nil {
    panic(err);
  }
  if output == nil {
    // fmt.Println("Output is nil");
  }
  outputFloatValue, err := strconv.ParseFloat(string(output.([]uint8)), 64);
  if err != nil {
    panic(err);
  }
  return outputFloatValue;
};

func (this *Connection) SetOutput(output float64) {
  // Atomically store the output of the connection.
  // fmt.Println("SET: " + strconv.FormatInt(this.GetTo().GetId(), 10) + ".to." + strconv.FormatInt(this.GetFrom().GetId(), 10));
  this.GetFrom().GetStore().Do("SET", strconv.FormatInt(this.GetTo().GetId(), 10) + "." + strconv.FormatInt(this.GetFrom().GetId(), 10), output * this.GetWeight());
};

func (this *Connection) GetTo() *SpikingNeuron {
  return this.to;
};

func (this *Connection) GetFrom() *SpikingNeuron {
  return this.from;
};

func (this *Connection) GetWeight() float64 {
  return this.weight;
};

func (this *Connection) IsWriteable() bool {
  return this.writeable;
};

func (this *Connection) SetWriteabel(writeable bool) {
  this.writeable = writeable;
};

func (this *Connection) GetWriteabel() bool {
  return this.writeable;
};

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
  redisConnection redis.Conn;
  spikes int;
  spikeRateMap map[float64]float64;
  input float64;
  outputs []float64;
  spikeMap []float64;
  connections []*Connection;
};

func NewSpikingNeuron(a, b, c, d float64, id int64, redisConnection redis.Conn) *SpikingNeuron {
  spikeRateMap := make(map[float64]float64);
  spikeMap := []float64{};
  return &SpikingNeuron{
    a: a,
    b: b,
    c: c,
    d: d,
    V: vars.GetDefaultV(),
    u: b * vars.GetDefaultV(),
    input: 0,
    outputs: nil,
    spikes: 0,
    redisConnection: redisConnection,
    spikeRateMap: spikeRateMap,
    spikeMap: spikeMap,
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
  this.V = vars.GetDefaultV();
  this.u = b * vars.GetDefaultV();
  this.outputs = nil;
  this.input = 0;
  this.spikes = 0;
};

func (this *SpikingNeuron) GetStore() redis.Conn {
  return this.redisConnection;
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

func (this *SpikingNeuron) SetTimeSpike(time float64) {
  this.spikeMap = append(this.spikeMap, time);
};

func (this *SpikingNeuron) GetTimeSpike() []float64 {
  return this.spikeMap;
};

func (this *SpikingNeuron) ScopedSimulation(I float64, i int, t, T1, tau float64, uu []float64, neuronManager *group.NeuronManager) {
  // Lock...
  if neuronManager != nil {
      // Increment to number of neurons...
    // fmt.Println("Neuron", this.GetId(), "incremented the wait counter...");
    time.Sleep(time.Millisecond);
    neuronManager.GetInnerWaitGroup().Add(1);
    defer neuronManager.GetInnerWaitGroup().Done();
    // fmt.Println("Neuron", this.GetId(), "grabbed the lock...");
    time.Sleep(time.Millisecond);
    neuronManager.OuterLock();
  }

  // Get all the outputs from each connection.
  if this.GetInputPredicate()(i, t, T1, this) {
    I = this.GetInputSuccess()(i, t, T1, this);
  } else {
    I = this.GetInputFail()(i, t, T1, this);
  }

  this.SetV(this.GetV() + tau * (vars.GetConstantV1() * (this.GetV() * this.GetV()) + vars.GetConstantV2() * this.GetV() + vars.GetConstantV3() - this.GetU() + I));
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
}

func (this *SpikingNeuron) Simulate(simulation *sim.Simulation, neuronManager *group.NeuronManager) {
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
    this.ScopedSimulation(I, i, t, T1, tau, uu, neuronManager);

    if neuronManager != nil {
      // fmt.Println("Neuron", this.GetId(), "released the lock");
      time.Sleep(time.Millisecond);
      neuronManager.OuterUnlock();
      // fmt.Println("Neuron", this.GetId(), "is waiting...");
      time.Sleep(time.Millisecond);
      neuronManager.GetInnerWaitGroup().Wait();
      time.Sleep(time.Millisecond * 2);
      // fmt.Println("Neuron", this.GetId(), "is finished at time", t);
    }
  }

  simulation.SetTimeSeries(timeSeries);

  if neuronManager != nil {
    // fmt.Println("Finished", this.GetId());
    neuronManager.DoneWaitGroup();
  }
};
