let getFile = document.getElementById("getFile");
let imprt = document.getElementById("importCSV");
if (getFile && imprt) {
  imprt.onchange = () => {
    getFile.submit();
    alert("Fichier en cours de transfert...");
  };
}
