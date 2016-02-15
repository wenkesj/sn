# Phasic Spiking Neuron #
**Homework Assignment for Intelligent Systems**

# Overview #
The objective of this experiment is to simulate and model the membrane potential of a neuron using **Izhikevich's model**.

**Izhikevich's neuron** is modeled by the following differential equations:

![Equation 1](/assets/eq1.png)

conditionally,

![Equation 2](/assets/eq2.png)

The model generates the full shape of a action potential. As seen in the [Plots](#plots) section, a **Izhikevich** neuron was modeled with 5 different constant _a_, _b_, _c_, and _d_ parameters and variable inputs _I_.

+ [Problem 1](#problem-1--single-neuron-analysis) simulates a neuron and drives it with external input(s).
+ [Problem 2](#problem-2--neuron-a-and-b) simulates a simple 2 neuron network.

# Results #

## Problem 1 – Single Neuron Analysis ##

**What kind of pattern do you see, and what does it say about the behavior of the neuron?**

**In what sense is the neuron performing any useful function here?**

## Problem 2 – Neuron A and B ##

**In particular, discuss the differences (if any) in the pattern of outputs seen for neurons A and B, and how they may arise**

# Plots #
#### [Figure 1](#figure-1) ####
Spiking Neuron at I = 1.0
![Spiking Neuron at I = 1.0](/tests/plots/spiking-neuron-1.000000.png)

#### [Figure 2](#figure-2) ####
Spiking Neuron at I = 5.0
![Spiking Neuron at I = 5.0](/tests/plots/spiking-neuron-5.000000.png)

#### [Figure 3](#figure-3) ####
Spiking Neuron at I = 10.0
![Spiking Neuron at I = 10.0](/tests/plots/spiking-neuron-10.000000.png)

#### [Figure 4](#figure-4) ####
Spiking Neuron at I = 15.0
![Spiking Neuron at I = 15.0](/tests/plots/spiking-neuron-15.000000.png)

#### [Figure 5](#figure-5) ####
Spiking Neuron at I = 20.0
![Spiking Neuron at I = 20.0](/tests/plots/spiking-neuron-20.000000.png)

#### [Figure 6](#figure-6) ####
Mean Spike Rate
![Mean Spike Rate](/tests/plots/spiking-neuron-mean-spike-rate.png)

# Conclusion #

# Code #
[All the code can be found here.](https://github.com/wenkesj/sn)

**This repository uses the following packages:**
+ `github.com/satori/go.uuid`

# Install #
```sh
go install github.com/wenkesj/sn/sn
```

# Testing and Plotting #
Run the test suite to generate the plots.
```sh
go test tests/sn_test.go
```
