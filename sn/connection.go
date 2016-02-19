package sn;

import (
  "fmt";
  "strconv";
  "github.com/garyburd/redigo/redis";
);

type Connection struct {
  weight float64;
  to *SpikingNeuron;
  from *SpikingNeuron;
  writeable bool;
  redisConnection redis.Conn;
};

func NewConnection(to *SpikingNeuron, from *SpikingNeuron, weight float64, writeable bool, redisConnection redis.Conn) *Connection {
  return &Connection{
    weight: weight,
    to: to,
    from: from,
    writeable: writeable,
    redisConnection: redisConnection,
  };
};

// func Float64frombytes(bytes []byte) float64 {
//   bits := binary.LittleEndian.Uint64(bytes)
//   float := math.Float64frombits(bits)
//   return float
// }

func (this *Connection) GetOutput() float64 {
  // Atomically load the output of the connection.
  fmt.Println("GET: " + strconv.FormatInt(this.GetFrom().GetId(), 10) + ".to." + strconv.FormatInt(this.GetTo().GetId(), 10));
  output, err := this.redisConnection.Do("GET", strconv.FormatInt(this.GetFrom().GetId(), 10) + "." + strconv.FormatInt(this.GetTo().GetId(), 10));
  if err != nil {
    panic(err);
  }
  if output == nil {
    fmt.Println("Output is nil");
  }
  outputFloatValue, err := strconv.ParseFloat(string(output.([]uint8)), 64);
  return outputFloatValue;
};

func (this *Connection) SetOutput(output float64) {
  // Atomically store the output of the connection.
  fmt.Println("SET: " + strconv.FormatInt(this.GetTo().GetId(), 10) + ".to." + strconv.FormatInt(this.GetFrom().GetId(), 10));
  this.redisConnection.Do("SET", strconv.FormatInt(this.GetTo().GetId(), 10) + "." + strconv.FormatInt(this.GetFrom().GetId(), 10), output * this.GetWeight());
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
