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
  let jsonLen = jsData.length;

  let nomSources = [];
  let codeGMAO = []
  let aRealiser = [];
  let enCours = [];
  let done = [];

  for (let i = 0; i < jsonLen; i++) {
    let jec = jsData[i].en_cours;
    let jar = jsData[i].a_realiser;
    let jd = jsData[i].done;

    nomSources.push(jsData[i].name);
    codeGMAO.push(jsData[i].code_GMAO);

    if (typeof jar == 'undefined') {
      aRealiser.push(0);
    } else {
      aRealiser.push(jsData[i].a_realiser);
    }

    if (typeof jec == 'undefined') {
      enCours.push(0);
    } else {
      enCours.push(jsData[i].en_cours);
    }

    if (typeof jd == 'undefined') {
      done.push(0);
    } else {
      done.push(jsData[i].done);
    }
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
      ],
      type: "bar",
    },

    color: {
      pattern: ['#ff0000', '#0000ff', '#00ff00']
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
