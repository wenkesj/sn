package sn;

type Connection struct {
  weight float64;
  target *NetworkNeuron;
};

func NewConnection(target *NetworkNeuron, weight float64) *Connection {
  return &Connection{
    weight: weight,
    target: target,
  };
};

func (this *Connection) GetTarget() *NetworkNeuron {
  return this.target;
};

func (this *Connection) GetWeight() float64 {
  return this.weight;
};
