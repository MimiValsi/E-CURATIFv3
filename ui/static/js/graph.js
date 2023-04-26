const data = [{
  x:nomSources,
  y:nbCuratifs,
  type:"bar"
}];

const layout = {
  title:"Nombre de Curatifs en cours",
  width: 1000,
  height: 400,
  autosize: true
};

var config = {
  modeBarButtonsToRemove: ['sendDataToCloud', 'autoScale2d', 'hoverClosestCartesian', 'hoverCompareCartesian', 'lasso2d', 'select2d', 'toggleSpikelines', 'pan2d', 'zoomIn2d', 'zoomOut2d', 'zoom2d', 'drawclosedpath'],
  displaylogo: false,
  showTips: true,
  responsive: true
};

Plotly.newPlot("myPlot", data, layout, config);
