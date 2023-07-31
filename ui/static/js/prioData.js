(async () => {
  const gJson = async () => {
    const response = await fetch("http://localhost:3001/prioData");

    const data = await response.json();
    const jsonS = JSON.stringify(data);
    const jsonP = await JSON.parse(jsonS);

    return jsonP;
  };

  const jsonData = await gJson();
  let jsLength = jsonData.length;

  let jsData = [];
  for (let i = 0; i < jsLength; i++) {
    let data = jsonData[i].material + ":" + " " + jsonData[i].detail
    jsData.push(data);
    document.getElementById("infoUrgent").innerHTML = jsData.join(' / ');
  }

  // document.getElementById("infoUrgent").innerHTML = jsData.join(' / ');
})();

