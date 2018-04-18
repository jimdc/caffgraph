// set the dimensions and csvMargins of the graph
var csvMargin = {top: 20, right: 20, bottom: 30, left: 50},
    csvWidth = 960 - csvMargin.left - csvMargin.right,
    csvHeight = 500 - csvMargin.top - csvMargin.bottom;

// parse the date / time
var csvParseTime = d3.timeParse("%Y-%m-%dT%H:%M:%SZ");

// set the ranges
var csvX = d3.scaleTime().range([0, csvWidth]);
var csvY = d3.scaleLinear().range([csvHeight, 0]);

// define the line
var valueline = d3.line()
    .x(function(d) { return csvX(d.date); })
    .y(function(d) { return csvY(d.close); });

// append the svg obgect to the body of the page
// appends a 'group' element to 'csvSvg'
// moves the 'group' element to the top left margin
var csvSvg = d3.select("body").append("svg")
    .attr("width", csvWidth + csvMargin.left + csvMargin.right)
    .attr("height", csvHeight + csvMargin.top + csvMargin.bottom)
  .append("g")
    .attr("transform",
          "translate(" + csvMargin.left + "," + csvMargin.top + ")");

// Get the data
d3.csv("caff.csv", function(error, data) {
  if (error) throw error;

  // format the data
  data.forEach(function(d) {
      d.date = csvParseTime(d.date);
      d.close = +d.close;
  });

  // Scale the range of the data
  csvX.domain(d3.extent(data, function(d) { return d.date; }));
  csvY.domain([0, d3.max(data, function(d) { return d.close; })]);

  // Add the valueline path.
  csvSvg.append("path")
      .data([data])
      .attr("class", "line")
      .attr("d", valueline);

  // Add the X Axis
  csvSvg.append("g")
      .attr("transform", "translate(0," + csvHeight + ")")
      .call(d3.axisBottom(csvX));

  // Add the Y Axis
  csvSvg.append("g")
      .call(d3.axisLeft(csvY));

});
