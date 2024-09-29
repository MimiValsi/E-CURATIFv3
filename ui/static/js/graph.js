(async () => {
  const getJson = async () => {
    const response = await fetch("http://localhost:3001/jsonGraph");

    const data = await response.json();
    const js = JSON.stringify(data);
    const jp = await JSON.parse(js)

    return jp;
  };


  const getJsonCD = async () => {
    const curatifDone = await fetch("http://localhost:3001/curatifDone");
    const data = await curatifDone.json();
    const js = JSON.stringify(data);
    const jp = await JSON.parse(js);

    return jp;
  }

  // Fetch every active 'Curatif'
  const jsData = await getJson();

  let nomSources = [];
  let codeGMAO = []
  let aRealiser = [];
  let enCours = [];
  let done = [];
  let total = [];

  for (let i = 0; i < jsData.length; i++) {
    aRealiser.push(jsData[i].a_realiser);
    enCours.push(jsData[i].en_cours);
    done.push(jsData[i].done);
    nomSources.push(jsData[i].name);
    codeGMAO.push(jsData[i].code_GMAO);
    total.push(jsData[i].curatifs);
  }

  var chart = bb.generate({
    bindto: "#myPlot",

    data: {
      names: {
        data1: "Curatifs",
        data2: "Sources"
      },
      columns: [
        ["Curatifs à réaliser", ...aRealiser],
        ["Curatifs en cours", ...enCours],
        ["Curatifs réalisées", ...done],
        ["Total", ...total]
      ],
      type: "bar",

      groups: [
        [
          "Curatifs à réaliser",
          "Curatifs en cours",
          "Curatifs réalisées"
        ]
      ]
    },

    color: {
      pattern: ["#cc1111", "#0080ff", "#99cc33", "#999999"]
    },

    axis: {
      x: {
        type: "category",
        categories: [...codeGMAO],
        height: 50,
      }
    },

    size: {
      width: 1000,
      height: 400
    },

    padding: true,

    // resize: {
    //   auto: true,
    //   timer: 100
    // },

    zoom: {
      enabled: true,
      type: "drag"
    },

    legend: {
      position: "bottom"
    },

    bar: {
      width: {
        ratio: 0.5
      }
    }
  });
})();
