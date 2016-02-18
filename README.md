# Phasic Spiking Neuron #
Model and Simulation of the Phasic Spiking Neuron.

```sh
go get -u github.com/wenkesj/sn
```

# Disclaimer #
This model has current issues:
+ Neurons only support one way connections.
+ Trying to find the perfect timing functions to make the simulations seem truly real time.
+ As the number of connections and neurons increase, more issues with `WaitGroup` and `connections` arise.

# Usage #
This package includes baseline simulation tools for creating and simulating [Neural Networks](https://en.wikipedia.org/wiki/Artificial_neural_network).
Currently, this package only has capabilities for simulating the Izhikevich Phasic Spiking Neuron model. Feel free to contribute!

For real examples, look at the [test directory](https://github.com/wenkesj/phasic-spiking-neuron/tree/master/tests).

# Testing and Plotting #
```sh
go test tests/sn_test.go
```

# License #
Copyright (c) 2016 Sam Wenke

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
