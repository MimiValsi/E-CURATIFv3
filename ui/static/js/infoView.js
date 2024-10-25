// Confirme si l'utilisateur veut vraiment supprimer une info
let del = document.getElementById("check-delete");
del.onclick = () => {
  return confirm("Êtes-vous sûr(e) de supprimer?");
};
