package sn;

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
  spikes int;
  previousState *SpikingNeuron;
};

func NewSpikingNeuron(a, b, c, d float64) *SpikingNeuron {
  return &SpikingNeuron{
    a: a,
    b: b,
    c: c,
    d: d,
    V: defaultV,
    u: b * defaultV,
    spikes: 0,
    previousState: &SpikingNeuron{
      a: a,
      b: b,
      c: c,
      d: d,
      V: defaultV,
      u: b * defaultV,
      spikes: 0,
      previousState: nil,
    },
  };
};

func (this *SpikingNeuron) Reset() {
  currentState := this;
  this = this.previousState;
  this.previousState = currentState;
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

func (this *SpikingNeuron) Simulate(input float64, simulation *Simulation, predicate, success, fail DecisionFunction) {
  I := float64(0);

  steps := simulation.GetSteps();
  tau := simulation.GetTau();
  timeSeries := simulation.GetTimeSeries();
  start := simulation.GetStart();
  T1 := simulation.GetT();

  uu := make([]float64, len(timeSeries));

  for t, i := start, 0; t < steps; t, i = t + tau, i + 1 {
    timeSeries[i] = t;

    if t > T1 {
      I = input;
    } else {
      I = defaultInputCase;
    }

    this.V = this.V + tau * (constantV1 * (this.V * this.V) + constantV2 * this.V + constantV3 - this.u + I);
    this.u = this.u + tau * this.a * (this.b * this.V - this.u);

    if predicate(t, i, this) {
      success(t, i, this);
    } else {
      fail(t, i, this);
    }

    uu[i] = this.u;
  }

  simulation.SetTimeSeries(timeSeries);
};
