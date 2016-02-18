package sn;

import (
  "fmt";
  "sync/atomic";
  "unsafe";
);

type Connection struct {
  output float64;
  weight float64;
  to *SpikingNeuron;
  from *SpikingNeuron;
  writeable bool;
  ready chan bool;
  maxChannelLength int;
  outputAddress *unsafe.Pointer;
};

func NewConnection(to *SpikingNeuron, from *SpikingNeuron, weight float64, writeable bool, outputAddress *unsafe.Pointer) *Connection {
  ready := make(chan bool, 2);
  maxChannelLength := 2;
  return &Connection{
    outputAddress: outputAddress,
    weight: weight,
    to: to,
    from: from,
    writeable: writeable,
    ready: ready,
    maxChannelLength: maxChannelLength,
  };
};

func (this *Connection) GetReady() bool {
  currentReady := <-this.ready;
  this.ready = make(chan bool, this.maxChannelLength + 1);
  return currentReady;
};

func (this *Connection) SetReady(ready bool) {
  this.ready = make(chan bool, this.maxChannelLength + 1);
  this.ready <- ready;
};

func (this *Connection) GetOutput() float64 {
  // Atomically load the output of the connection.
  fmt.Println("Getting pointer", atomic.LoadPointer(this.outputAddress));
  outputPointer := (*float64)(atomic.LoadPointer(this.outputAddress));
  return *outputPointer;
};

func (this *Connection) SetOutput(output float64) {
  // Atomically store the output of the connection.
  output = this.GetWeight() * output;
  outputPointer := unsafe.Pointer(&output);
  fmt.Println("Storing pointer", this.outputAddress);
  atomic.StorePointer(this.outputAddress, outputPointer);
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
