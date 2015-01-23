$(document).ready(function() {
  $(".hide-properties-with-defaults").on("click", function() {
    $("body").toggleClass("properties-with-defaults-hidden");
    return false;
  });

  $(".hide-property-descriptions").on("click", function() {
    $("body").toggleClass("property-descriptions-hidden");
    return false;
  });
});
