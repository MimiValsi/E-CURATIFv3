(async () => {
  const getJson = async () => {
    const response = await fetch("http://localhost:3001/jsonGraph");
    const data = await response.json();
    const js = await JSON.stringify(data);
    const jj = await JSON.parse(js)

    return jj;
  };

  const jsData = await getJson();

  let nbCuratifs = [];
  for (let i = 0; i < jsData.length; i++) {
    nbCuratifs.push(jsData[i].curatifs);
  }

  let nomSources = [];
  for (let i = 0; i < jsData.length; i++) {
    nomSources.push(jsData[i].name);
  }

  let codeGMAO = []
  for (let i = 0; i < jsData.length; i++) {
    codeGMAO.push(jsData[i].code_GMAO)
  }

  var chart = bb.generate({
    bindto: "#myPlot",

    data: {
      names: {
        data1: "Curatifs",
        data2: "Sources"
      },
      columns: [
        ["Curatifs en cours", ...nbCuratifs],
      ],
      type: "bar",
    },

    axis: {
      x: {
        type: "category",
        categories: [...codeGMAO],
        height: 50,
        tick: {
          rotate: 75,
          multiline: false,
        }
      }
    },

    size: {
      width: 1000,
      height: 400
    },

    padding: true,

    resize: true,

    zoom: {
      enabled: true,
      type: "drag"
    },

    legend: {
      position: "inset"
    },

    bar: {
      width: {
        ratio: 0.5
      }
    }
  });
})();
