package sn;

// Simulation constants.
var defaultInputCase = float64(0);
var defaultVCutoff = float64(30);
var defaultV = float64(-64);
var constantV1 = float64(0.04);
var constantV2 = float64(5);
var constantV3 = float64(140);

type SpikingNeuron struct {
  a float64;
  b float64;
  c float64;
  d float64;
  V float64;
  u float64;
  spikeRate float64;
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
    spikeRate: float64(0),
    previousState: &SpikingNeuron{
      a: a,
      b: b,
      c: c,
      d: d,
      V: defaultV,
      u: b * defaultV,
      spikeRate: float64(0),
      previousState: nil,
    },
  };
};

func (this *SpikingNeuron) Reset() {
  currentState := this;
  this = this.previousState;
  this.previousState = currentState;
};

func (this *SpikingNeuron) SetSpikeRate(mean float64) {
  this.spikeRate = mean;
};

func (this *SpikingNeuron) GetSpikeRate() float64 {
  return this.spikeRate;
};

func (this *SpikingNeuron) Simulate(input float64, simulation *Simulation) []float64 {
  spikes := 0;
  I := float64(0);

  steps := simulation.GetSteps();
  tau := simulation.GetTau();
  timeSeries := simulation.GetTimeSeries();
  start := simulation.GetStart();
  measureStart := simulation.GetMeasureStart();
  T1 := simulation.GetT();

  VV, uu := make([]float64, len(timeSeries)), make([]float64, len(timeSeries));

  for t, i := start, 0; t < steps; t, i = t + tau, i + 1 {
    timeSeries[i] = t;

    if t > T1 {
      I = input;
    } else {
      I = defaultInputCase;
    }

    this.V = this.V + tau * (constantV1 * (this.V * this.V) + constantV2 * this.V + constantV3 - this.u + I);
    this.u = this.u + tau * this.a * (this.b * this.V - this.u);

    if (this.V > defaultVCutoff) {
      VV[i] = defaultVCutoff;
      this.V = this.c;
      this.u = this.u + this.d;
      if t > measureStart {
        spikes++;
      }
    } else {
      VV[i] = this.V;
    }
    uu[i] = this.u;
  }
  this.SetSpikeRate(float64(spikes) / measureStart);
  simulation.SetTimeSeries(timeSeries);
  return VV;
};
