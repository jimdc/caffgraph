// set the dimensions and margins of the graph
var margin = {top: 20, right: 40, bottom: 30, left: 50},
    width = 800 - margin.left - margin.right,
    height = 500 - margin.top - margin.bottom;

// parse the date / time
var parseTime = d3.timeParse("%Y-%m-%dT%H:%M:%SZ");

// set the ranges
var x = d3.scaleTime().range([0, width]);
var y = d3.scaleLinear().range([height, 0]);

// define the lines
var dosageline = d3.line()
    .x(function(d) { return x(d.time); })
    .y(function(d) { return y(d.dosage); });

var remnantlines = [];
var remnantdata = [];

// append the svg obgect to the body of the page
// appends a 'group' element to 'svg'
// moves the 'group' element to the top left margin
var svg = d3.select("body").append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
  .append("g")
    .attr("transform",
          "translate(" + margin.left + "," + margin.top + ")");

// Get the data
d3.json("caffeine.json", function(error, data) {
  if (error) throw error;

  // format the data
  data.forEach(function(d) {
      d.time = parseTime(d.time);
      d.dosage = +d.dosage;

      var remnants = d.remnants;
      if (remnants) {

        remnants.forEach(function(r) {
          r.time = parseTime(r.time);  
        });

        // d.remnants.remnant does not work
        // maybe the domain needs to be scaled for max?
        remnantlines.push(
          d3.line()
          .x(function(d) { return x(d.time); })
          .y(function(d) { return y(d.remnants); })
        );
        remnantdata.push(remnants)
      }
  });

  // Scale the range of the data
  x.domain(d3.extent(data, function(d) { return d.time; }));
  // Not the true max, which could come from summing up dosages
  y.domain([0, d3.max(data, function(d) { return d.dosage; })]);

  // Add the dosageline path.
  svg.append("path")
      .data([data])
      .attr("class", "line")
      .attr("d", dosageline);

  // Setup for dosing circles
  var color = d3.scaleOrdinal(d3.schemeCategory10);
  var doses = color.domain().map(function(name) {
    return {
      name: name,
      values: data.map(function(d) {
        return {date: d.time, dosage: +d[dosage]};
      })
    };
  });
  
  var dose = svg.selectAll(".dose")
      .data(doses)
      .enter()
      .append("g")
      .attr("class", "dose");

  // Add the remnant lines' path: does not work.
  // Loops twice which is correct, but data doesn't work... 
  remnantlines.forEach(function(r, index) {
      svg.append("path")
          .data([data])
          .attr("class", "line")
          .attr("d", r);
  });

  // Add the X Axis
  svg.append("g")
      .attr("transform", "translate(0," + height + ")")
      .call(d3.axisBottom(x));

  // text label for the x axis
  svg.append("text")             
      .attr("transform",
            "translate(" + (width/2) + " ," + 
                           (height + margin.top + 10) + ")")
      .style("text-anchor", "middle")
      .text("Time");

  // Add the Y Axis
  svg.append("g")
      .call(d3.axisLeft(y));

  // text label for the y axis
  svg.append("text")
      .attr("transform", "rotate(-90)")
      .attr("y", 0 - margin.left)
      .attr("x",0 - (height / 2))
      .attr("dy", "1em")
      .style("text-anchor", "middle")
      .text("Caffeine (mg)");

  // text label for dosage line
  var lowestDose = d3.min(data, function(d) { return d.dosage; });
  svg.append("text")
      .attr("transform", "translate(" + (width+3) + "," + y(lowestDose) + ")")
      .attr("dy", ".35em")
      .attr("text-anchor", "start")
      .style("fill", "steelblue")
      .text("dose");
});
