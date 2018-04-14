# caffgraph
This program consists of two components: 

1. a Go program that converts a timestamped dosage of caffeine to a list of estimates for how much remains in the system 
2. a D3.js page that visualizes in line graph form, from that Go output, how much caffeine is in the body at one time

# adding entries
```
> go run halflife.go
> Enter mg: 200
2018-04-14T18:14 200mg
2018-04-14T23:56 100mg
2018-04-15T05:38 50mg
2018-04-15T11:20 25mg
2018-04-15T17:02 12mg
2018-04-15T22:44 6mg
2018-04-16T04:26 3mg
2018-04-16T10:08 1mg
```

Features to add:

* output to a file (TSV, JSON) that the js can read
* read in from file, add existing caffeine to dosage
* webapp to write to dosage file then run halflife.go
