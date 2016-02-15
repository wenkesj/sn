package sn;

import (
  // "github.com/CHH/eventemitter";
  uuid "github.com/satori/go.uuid"
);

type NetworkNeuron struct {
  neuron *SpikingNeuron;
  id string;
  connections []*Connection;
};

func NewNetworkNeuron(neuron *SpikingNeuron) *NetworkNeuron {
  id := uuid.NewV4().String();
  return &NetworkNeuron{
    id: id,
    neuron: neuron,
    connections: nil,
  };
};

func (this *NetworkNeuron) GetConnections() []*Connection {
  return this.connections;
};

func (this *NetworkNeuron) GetId() string {
  return this.id;
};

func (this *NetworkNeuron) CreateConnection(targetNeuron *NetworkNeuron, weight float64, once int) {
  if this.connections == nil {
    this.connections = []*Connection{};
  }
  newConnection := NewConnection(targetNeuron, weight);
  this.connections = append(this.connections, newConnection);
  if once == 1 {
    return;
  }
  targetNeuron.CreateConnection(this, weight, 1);
};

func (this *NetworkNeuron) RemoveConnection(targetNeuron *NetworkNeuron, once int) {
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
