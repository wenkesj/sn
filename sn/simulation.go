package sn;

type Simulation struct {
  steps float64;
  tau float64;
  start float64;
  T float64;
  sketch map[float64]float64;
  timeSeries []float64;
};

func NewSimulation(steps, tau, start, T float64) *Simulation {
  numberOfIterations := int(steps * (1 / tau));
  timeSeries := make([]float64, numberOfIterations);
  sketch := make(map[float64]float64);
  return &Simulation{
    steps: steps,
    tau: tau,
    start: start,
    T: T,
    timeSeries: timeSeries,
    sketch: sketch,
  };
};

func (this *Simulation) GetSketch() map[float64]float64 {
  return this.sketch;
};

func (this *Simulation) SetSketch(sketch map[float64]float64) {
  this.sketch = sketch;
};

func (this *Simulation) SetSpikeRate(key, val float64) {
  this.sketch[key] = val;
};

func (this *Simulation) GetSpikeRate(key float64) float64 {
  return this.sketch[key];
};

func (this *Simulation) GetTimeSeries() []float64 {
  return this.timeSeries;
};

func (this *Simulation) SetTimeSeries(timeSeries []float64) {
  this.timeSeries = timeSeries;
};

func (this *Simulation) GetTau() float64 {
  return this.tau;
};

func (this *Simulation) SetTau(tau float64) {
  this.tau = tau;
};

func (this *Simulation) GetStart() float64 {
  return this.start;
};

func (this *Simulation) SetStart(start float64) {
  this.start = start;
};

func (this *Simulation) GetT() float64 {
  return this.T;
};

func (this *Simulation) SetT(T float64) {
  this.T = T;
};

func (this *Simulation) GetSteps() float64 {
  return this.steps;
};

func (this *Simulation) SetSteps(steps float64) {
  this.steps = steps;
};
