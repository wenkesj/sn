package sn;

type Connection struct {
  output float64;
  weight float64;
  to *SpikingNeuron;
  from *SpikingNeuron;
  writeable bool;
  ready bool;
};

func NewConnection(to *SpikingNeuron, from *SpikingNeuron, weight float64, writeable bool) *Connection {
  return &Connection{
    output: 0,
    weight: weight,
    to: to,
    from: from,
    writeable: writeable,
    ready: false,
  };
};

func (this *Connection) GetOutput() float64 {
  return this.output;
};

func (this *Connection) SetOutput(output float64) {
  this.output = output * this.GetWeight();
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

func (this *Connection) IsReady() bool {
  return ready;
};

func (this *Connection) SetReady(ready bool) {
  this.ready = ready;
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
