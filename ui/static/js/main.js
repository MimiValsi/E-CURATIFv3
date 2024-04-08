function searchStatus() {
  let input, filter, table, tr, td, i, txtValue;

  // Fetch for the <input> searcg tag
  input = document.getElementById("searchStatus");
  filter = input.value.toUpperCase(); // It simplifies the research
  table = document.getElementById("myTable");
  tr = table.getElementsByTagName("tr"); // Fetch every <tr> inside table

  /* Inside "sourceView" <table id="myTable"> there's 2 <tr>
   * tr[1] contains 2 <td>. We're going to look inside td[1]
   * which represents "Status" column.
   * And coz td[0] doesn't work... Added to TODO list
   */

  for (i = 1; i < tr.length; i++) {
    // For Info search change to ...("td")[0];
    // For Priority search change to ...("td")[1];
    // For Status search change to ...("td")[2];
    td = tr[i].getElementsByTagName("td")[2];
    //                               here ^

    if (td) {
      txtValue = td.textContent;
      console.log(txtValue);

      if (txtValue.toUpperCase().indexOf(filter) > -1) {
        tr[i].style.display = "";
      } else {
        tr[i].style.display = "none";
      }
    }
  }
}

function searchPS() {
  let input, filter, table, tr, td, i, txtValue;

  // Same function as searchStatus() but for home page
  input = document.getElementById("searchPS");
  filter = input.value.toUpperCase();
  table = document.getElementById("homeTable");
  tr = table.getElementsByTagName("tr");

  for (i = 1; i < tr.length; i++) {
    td = tr[i].getElementsByTagName("td")[0];

    if (td) {
      txtValue = td.textContent || td.innerText;

      if (txtValue.toUpperCase().indexOf(filter) > -1) {
        tr[i].style.display = "";
      } else {
        tr[i].style.display = "none";
      }
    }
  }
}

// function checkInpt() {
//   let input = document.forms["srcInpt"]["name"].value;

//   if (input == "" || input == null) {
//     alert("Le champ ne doit pas être vide!");
//     document.getElementById("name").classList.remove('inpt')
//     document.getElementById("name").classList.add('inptAlert')
//     return false;
//   }
// }

document.addEventListener("DOMContentLoaded", function () {
  let elements = document.getElementsByTagName("INPUT");
  let srcName = document.getElementById("name");

  let infoAgent = document.getElementById("agent");
  let infoMaterial = document.getElementById("material");
  let infoDetail = document.getElementById("detail");
  let infoEvent = document.getElementById("event");
  let infoPriority = document.getElementById("priority");

  for (let i = 0; i < elements.length; i++) {
    elements[i].oninvalid = function (e) {
      e.target.setCustomValidity("");
      if (!e.target.validity.valid) {
        e.target.setCustomValidity("Ce champ ne doit pas être vide");
        if (srcName != null) {
          srcName.classList.add("inptAlert");
        }

        if (infoAgent != null) {
          infoAgent.classList.add("inptAlert");
        }
        if (infoMaterial != null) {
          infoMaterial.classList.add("inptAlert");
        }
        if (infoDetail != null) {
          infoDetail.classList.add("inptAlert");
        }
        if (infoEvent != null) {
          infoEvent.classList.add("inptAlert");
        }
        if (infoPriority != null) {
          infoPriority.classList.add("inptAlert");
        }
      }
    };
    elements[i].oninput = function (e) {
      e.target.setCustomValidity("");
    };
  }
});
