/*
 * Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
 * This file is licensed under the MIT license.
 * See the LICENSE file.
 */

window.onload = function() {
  // The user wants to create a new topic.
  document.getElementById("create-button").onclick = function() {
    this.style.display = "none";

    var form = this.parentNode.getElementsByTagName("form")[0];
    form.style.display = "block";
    form.getElementsByClassName("tiny-text")[0].focus();
  };

  // When the user click the "rename" link, show the rename form.
  document.getElementById("rename").onclick = function() {
    var f = document.getElementById("rename-form");
    f.style.display = "block";
    f.getElementsByClassName("tiny-text")[0].focus();
    return false;
  };

  // When the user clicks the "edit" link, show the form to edit the topic.
  document.getElementById("edit").onclick = function() {
    var body = document.getElementById("edit-body");
    body.getElementsByClassName("contents-body")[0].style.display = "none";
    body.getElementsByClassName("contents-edit")[0].style.display = "block";
    document.getElementsByTagName("textarea")[0].focus();
    return false;
  };

  // When editing the topic, if the user clicks "Cancel", let's hide the whole
  // thing.
  document.getElementById("cancel-btn").onclick = function() {
    var body = document.getElementById("edit-body");
    body.getElementsByClassName("contents-edit")[0].style.display = "none";
    body.getElementsByClassName("contents-body")[0].style.display = "block";
  };

  // When the user clicks the "delete" link, show the confirmation form.
  document.getElementById("delete").onclick = function() {
    this.style.display = "none";
    document.getElementById("confirmation").style.display = "inline";
    return false;
  };

  // The current user really wants to delete this topic, do it.
  document.getElementById("delete-yes").onclick = function() {
    document.getElementById("delete-form").submit();
    return false;
  };

  // The current chickened out on deleting this topic.
  document.getElementById("delete-no").onclick = function() {
    document.getElementById("confirmation").style.display = "none";
    document.getElementById("delete").style.display = "inline";
    return false;
  };
}
