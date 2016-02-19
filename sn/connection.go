package sn;

import (
  // "fmt";
  "strconv";
);

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
